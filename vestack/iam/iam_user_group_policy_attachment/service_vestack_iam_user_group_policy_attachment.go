package iam_user_group_policy_attachment

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

type VestackIamUserGroupPolicyAttachmentService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewIamUserGroupPolicyAttachmentService(c *bp.SdkClient) *VestackIamUserGroupPolicyAttachmentService {
	return &VestackIamUserGroupPolicyAttachmentService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackIamUserGroupPolicyAttachmentService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamUserGroupPolicyAttachmentService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithSimpleQuery(condition, func(m map[string]interface{}) ([]interface{}, error) {
		action := "ListAttachedUserGroupPolicies"
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
		results, err = bp.ObtainSdkValue("Result.AttachedPolicyMetadata", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.AttachedPolicyMetadata is not slice")
		}
		return data, err
	})
}

func (s *VestackIamUserGroupPolicyAttachmentService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		tempData = map[string]interface{}{}
		results  []interface{}
		ok       bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	req := map[string]interface{}{
		"UserGroupName": ids[0],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if tempData, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
		pName, ok := tempData["PolicyName"].(string)
		if !ok {
			return data, errors.New("PolicyName is not string")
		}
		if pName == ids[1] {
			data = tempData
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("iam_user_group_policy_attachment %s not exist ", id)
	}
	return data, err
}

func (s *VestackIamUserGroupPolicyAttachmentService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *VestackIamUserGroupPolicyAttachmentService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AttachUserGroupPolicy",
			ConvertMode: bp.RequestConvertAll,
			Convert:     map[string]bp.RequestConvert{},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				userGroupName, ok := d.Get("user_group_name").(string)
				if !ok {
					return errors.New("user_group_name is not string")
				}
				policyName, ok := d.Get("policy_name").(string)
				if !ok {
					return errors.New("policy_name is not string")
				}
				d.SetId(fmt.Sprintf("%s:%s", userGroupName, policyName))
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				userGroupName, ok := d.Get("user_group_name").(string)
				if !ok {
					return ""
				}
				return userGroupName
			},
		},
	}
	return []bp.Callback{callback}
}

func (VestackIamUserGroupPolicyAttachmentService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackIamUserGroupPolicyAttachmentService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackIamUserGroupPolicyAttachmentService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	ids := strings.Split(resourceData.Id(), ":")
	policyType, ok := resourceData.Get("policy_type").(string)
	if !ok {
		return []bp.Callback{{Err: errors.New("policy_type is not string")}}
	}
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DetachUserGroupPolicy",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"UserGroupName": ids[0],
				"PolicyName":    ids[1],
				"PolicyType":    policyType,
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackIamUserGroupPolicyAttachmentService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
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
		NameField:    "UserGroupName",
		CollectField: "policies",
	}
}

func (s *VestackIamUserGroupPolicyAttachmentService) ReadResourceId(id string) string {
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
