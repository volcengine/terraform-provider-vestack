package image_share_permission

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

func DataSourceVestackImageSharePermissions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVestackImageSharePermissionsRead,
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the image.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of query.",
			},

			"accounts": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The shared account id of the image.",
						},
					},
				},
			},
		},
	}
}

func dataSourceVestackImageSharePermissionsRead(d *schema.ResourceData, meta interface{}) error {
	service := NewImageSharePermissionService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceVestackImageSharePermissions())
}
