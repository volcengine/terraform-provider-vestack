package security_group

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

func DataSourceVestackSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVestackSecurityGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of SecurityGroup IDs.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of vpc where security group is located.",
			},
			"security_group_names": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The list of security group name to query.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of SecurityGroup.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ProjectName of SecurityGroup.",
			},
			"tags": bp.TagsSchema(),
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of SecurityGroup query.",
			},
			"security_groups": {
				Description: "The collection of SecurityGroup query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of SecurityGroup.",
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of SecurityGroup.",
						},
						"security_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Name of SecurityGroup.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of Vpc.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Status of SecurityGroup.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of SecurityGroup.",
						},
						"creation_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of SecurityGroup.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A Name Regex of SecurityGroup.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ProjectName of SecurityGroup.",
						},
						"tags": bp.TagsSchemaComputed(),
					},
				},
			},
		},
	}
}

func dataSourceVestackSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	securityGroupService := NewSecurityGroupService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(securityGroupService, d, DataSourceVestackSecurityGroups())
}
