package iam_account_summer

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
	"time"
)

type VestackIamAccountSummaryService struct {
	Client *bp.SdkClient
}

func NewIamAccountSummaryService(c *bp.SdkClient) *VestackIamAccountSummaryService {
	return &VestackIamAccountSummaryService{
		Client: c,
	}
}

func (s *VestackIamAccountSummaryService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamAccountSummaryService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	action := "GetAccountSummary"
	logger.Debug(logger.ReqFormat, action, m)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
	if err != nil {
		return data, err
	}

	result, err := bp.ObtainSdkValue("Result.SummaryMap", *resp)
	if err != nil {
		return data, err
	}
	if result == nil {
		return []interface{}{}, nil
	}

	return []interface{}{result}, nil
}

func (s *VestackIamAccountSummaryService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ResponseConverts: map[string]bp.ResponseConvert{
			"AccessKeysPerUserQuota":              {TargetField: "access_keys_per_user_quota"},
			"AccessKeysPerAccountQuota":           {TargetField: "access_keys_per_account_quota"},
			"AttachedPoliciesPerGroupQuota":       {TargetField: "attached_policies_per_group_quota"},
			"AttachedSystemPoliciesPerGroupQuota": {TargetField: "attached_system_policies_per_group_quota"},
			"AttachedPoliciesPerRoleQuota":        {TargetField: "attached_policies_per_role_quota"},
			"AttachedSystemPoliciesPerRoleQuota":  {TargetField: "attached_system_policies_per_role_quota"},
			"AttachedSystemPoliciesPerUserQuota":  {TargetField: "attached_system_policies_per_user_quota"},
			"AttachedPoliciesPerUserQuota":        {TargetField: "attached_policies_per_user_quota"},
			"GroupsQuota":                         {TargetField: "groups_quota"},
			"PoliciesQuota":                       {TargetField: "policies_quota"},
			"RolesQuota":                          {TargetField: "roles_quota"},
			"UsersQuota":                          {TargetField: "users_quota"},
			"PolicySize":                          {TargetField: "policy_size"},
			"RolesUsage":                          {TargetField: "roles_usage"},
			"UsersUsage":                          {TargetField: "users_usage"},
			"GroupsUsage":                         {TargetField: "groups_usage"},
			"PoliciesUsage":                       {TargetField: "policies_usage"},
			"GroupsPerUserQuota":                  {TargetField: "groups_per_user_quota"},
		},
		CollectField: "account_summaries",
	}
}

func (s *VestackIamAccountSummaryService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	return nil, nil
}
func (s *VestackIamAccountSummaryService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}
func (s *VestackIamAccountSummaryService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	return nil
}
func (s *VestackIamAccountSummaryService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}
func (s *VestackIamAccountSummaryService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}
func (s *VestackIamAccountSummaryService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return nil
}
func (s *VestackIamAccountSummaryService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "iam",
		Action:      actionName,
		Version:     "2018-01-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		RegionType:  bp.Global,
	}
}
