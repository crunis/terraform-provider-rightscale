package rs

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rs/rs/rsc"
)

// Example:
//
// data "rs_cm_server_template" "mysql" {
//     filter {
//         name = "Database Manager for MySQL"
//         revision = 24
//     }
// }

func dataSourceCMServerTemplate() *schema.Resource {
	return &schema.Resource{
		Read: resourceServerTemplateRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "name of ServerTemplate, partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"revision": {
							Type:        schema.TypeInt,
							Description: "revision of ServerTemplate, use 0 to match latest non-committed version",
							Optional:    true,
							ForceNew:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "description of ServerTemplate, partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"lineage": {
							Type:        schema.TypeString,
							Description: "lineage of ServerTemplate",
							Optional:    true,
							ForceNew:    true,
						},
						"multi_cloud_image": {
							Type:        schema.TypeString,
							Description: "ID of ServerTemplate multi cloud image resource",
							Optional:    true,
							ForceNew:    true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"revision": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"lineage": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceServerTemplateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	loc := &rsc.Locator{Namespace: "rs_cm", Type: "server_templates"}

	res, err := client.List(loc, "", cmFilters(d))
	if err != nil {
		return err
	}

	if len(res) == 0 {
		return nil
	}
	for k, v := range res[0].Fields {
		d.Set(k, v)
	}
	d.SetId(res[0].Locator.Href)
	return nil
}