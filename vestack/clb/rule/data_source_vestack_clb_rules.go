package rule

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

func DataSourceVestackRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVestackRulesRead,
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Id of listener.",
			},
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Rule IDs.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"rules": {
				Description: "The collection of Rule query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Id of Rule.",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Id of Rule.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Domain of Rule.",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Url of Rule.",
						},
						"server_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Id of Server Group.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Description of Rule.",
						},
					},
				},
			},
		},
	}
}

func dataSourceVestackRulesRead(d *schema.ResourceData, meta interface{}) error {
	ruleService := NewRuleService(meta.(*ve.SdkClient))
	return ve.DefaultDispatcher().Data(ruleService, d, DataSourceVestackRules())
}
