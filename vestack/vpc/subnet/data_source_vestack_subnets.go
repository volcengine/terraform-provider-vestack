package subnet

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

func DataSourceVestackSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVestackSubnetsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Subnet IDs.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of Subnet.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of Subnet query.",
			},
			"subnets": {
				Description: "The collection of Subnet query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of Subnet.",
						},
						"subnet_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Name of Subnet.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Vpc ID of Subnet.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Status of Subnet.",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cidr block of Subnet.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of Subnet.",
						},
						"available_ip_address_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The count of available ip address.",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of Zone.",
						},
						"creation_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of Subnet.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time of Subnet.",
						},
						"total_ipv4_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Count of ipv4.",
						},
						"route_table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of route table.",
						},
						"route_table_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of route table.",
						},
					},
				},
			},
		},
	}
}

func dataSourceVestackSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	subnetService := NewSubnetService(meta.(*ve.SdkClient))
	return subnetService.Dispatcher.Data(subnetService, d, DataSourceVestackSubnets())
}