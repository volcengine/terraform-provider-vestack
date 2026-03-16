package iam_login_profile

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamLoginProfileService struct {
	Client *bp.SdkClient
}

func NewIamLoginProfileService(c *bp.SdkClient) *VestackIamLoginProfileService {
	return &VestackIamLoginProfileService{
		Client: c,
	}
}

func (s *VestackIamLoginProfileService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamLoginProfileService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		result interface{}
	)
	action := "GetLoginProfile"
	logger.Debug(logger.ReqFormat, action, m)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &m)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp)
	result, err = bp.ObtainSdkValue("Result.LoginProfile", *resp)
	if err != nil {
		return data, err
	}
	if dataMap, ok := result.(map[string]interface{}); ok {
		delete(dataMap, "Password")
		data = append(data, dataMap)
	} else {
		return data, errors.New("Value is not map ")
	}

	return data, nil
}

func (s *VestackIamLoginProfileService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		result interface{}
		ok     bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	condition := map[string]interface{}{"UserName": id}
	action := "GetLoginProfile"
	logger.Debug(logger.ReqFormat, action, condition)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp)
	result, err = bp.ObtainSdkValue("Result.LoginProfile", *resp)
	if err != nil {
		return data, err
	}
	if data, ok = result.(map[string]interface{}); !ok {
		return data, errors.New("Value is not map ")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("login profile %s not exist ", id)
	}

	return data, err
}

func (s *VestackIamLoginProfileService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (VestackIamLoginProfileService) WithResourceResponseHandlers(v map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		delete(v, "Password")
		return v, map[string]bp.ResponseConvert{}, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackIamLoginProfileService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateLoginProfile",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"user_name": {
					ForceGet: true,
				},
				"password": {
					ForceGet: true,
				},
				"login_allowed": {
					TargetField: "LoginAllowed",
				},
				"password_reset_required": {
					TargetField: "PasswordResetRequired",
				},
				"safe_auth_flag": {
					TargetField: "SafeAuthFlag",
				},
				"safe_auth_type": {
					TargetField: "SafeAuthType",
				},
				"safe_auth_exempt_required": {
					TargetField: "SafeAuthExemptRequired",
				},
				"safe_auth_exempt_unit": {
					TargetField: "SafeAuthExemptUnit",
				},
				"safe_auth_exempt_duration": {
					TargetField: "SafeAuthExemptDuration",
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, resp)
				time.Sleep(5 * time.Second)
				if v, ok := d.Get("user_name").(string); ok && v != "" {
					d.SetId(v)
				}
				return nil
			},
		},
	}

	return []bp.Callback{callback}
}

func (s *VestackIamLoginProfileService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateLoginProfile",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"password": {
					ForceGet:    true,
					TargetField: "Password",
				},
				"login_allowed": {
					ForceGet:    true,
					TargetField: "LoginAllowed",
				},
				"password_reset_required": {
					ForceGet:    true,
					TargetField: "PasswordResetRequired",
				},
				"safe_auth_flag": {
					ForceGet:    true,
					TargetField: "SafeAuthFlag",
				},
				"safe_auth_type": {
					ForceGet:    true,
					TargetField: "SafeAuthType",
				},
				"safe_auth_exempt_required": {
					ForceGet:    true,
					TargetField: "SafeAuthExemptRequired",
				},
				"safe_auth_exempt_unit": {
					ForceGet:    true,
					TargetField: "SafeAuthExemptUnit",
				},
				"safe_auth_exempt_duration": {
					ForceGet:    true,
					TargetField: "SafeAuthExemptDuration",
				},
			},
			RequestIdField: "UserName",
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				time.Sleep(5 * time.Second)
				return resp, err
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackIamLoginProfileService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:         "DeleteLoginProfile",
			ConvertMode:    bp.RequestConvertIgnore,
			RequestIdField: "UserName",
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}

	return []bp.Callback{callback}
}

func (s *VestackIamLoginProfileService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"user_name": {
				TargetField: "UserName",
			},
		},
		NameField:    "UserName",
		IdField:      "UserName",
		CollectField: "login_profiles",
		ResponseConverts: map[string]bp.ResponseConvert{
			"SafeAuthFlag": {
				TargetField: "safe_auth_flag",
			},
			"SafeAuthType": {
				TargetField: "safe_auth_type",
			},
			"SafeAuthExemptRequired": {
				TargetField: "safe_auth_exempt_required",
			},
			"SafeAuthExemptUnit": {
				TargetField: "safe_auth_exempt_unit",
			},
			"SafeAuthExemptDuration": {
				TargetField: "safe_auth_exempt_duration",
			},
		},
	}
}

func (s *VestackIamLoginProfileService) ReadResourceId(id string) string {
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
