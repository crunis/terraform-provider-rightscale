package rs

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rs/rs/rsc"
)

// Example:
//
// data "rs_cm_network" "infra_vpc" {
//     filter {
//         resource_uid = "vpc-c31ee987"
//         cloud = ${data.rs_cm_cloud.ec2_us_east_1.id}
//     }
// }

func dataSourceNetworks() *schema.Resource {
	return &schema.Resource{
		Read: resourceNetworkRead,

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
							Description: "name of network, uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"cloud": {
							Type:        schema.TypeString,
							Description: "ID of the network cloud",
							Optional:    true,
							ForceNew:    true,
						},
						"deployment": {
							Type:        schema.TypeString,
							Description: "ID of deployment resource that owns network",
							Optional:    true,
							ForceNew:    true,
						},
						"resource_uid": {
							Type:        schema.TypeString,
							Description: "cloud ID of network, e.g. 'vpc-2124fe46'",
							Optional:    true,
							ForceNew:    true,
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Description: "CIDR of the network resource",
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
			"resource_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_tenancy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceNetworkRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	loc := &rsc.Locator{Namespace: "rs_cm", Type: "networks"}

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
