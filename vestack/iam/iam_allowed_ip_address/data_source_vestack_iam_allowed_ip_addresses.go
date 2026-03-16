package iam_allowed_ip_address

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

func DataSourceVestackIamAllowedIpAddresses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVestackIamAllowedIpAddressesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"allowed_ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The collection of query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_ip_list": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable the IP whitelist.",
						},
						"quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The quota of the IP whitelist.",
						},
						"ip_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The IP whitelist list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description of the IP address.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceVestackIamAllowedIpAddressesRead(d *schema.ResourceData, meta interface{}) error {
	service := NewIamAllowedIpAddressService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(service, d, DataSourceVestackIamAllowedIpAddresses())
}
