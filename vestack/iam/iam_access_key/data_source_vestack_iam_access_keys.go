package iam_access_key

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

func DataSourceVestackIamAccessKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVestackIamAccessKeysRead,
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user name.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of user query.",
			},
			"access_key_metadata": {
				Description: "The collection of access keys.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user name.",
						},
						"access_key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user access key id.",
						},
						"create_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user access key create date.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user access key status.",
						},
						"update_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user access key update date.",
						},
					},
				},
			},
		},
	}
}

func dataSourceVestackIamAccessKeysRead(d *schema.ResourceData, meta interface{}) error {
	iamAccessKeyService := NewIamAccessKeyService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(iamAccessKeyService, d, DataSourceVestackIamAccessKeys())
}
