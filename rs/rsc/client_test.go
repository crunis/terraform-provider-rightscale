package rsc

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/rightscale/rsc.v6/rsapi"
)

// NOTE: the tests below make "real" API requests to the RightScale platform.
// They use credentials passed in environment variables to perform auth. The
// tests are skipped if the environment variables are missing.
//
// The tests use the following environment variables:
//
//     * RS_API_TOKEN is the API token used to auth API requests made to RightScale.
//     * RS_PROJECT_ID is the RightScale project used to run the tests.
//     * DEBUG causes additional output useful to troubleshoot issues.

func TestAuthenticate(t *testing.T) {
	token := validToken(t)
	project := validProjectID(t)
	cases := []struct {
		Name      string
		Token     string
		ProjectID int
		Error     string
	}{
		{"valid", token, project, ""},
		{"invalid-token", "foo", project, "failed to authenticate"},
		{"invalid-project-id", token, 0, "session does not give access to project 0"},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			_, err := New(c.Token, c.ProjectID)
			if c.Error == "" {
				if err != nil {
					t.Errorf("got error %q, expected none", err)
				}
			} else {
				if err == nil {
					t.Errorf("got no error, expected %q", c.Error)
				} else {
					if err.Error() != c.Error {
						t.Errorf("got error %q, expected %q", err.Error(), c.Error)
					}
				}
			}
		})
	}
}

func TestCreate(t *testing.T) {
	const (
		namespace = "rs_cm"
		typ       = "deployment"
		depl      = "Terraform Provider Test Deployment"
		deplDesc  = "Created by tests"
	)
	token := validToken(t)
	project := validProjectID(t)
	cl, err := New(token, project)
	if err != nil {
		t.Fatal(err)
	}
	rs := cl.(*client).rs
	cleanDeployment(t, depl, rs)
	defer cleanDeployment(t, depl, rs)

	_, err = cl.Create(namespace, typ, Fields{"name": depl, "description": deplDesc})

	if err != nil {
		t.Errorf("got error %q, expected none", err)
		return
	}
	d := showDeployment(t, depl, rs)
	if d == nil {
		t.Errorf("deployment not created")
		return
	}
	if d["name"].(string) != depl {
		t.Errorf("got deployment with name %v, expected %q", d["name"], depl)
	}
	if d["description"].(string) != deplDesc {
		t.Errorf("got deployment with description %v, expected %q", d["description"], deplDesc)
	}
}

func TestDelete(t *testing.T) {
	const (
		namespace = "rs_cm"
		typ       = "deployment"
		depl      = "Terraform Provider Test Deployment"
		deplDesc  = "Created by tests"
	)
	token := validToken(t)
	project := validProjectID(t)
	cl, err := New(token, project)
	if err != nil {
		t.Fatal(err)
	}
	rs := cl.(*client).rs
	cleanDeployment(t, depl, rs)
	defer cleanDeployment(t, depl, rs)

	res, err := cl.Create(namespace, typ, Fields{"name": depl, "description": deplDesc})
	if err != nil {
		t.Fatal(err)
	}

	err = cl.Delete(res.Locator)

	if err != nil {
		t.Errorf("got error %q, expected none", err)
		return
	}
	d := showDeployment(t, depl, rs)
	if d != nil {
		t.Errorf("deployment not deleted")
		return
	}
}

func validToken(t *testing.T) string {
	tok := os.Getenv("RS_API_TOKEN")
	if tok == "" {
		t.Skip("RS_API_TOKEN environment variable not defined, skipping authentication test")
	}
	return tok
}

func validProjectID(t *testing.T) int {
	pid := os.Getenv("RS_PROJECT_ID")
	if pid == "" {
		t.Skip("RS_PROJECT_ID environment variable not defined")
	}
	projectID, err := strconv.Atoi(pid)
	if err != nil {
		t.Fatal(err)
	}
	return projectID
}

func showDeployment(t *testing.T, depl string, rs *rsapi.API) map[string]interface{} {
	req, err := rs.BuildHTTPRequest("GET", "/api/deployments", "1.5", rsapi.APIParams{"filter[]": "name==" + depl}, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := rs.PerformRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to retrieve deployment: index returned %q", resp.Status)
	}
	ms, err := rs.LoadResponse(resp)
	if err != nil {
		t.Fatal(err)
	}
	if len(ms.([]interface{})) == 0 {
		return nil
	}
	return ms.([]interface{})[0].(map[string]interface{})
}

func cleanDeployment(t *testing.T, depl string, rs *rsapi.API) {
	var id string
	{
		m := showDeployment(t, depl, rs)
		if m == nil {
			return
		}
		links := m["links"].([]interface{})
		var href string
		for _, l := range links {
			rel := l.(map[string]interface{})["rel"].(string)
			if rel != "self" {
				continue
			}
			href = l.(map[string]interface{})["href"].(string)
			break
		}
		idx := strings.LastIndex(href, "/")
		id = href[idx+1:]
		if id == "" {
			t.Fatalf("failed to retrieve deployment id, href: %q", href)
		}
	}
	req, err := rs.BuildHTTPRequest("DELETE", "/api/deployments/"+id, "1.5", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := rs.PerformRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("failed to delete deployment: destroy returned %q", resp.Status)
	}
}