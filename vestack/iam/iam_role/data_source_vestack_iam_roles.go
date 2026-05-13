package iam_role

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

func DataSourceVestackIamRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVestackIamRolesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of Role.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of Role query.",
			},
			"query": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Fuzzy query. Can query by role name, display name or description.",
			},
			"roles": {
				Description: "The collection of Role query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource name of the Role.",
						},
						"role_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Role.",
						},
						"create_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The create time of the Role.",
						},
						"update_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The update time of the Role.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the Role.",
						},
						"trust_policy_document": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The trust policy document of the Role.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the Role.",
						},
						"max_session_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The max session duration of the Role.",
						},
						"role_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The id of the Role.",
						},
						"is_service_linked_role": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the Role is a service linked role.",
						},
						"tags": bp.TagsSchemaComputed(),
					},
				},
			},
		},
	}
}

func dataSourceVestackIamRolesRead(d *schema.ResourceData, meta interface{}) error {
	iamRoleService := NewIamRoleService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(iamRoleService, d, DataSourceVestackIamRoles())
}
