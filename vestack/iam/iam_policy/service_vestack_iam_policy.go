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
		resp         *map[string]interface{}
		results      interface{}
		ok           bool
		allPolicies  []interface{}
		userPolicies []interface{}
		rolePolicies []interface{}
		temp         interface{}
		userName     string
		roleName     string
	)
	if userName, ok = m["UserName"].(string); ok {
		action := "ListAttachedUserPolicies"
		param := map[string]interface{}{
			"UserName": userName,
		}
		logger.Debug(logger.ReqFormat, action, param)
		resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &param)
		if err != nil {
			return data, err
		}
		temp, err = bp.ObtainSdkValue("Result.AttachedPolicyMetadata", *resp)
		if err != nil {
			return data, err
		}
		if temp != nil {
			if userPolicies, ok = temp.([]interface{}); !ok {
				return data, fmt.Errorf("%s Response AttachedPolicyMetadata not []interface{}", action)
			}
		}
		delete(m, "UserName")
	}

	if roleName, ok = m["RoleName"].(string); ok {
		action := "ListAttachedRolePolicies"
		param := map[string]interface{}{
			"RoleName": roleName,
		}
		logger.Debug(logger.ReqFormat, action, param)
		resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &param)
		if err != nil {
			return data, err
		}
		temp, err = bp.ObtainSdkValue("Result.AttachedPolicyMetadata", *resp)
		if err != nil {
			return data, err
		}
		if temp != nil {
			if rolePolicies, ok = temp.([]interface{}); !ok {
				return data, fmt.Errorf("%s Response AttachedPolicyMetadata not []interface{}", action)
			}
		}
		delete(m, "RoleName")
	}

	allPolicies, err = bp.WithPageOffsetQuery(m, "Limit", "Offset", 100, 0, func(condition map[string]interface{}) (data []interface{}, err error) {
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
	if err != nil {
		return data, err
	}

	data = allPolicies

	if len(userPolicies) > 0 {
		data, err = s.MergeAttachedPolicies(userPolicies, data, "UserName", userName, "UserAttachDate")
		if err != nil {
			return data, err
		}
	}

	if len(rolePolicies) > 0 {
		data, err = s.MergeAttachedPolicies(rolePolicies, data, "RoleName", roleName, "RoleAttachDate")
		if err != nil {
			return data, err
		}
	}

	return data, err
}

func (s *VestackIamPolicyService) ReadResource(resourceData *schema.ResourceData, policyId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
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
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("value is not map")
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
				d.SetId(policyName.(string))
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
					d.SetId(policyName.(string))
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
		ResponseConverts: map[string]bp.ResponseConvert{
			"PolicyName": {
				TargetField: "id",
				KeepDefault: true,
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

func (s *VestackIamPolicyService) MergeAttachedPolicies(attached []interface{}, source []interface{}, k, v, attachKey string) (data []interface{}, err error) {
	var (
		temp interface{}
	)
	for _, p0 := range attached {
		temp, err = bp.ObtainSdkValue("PolicyName", p0)
		if err != nil {
			return data, err
		}
		p0PolicyName := temp.(string)
		temp, err = bp.ObtainSdkValue("AttachDate", p0)
		if err != nil {
			return data, err
		}
		attachDate := temp.(string)
		for _, p1 := range source {
			temp, err = bp.ObtainSdkValue("PolicyName", p0)
			if err != nil {
				return data, err
			}
			p1PolicyName := temp.(string)
			if p0PolicyName == p1PolicyName {
				p1.(map[string]interface{})[attachKey] = attachDate
				p1.(map[string]interface{})[k] = v
				data = append(data, p1)
				break
			}
		}
	}
	return data, err
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "iam",
		Version:     "2018-01-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}
