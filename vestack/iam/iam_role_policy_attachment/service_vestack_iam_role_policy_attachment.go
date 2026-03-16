package iam_role_policy_attachment

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamRolePolicyAttachmentService struct {
	Client *bp.SdkClient
}

func NewIamRolePolicyAttachmentService(c *bp.SdkClient) *VestackIamRolePolicyAttachmentService {
	return &VestackIamRolePolicyAttachmentService{
		Client: c,
	}
}

func (s *VestackIamRolePolicyAttachmentService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamRolePolicyAttachmentService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	action := "ListAttachedRolePolicies"
	logger.Debug(logger.ReqFormat, action, m)
	if m == nil {
		resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &m)
		if err != nil {
			return data, err
		}
	}

	logger.Debug(logger.RespFormat, action, m, *resp)

	results, err = bp.ObtainSdkValue("Result.AttachedPolicyMetadata", *resp)
	if err != nil {
		return data, err
	}
	if results == nil {
		results = []interface{}{}
	}
	if data, ok = results.([]interface{}); !ok {
		return data, errors.New("Result.AttachedPolicyMetadata is not Slice")
	}
	return data, err
}

func (s *VestackIamRolePolicyAttachmentService) ReadResource(resourceData *schema.ResourceData, roleId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if roleId == "" {
		roleId = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(roleId, ":")
	if len(ids) != 3 {
		return data, fmt.Errorf("import id is invalid")
	}
	req := map[string]interface{}{
		"RoleName": ids[0],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("value is not map")
		}
		pName, ok := data["PolicyName"].(string)
		if !ok {
			continue
		}
		pType, ok := data["PolicyType"].(string)
		if !ok {
			continue
		}
		if ids[1] == pName && ids[2] == pType {
			return data, err
		}
	}
	return data, fmt.Errorf("Role policy attachment %s not exist ", roleId)
}

func (s *VestackIamRolePolicyAttachmentService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *VestackIamRolePolicyAttachmentService) WithResourceResponseHandlers(rolePolicyAttachment map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return rolePolicyAttachment, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackIamRolePolicyAttachmentService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	createIamRolePolicyAttachmentCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AttachRolePolicy",
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				roleName, ok := d.Get("role_name").(string)
				if !ok {
					return errors.New("role_name is not string")
				}
				policyName, ok := d.Get("policy_name").(string)
				if !ok {
					return errors.New("policy_name is not string")
				}
				policyType, ok := d.Get("policy_type").(string)
				if !ok {
					return errors.New("policy_type is not string")
				}
				d.SetId(fmt.Sprintf("%s:%s:%s", roleName, policyName, policyType))
				return nil
			},
		},
	}
	return []bp.Callback{createIamRolePolicyAttachmentCallback}
}

func (s *VestackIamRolePolicyAttachmentService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackIamRolePolicyAttachmentService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	deleteRoleCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DetachRolePolicy",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ids := strings.Split(d.Id(), ":")
				(*call.SdkParam)["RoleName"] = ids[0]
				(*call.SdkParam)["PolicyName"] = ids[1]
				(*call.SdkParam)["PolicyType"] = ids[2]
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
							return resource.NonRetryableError(fmt.Errorf("error on reading iam role on delete %q, %w", d.Id(), callErr))
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
	return []bp.Callback{deleteRoleCallback}
}

func (s *VestackIamRolePolicyAttachmentService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ResponseConverts: map[string]bp.ResponseConvert{
			"PolicyName": {
				TargetField: "policy_name",
			},
			"PolicyType": {
				TargetField: "policy_type",
			},
			"PolicyTrn": {
				TargetField: "policy_trn",
			},
			"Description": {
				TargetField: "description",
			},
			"AttachDate": {
				TargetField: "attach_date",
			},
			"PolicyScope": {
				TargetField: "policy_scope",
			},
		},
		NameField:    "RoleName",
		CollectField: "policies",
	}
}

func (s *VestackIamRolePolicyAttachmentService) ReadResourceId(id string) string {
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
