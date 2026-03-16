package ebs_max_extra_performance

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackEbsMaxExtraPerformanceService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewEbsMaxExtraPerformanceService(c *bp.SdkClient) *VestackEbsMaxExtraPerformanceService {
	return &VestackEbsMaxExtraPerformanceService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackEbsMaxExtraPerformanceService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackEbsMaxExtraPerformanceService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithSimpleQuery(condition, func(m map[string]interface{}) ([]interface{}, error) {
		action := "DescribeMaxExtraPerformanceCanPurchase"

		bytes, _ := json.Marshal(condition)
		logger.Debug(logger.ReqFormat, action, string(bytes))
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
		respBytes, _ := json.Marshal(resp)
		logger.Debug(logger.RespFormat, action, condition, string(respBytes))

		result, err := bp.ObtainSdkValue("Result", *resp)
		if err != nil {
			return data, err
		}
		if resultMap, ok := result.(map[string]interface{}); ok {
			results = []interface{}{resultMap}
		}

		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result is not Slice ")
		}
		return data, err
	})
}

func (s *VestackEbsMaxExtraPerformanceService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	return data, err
}

func (s *VestackEbsMaxExtraPerformanceService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (VestackEbsMaxExtraPerformanceService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackEbsMaxExtraPerformanceService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackEbsMaxExtraPerformanceService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackEbsMaxExtraPerformanceService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackEbsMaxExtraPerformanceService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		CollectField: "performances",
		ResponseConverts: map[string]bp.ResponseConvert{
			"IOPS": {
				TargetField: "iops",
			},
		},
	}
}

func (s *VestackEbsMaxExtraPerformanceService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "storage_ebs",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
