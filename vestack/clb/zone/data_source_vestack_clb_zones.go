package zone

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

func DataSourceVestackClbZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVestackClbZonesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of zone query.",
			},
			"master_zones": {
				Description: "The master zones list.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The master zone id.",
						},
						"slave_zones": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The slave zones list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The slave zone id.",
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

func dataSourceVestackClbZonesRead(d *schema.ResourceData, meta interface{}) error {
	service := NewClbZoneService(meta.(*ve.SdkClient))
	return ve.DefaultDispatcher().Data(service, d, DataSourceVestackClbZones())
}
