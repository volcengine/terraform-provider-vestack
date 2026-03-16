package listener_health

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackListenerHealthService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewListenerHealthService(c *bp.SdkClient) *VestackListenerHealthService {
	return &VestackListenerHealthService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackListenerHealthService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackListenerHealthService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp           *map[string]interface{}
		unHealthyCount interface{}
		status         interface{}
	)

	results, err := bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) (data []interface{}, err error) {
		action := "DescribeListenerHealth"
		logger.Debug(logger.ReqFormat, action, m)
		resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &m)
		if err != nil {
			return nil, err
		}
		respBytes, _ := json.Marshal(resp)
		logger.Debug(logger.RespFormat, action, m, string(respBytes))

		// 获取后端服务器健康状态结果
		pageResults, err := bp.ObtainSdkValue("Result.Results", *resp)
		if err != nil {
			return nil, err
		}

		// 同时获取汇总信息
		unHealthyCount, _ = bp.ObtainSdkValue("Result.UnHealthyCount", *resp)
		status, _ = bp.ObtainSdkValue("Result.Status", *resp)

		if pageResults == nil {
			return []interface{}{}, nil
		}
		if slice, ok := pageResults.([]interface{}); ok {
			return slice, nil
		}
		return nil, errors.New("Result.Results is not Slice")
	})

	if err != nil {
		return data, err
	}

	// 构造返回结构，包含汇总信息和服务器列表
	healthInfo := map[string]interface{}{
		"UnHealthyCount": unHealthyCount,
		"ListenerStatus": status,
		"Results":        results,
	}

	return []interface{}{healthInfo}, nil
}

func (s *VestackListenerHealthService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	return data, nil
}

func (s *VestackListenerHealthService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *VestackListenerHealthService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{}
	return []bp.Callback{callback}
}

func (VestackListenerHealthService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackListenerHealthService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{}
	return []bp.Callback{callback}
}

func (s *VestackListenerHealthService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{}
	return []bp.Callback{callback}
}

func (s *VestackListenerHealthService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{},
		CollectField:    "health_info",
	}
}

func (s *VestackListenerHealthService) ReadResourceId(id string) string {
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
