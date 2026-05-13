package iam_caller_identity

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

func DataSourceVestackIamCallerIdentities() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVestackIamCallerIdentitiesRead,
		Schema: map[string]*schema.Schema{
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
			"caller_identities": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The collection of caller identities.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The account id.",
						},
						"trn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The trn.",
						},
						"identity_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The identity type.",
						},
						"identity_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The identity id.",
						},
					},
				},
			},
		},
	}
}

func dataSourceVestackIamCallerIdentitiesRead(d *schema.ResourceData, meta interface{}) error {
	service := NewIamCallerIdentityService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(service, d, DataSourceVestackIamCallerIdentities())
}
