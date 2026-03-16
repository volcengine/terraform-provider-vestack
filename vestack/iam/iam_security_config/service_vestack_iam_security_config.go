package iam_security_config

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamSecurityConfigService struct {
	Client *bp.SdkClient
}

func NewIamSecurityConfigService(c *bp.SdkClient) *VestackIamSecurityConfigService {
	return &VestackIamSecurityConfigService{
		Client: c,
	}
}

func (s *VestackIamSecurityConfigService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamSecurityConfigService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		result interface{}
	)
	action := "GetSecurityConfig"
	logger.Debug(logger.ReqFormat, action, m)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &m)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp)
	result, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return data, err
	}
	if dataMap, ok := result.(map[string]interface{}); ok {
		// API returns UserID, SafeAuthType, etc. We need to inject UserName because API doesn't return it in Result (only in request)
		// Wait, ReadResources is usually for List. GetSecurityConfig is singular.
		// If m contains UserName, we can inject it back.
		if userName, ok := m["UserName"]; ok {
			dataMap["UserName"] = userName
		}
		data = append(data, dataMap)
	} else {
		return data, errors.New("Value is not map ")
	}

	return data, nil
}

func (s *VestackIamSecurityConfigService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		result interface{}
		ok     bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	condition := map[string]interface{}{"UserName": id}
	action := "GetSecurityConfig"
	logger.Debug(logger.ReqFormat, action, condition)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp)
	result, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return data, err
	}
	if data, ok = result.(map[string]interface{}); !ok {
		return data, errors.New("Value is not map ")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("security config for user %s not exist ", id)
	}

	return data, err
}

func (s *VestackIamSecurityConfigService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (VestackIamSecurityConfigService) WithResourceResponseHandlers(v map[string]interface{}) []bp.ResourceResponseHandler {
	return []bp.ResourceResponseHandler{}
}

func (s *VestackIamSecurityConfigService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "SetSecurityConfig",
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				if v, ok := d.Get("user_name").(string); ok && v != "" {
					d.SetId(v)
				}
				return nil
			},
		},
	}

	return []bp.Callback{callback}
}

func (s *VestackIamSecurityConfigService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:         "SetSecurityConfig",
			ConvertMode:    bp.RequestConvertAll,
			RequestIdField: "UserName",
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackIamSecurityConfigService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	// No delete API available.
	return []bp.Callback{}
}

func (s *VestackIamSecurityConfigService) DatasourceResources(resourceData *schema.ResourceData, r *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"user_name": {
				TargetField: "UserName",
			},
		},
		NameField:    "UserName",
		IdField:      "UserName",
		CollectField: "security_configs",
		ResponseConverts: map[string]bp.ResponseConvert{
			"UserID": {
				TargetField: "user_id",
			},
			"SafeAuthType": {
				TargetField: "safe_auth_type",
			},
			"SafeAuthExemptDuration": {
				TargetField: "safe_auth_exempt_duration",
			},
			"SafeAuthClose": {
				TargetField: "safe_auth_close",
			},
		},
	}
}

func (s *VestackIamSecurityConfigService) ReadResourceId(id string) string {
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
