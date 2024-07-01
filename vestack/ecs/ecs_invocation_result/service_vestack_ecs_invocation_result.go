package ecs_invocation_result

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackEcsInvocationResultService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewEcsInvocationResultService(c *bp.SdkClient) *VestackEcsInvocationResultService {
	return &VestackEcsInvocationResultService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackEcsInvocationResultService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackEcsInvocationResultService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeInvocationResults"
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
		results, err = bp.ObtainSdkValue("Result.InvocationResults", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.InvocationResults is not Slice")
		}

		return data, err
	})
}

func (s *VestackEcsInvocationResultService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	return data, err
}

func (s *VestackEcsInvocationResultService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (VestackEcsInvocationResultService) WithResourceResponseHandlers(data map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return data, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackEcsInvocationResultService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackEcsInvocationResultService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackEcsInvocationResultService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackEcsInvocationResultService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"invocation_result_status": {
				TargetField: "InvocationResultStatus",
				Convert: func(data *schema.ResourceData, i interface{}) interface{} {
					var status string
					statusSet, ok := data.GetOk("invocation_result_status")
					if !ok {
						return status
					}
					statusList := statusSet.(*schema.Set).List()
					statusArr := make([]string, 0)
					for _, value := range statusList {
						statusArr = append(statusArr, value.(string))
					}
					status = strings.Join(statusArr, ",")
					return status
				},
			},
		},
		IdField:      "InvocationResultId",
		CollectField: "invocation_results",
		ResponseConverts: map[string]bp.ResponseConvert{
			"InvocationResultId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *VestackEcsInvocationResultService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "ecs",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
