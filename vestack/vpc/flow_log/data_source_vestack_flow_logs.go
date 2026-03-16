package flow_log

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

func DataSourceVestackFlowLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVestackFlowLogsRead,
		Schema: map[string]*schema.Schema{
			"flow_log_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of flow log IDs.",
			},
			"flow_log_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of flow log.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of flow log.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of resource. Valid values: `vpc`, `subnet`, `eni`.",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of resource.",
			},
			"traffic_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of traffic. Valid values: `All`, `Allow`, `Drop`.",
			},
			"log_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of log project.",
			},
			"log_topic_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of log topic.",
			},
			"aggregation_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The aggregation interval of flow log. Unit: minute. Valid values: `1`, `5`, `10`.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of flow log. Valid values: `Active`, `Pending`, `Inactive`, `Creating`, `Deleting`.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of VPC.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The project name of flow log.",
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
			"flow_logs": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of flow log.",
						},
						"flow_log_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of flow log.",
						},
						"flow_log_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of flow log.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of flow log.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource. Valid values: `vpc`, `subnet`, `eni`.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of resource.",
						},
						"traffic_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of traffic. Valid values: `All`, `Allow`, `Drop`.",
						},
						"log_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of log project.",
						},
						"log_topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of log topic.",
						},
						"aggregation_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The aggregation interval of flow log. Unit: minute. Valid values: `1`, `5`, `10`.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of flow log. Valid values: `Active`, `Pending`, `Inactive`, `Creating`, `Deleting`.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of VPC.",
						},
						"business_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The business status of flow log.",
						},
						"lock_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reason why flow log is locked.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The created time of flow log.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The updated time of flow log.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project name of flow log.",
						},
						"tags": bp.TagsSchemaComputed(),
					},
				},
			},
		},
	}
}

func dataSourceVestackFlowLogsRead(d *schema.ResourceData, meta interface{}) error {
	service := NewFlowLogService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceVestackFlowLogs())
}
