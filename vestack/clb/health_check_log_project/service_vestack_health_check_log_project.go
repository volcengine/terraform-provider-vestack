package health_check_log_project

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackHealthCheckLogProjectService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewHealthCheckLogProjectService(c *bp.SdkClient) *VestackHealthCheckLogProjectService {
	return &VestackHealthCheckLogProjectService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackHealthCheckLogProjectService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackHealthCheckLogProjectService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	action := "DescribeHealthCheckLogProjectAttributes"

	bytes, _ := json.Marshal(m)
	logger.Debug(logger.ReqFormat, action, string(bytes))

	resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
	if err != nil {
		return data, err
	}

	respBytes, _ := json.Marshal(resp)
	logger.Debug(logger.RespFormat, action, string(respBytes))

	results, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return data, err
	}

	if results == nil {
		results = []interface{}{}
	}
	// 单个结果，包装成列表返回
	if resultMap, ok := results.(map[string]interface{}); ok {
		data = []interface{}{resultMap}
	}
	return data, err
}

func (s *VestackHealthCheckLogProjectService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		resp *map[string]interface{}
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	action := "DescribeHealthCheckLogProjectAttributes"
	logger.Debug(logger.ReqFormat, action, id)

	resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
	if err != nil {
		return data, err
	}

	respBytes, _ := json.Marshal(resp)
	logger.Debug(logger.RespFormat, action, string(respBytes))

	// 从响应中提取Result
	results, err := bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return data, err
	}

	if results == nil {
		return data, fmt.Errorf("health_check_log_project %s not exist", id)
	}

	// 处理Result - 它应该是一个map包含LogProjectId
	if resultMap, ok := results.(map[string]interface{}); ok {
		data = make(map[string]interface{})
		data["log_project_id"] = resultMap["LogProjectId"]
		data["id"] = resultMap["LogProjectId"]
	} else {
		return data, fmt.Errorf("invalid result format for health_check_log_project %s", id)
	}

	if data["log_project_id"] == nil || data["log_project_id"].(string) == "" {
		return data, fmt.Errorf("health_check_log_project %s not exist", id)
	}
	return data, err
}

func (s *VestackHealthCheckLogProjectService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateHealthCheckLogProject",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam:    &map[string]interface{}{},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.LogProjectId", *resp)
				d.SetId(id.(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackHealthCheckLogProjectService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteHealthCheckLogProject",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam:    &map[string]interface{}{},
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

func (s *VestackHealthCheckLogProjectService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackHealthCheckLogProjectService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				d map[string]interface{}
			)
			d, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			return d, "Available", err
		},
	}
}

func (s *VestackHealthCheckLogProjectService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackHealthCheckLogProjectService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		NameField:    "LogProjectId",
		IdField:      "LogProjectId",
		CollectField: "health_check_log_projects",
		ResponseConverts: map[string]bp.ResponseConvert{
			"LogProjectId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *VestackHealthCheckLogProjectService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "clb",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
