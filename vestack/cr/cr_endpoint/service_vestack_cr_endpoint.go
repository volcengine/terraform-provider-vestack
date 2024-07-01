package cr_endpoint

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackCrEndpointService struct {
	Client *bp.SdkClient
}

func NewCrEndpointService(c *bp.SdkClient) *VestackCrEndpointService {
	return &VestackCrEndpointService{
		Client: c,
	}
}

func (s *VestackCrEndpointService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackCrEndpointService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	action := "GetPublicEndpoint"

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

	logger.Debug(logger.RespFormat, action, resp)
	results, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return data, err
	}
	if results == nil {
		return data, fmt.Errorf("GetPublicEndpoint return an empty result")
	}

	registry, err := bp.ObtainSdkValue("Result.Registry", *resp)
	if err != nil {
		return data, err
	}
	enabled, err := bp.ObtainSdkValue("Result.Enabled", *resp)
	if err != nil {
		return data, err
	}
	status, err := bp.ObtainSdkValue("Result.Status", *resp)
	if err != nil {
		return data, err
	}
	endpoint := map[string]interface{}{
		"Registry": registry,
		"Enabled":  enabled,
		"Status":   status,
	}

	return []interface{}{endpoint}, err
}

func (s *VestackCrEndpointService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)

	registry := resourceData.Get("registry").(string)
	req := map[string]interface{}{
		"Registry": registry,
	}

	results, err = s.ReadResources(req)

	if err != nil {
		return data, err
	}

	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("GetPublicEndpoint value is not a map")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("cr endpoint %s is not exist", id)
	}
	return data, err
}

func (s *VestackCrEndpointService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,

		Refresh: func() (result interface{}, state string, err error) {
			var (
				demo   map[string]interface{}
				status interface{}
			)
			failedStatus := []string{"Failed"}

			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			logger.DebugInfo("Refresh CrEndpoint status resp:%v", demo)

			status, err = bp.ObtainSdkValue("Status", demo)
			if err != nil {
				return nil, "", err
			}

			for _, v := range failedStatus {
				if v == status.(string) {
					return nil, "", fmt.Errorf("CrEndpoint status error,status %s", status.(string))
				}
			}

			return demo, status.(string), err
		},
	}
}

func (VestackCrEndpointService) WithResourceResponseHandlers(instance map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return instance, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackCrEndpointService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	target := "Disabled"
	enabled := resourceData.Get("enabled").(bool)
	if enabled {
		target = "Enabled"
	}
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdatePublicEndpoint",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Registry"] = d.Get("registry")
				(*call.SdkParam)["Enabled"] = d.Get("enabled")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				registry := d.Get("registry").(string)
				id := "endpoint:" + registry
				d.SetId(id)
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{target},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackCrEndpointService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	target := "Disabled"
	enabled := resourceData.Get("enabled").(bool)
	if enabled {
		target = "Enabled"
	}
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdatePublicEndpoint",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Registry"] = d.Get("registry")
				(*call.SdkParam)["Enabled"] = d.Get("enabled")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{target},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackCrEndpointService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdatePublicEndpoint",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Registry"] = d.Get("registry")
				(*call.SdkParam)["Enabled"] = false
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Disabled"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackCrEndpointService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType:  bp.ContentTypeJson,
		CollectField: "endpoints",
	}
}

func (s *VestackCrEndpointService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "cr",
		Version:     "2022-05-12",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
