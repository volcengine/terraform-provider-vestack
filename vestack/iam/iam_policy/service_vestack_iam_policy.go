package iam_policy

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamPolicyService struct {
	Client *bp.SdkClient
}

func NewIamPolicyService(c *bp.SdkClient) *VestackIamPolicyService {
	return &VestackIamPolicyService{
		Client: c,
	}
}

func (s *VestackIamPolicyService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamPolicyService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	// remove unsupport params
	delete(m, "UserName")
	delete(m, "RoleName")

	return bp.WithPageOffsetQuery(m, "Limit", "Offset", 100, 0, func(condition map[string]interface{}) (data []interface{}, err error) {
		action := "ListPolicies"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, err
			}
		}

		results, err = bp.ObtainSdkValue("Result.PolicyMetadata", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.PolicyMetadata is not Slice")
		}
		return data, err
	})
}

func (s *VestackIamPolicyService) ReadResource(resourceData *schema.ResourceData, policyId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if policyId == "" {
		policyId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"Query": policyId,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if temp, ok := v.(map[string]interface{}); !ok {
			return data, errors.New("value is not map")
		} else {
			if pName, ok := temp["PolicyName"].(string); ok && pName == policyId {
				data = temp
				break
			}
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Policy %s not exist ", policyId)
	}
	return data, err
}

func (s *VestackIamPolicyService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *VestackIamPolicyService) WithResourceResponseHandlers(policy map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		policy["Id"] = policy["PolicyName"]
		return policy, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackIamPolicyService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	createIamPolicyCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreatePolicy",
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				policyName, err := bp.ObtainSdkValue("Result.Policy.PolicyName", *resp)
				if err != nil {
					return err
				}
				if nameStr, ok := policyName.(string); ok && nameStr != "" {
					d.SetId(nameStr)
				}
				return nil
			},
		},
	}
	return []bp.Callback{createIamPolicyCallback}
}

func (s *VestackIamPolicyService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	updatePolicyCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdatePolicy",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["PolicyName"] = d.Get("policy_name")
				if d.HasChange("policy_name") {
					oldPolicyName, newPolicyName := d.GetChange("policy_name")
					(*call.SdkParam)["PolicyName"] = oldPolicyName
					(*call.SdkParam)["NewPolicyName"] = newPolicyName
				}
				if d.HasChange("policy_document") {
					(*call.SdkParam)["NewPolicyDocument"] = d.Get("policy_document")
				}
				if d.HasChange("description") {
					(*call.SdkParam)["NewDescription"] = d.Get("description")
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				if d.HasChange("policy_name") {
					policyName, err := bp.ObtainSdkValue("Result.Policy.PolicyName", *resp)
					if err != nil {
						return err
					}
					if nameStr, ok := policyName.(string); ok && nameStr != "" {
						d.SetId(nameStr)
					}
				}
				return nil
			},
		},
	}
	return []bp.Callback{updatePolicyCallback}
}

func (s *VestackIamPolicyService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	deletePolicyCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeletePolicy",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["PolicyName"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading iam policy on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
		},
	}
	return []bp.Callback{deletePolicyCallback}
}

func (s *VestackIamPolicyService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"with_service_role_policy": {
				TargetField: "WithServiceRolePolicy",
			},
			"scope": {
				TargetField: "Scope",
			},
		},
		ResponseConverts: map[string]bp.ResponseConvert{
			"PolicyName": {
				TargetField: "id",
				KeepDefault: true,
			},
			"Category": {
				TargetField: "category",
			},
			"AttachmentCount": {
				TargetField: "attachment_count",
			},
			"IsServiceRolePolicy": {
				TargetField: "is_service_role_policy",
			},
		},
		NameField:    "PolicyName",
		IdField:      "PolicyName",
		CollectField: "policies",
	}
}

func (s *VestackIamPolicyService) ReadResourceId(id string) string {
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
