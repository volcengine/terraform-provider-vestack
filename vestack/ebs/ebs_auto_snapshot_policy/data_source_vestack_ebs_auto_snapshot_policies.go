package ebs_auto_snapshot_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

func DataSourceVestackEbsAutoSnapshotPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVestackEbsAutoSnapshotPoliciesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of auto snapshot policy IDs.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The project name of auto snapshot policy.",
			},
			"tags": bp.TagsSchema(),
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of Resource.",
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

			"auto_snapshot_policies": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the auto snapshot policy.",
						},
						"auto_snapshot_policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the auto snapshot policy.",
						},
						"auto_snapshot_policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the auto snapshot policy.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the auto snapshot policy.",
						},
						"volume_nums": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of volumes associated with the auto snapshot policy.",
						},
						"time_points": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The creation time points of the auto snapshot policy. The value range is `0~23`, representing a total of 24 time points from 00:00 to 23:00, for example, 1 represents 01:00.",
						},
						"repeat_weekdays": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The date of creating snapshot repeatedly by week. The value range is `1-7`, for example, 1 represents Monday.",
						},
						"repeat_days": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Create snapshots repeatedly on a daily basis, with intervals of a certain number of days between each snapshot.",
						},
						"retention_days": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The retention days of the auto snapshot. `-1` means permanently preserving the snapshot.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the auto snapshot policy.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The updated time of the auto snapshot policy.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project name of the auto snapshot policy.",
						},
						"tags": bp.TagsSchemaComputed(),
					},
				},
			},
		},
	}
}

func dataSourceVestackEbsAutoSnapshotPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	service := NewEbsAutoSnapshotPolicyService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceVestackEbsAutoSnapshotPolicies())
}
