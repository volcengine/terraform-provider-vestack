package ecs_invocation

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackEcsInvocationService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewEcsInvocationService(c *bp.SdkClient) *VestackEcsInvocationService {
	return &VestackEcsInvocationService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackEcsInvocationService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackEcsInvocationService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeInvocations"
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
		results, err = bp.ObtainSdkValue("Result.Invocations", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Invocations is not Slice")
		}

		for _, v := range data {
			instanceIds := make([]string, 0)
			invocation, ok := v.(map[string]interface{})
			if !ok {
				return data, fmt.Errorf(" Invocation is not map ")
			}
			action := "DescribeInvocationInstances"
			req := map[string]interface{}{
				"InvocationId": invocation["InvocationId"],
			}
			logger.Debug(logger.ReqFormat, action, req)
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
			logger.Debug(logger.RespFormat, action, req, resp)
			results, err := bp.ObtainSdkValue("Result.InvocationInstances", *resp)
			if err != nil {
				return data, err
			}
			if results == nil {
				results = []interface{}{}
			}
			instances, ok := results.([]interface{})
			if !ok {
				return data, errors.New("Result.InvocationInstances is not Slice")
			}
			if len(instances) == 0 {
				return data, fmt.Errorf("invocation %s does not contain any instances", invocation["InvocationId"])
			}
			for _, v1 := range instances {
				instance, ok := v1.(map[string]interface{})
				if !ok {
					return data, fmt.Errorf(" invocation instance is not map ")
				}
				instanceIds = append(instanceIds, instance["InstanceId"].(string))
			}
			invocation["InstanceIds"] = instanceIds
		}

		return data, err
	})
}

func (s *VestackEcsInvocationService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"InvocationId": id,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("ecs invocation %s is not exist ", id)
	}

	// 处理 launch_time、recurrence_end_time 传参与查询结果不一致的问题
	if mode := resourceData.Get("repeat_mode"); mode.(string) != "Once" {
		layout := "2006-01-02T15:04:05Z"
		launchTimeExpr, exist1 := resourceData.GetOkExists("launch_time")
		endTimeExpr, exist2 := resourceData.GetOkExists("recurrence_end_time")
		if exist1 && launchTimeExpr.(string) != "" {
			launchTime, err := ParseUTCTime(launchTimeExpr.(string))
			if err != nil {
				return data, err
			}
			lt := launchTime.Format(layout)
			if lt == data["LaunchTime"].(string) {
				data["LaunchTime"] = launchTimeExpr
			}
		}
		if exist2 && endTimeExpr.(string) != "" {
			endTime, err := ParseUTCTime(endTimeExpr.(string))
			if err != nil {
				return data, err
			}
			et := endTime.Format(layout)
			if et == data["RecurrenceEndTime"].(string) {
				data["RecurrenceEndTime"] = endTimeExpr
			}
		}
	}

	return data, err
}

func ParseUTCTime(timeExpr string) (time.Time, error) {
	timeWithoutSecond, err := ParseUTCTimeWithoutSecond(timeExpr)
	if err != nil {
		timeWithSecond, err := ParseUTCTimeWithSecond(timeExpr)
		if err != nil {
			return time.Time{}, err
		} else {
			return timeWithSecond, nil
		}
	} else {
		return timeWithoutSecond, nil
	}
}

func ParseUTCTimeWithoutSecond(timeExpr string) (time.Time, error) {
	t, err := time.Parse("2006-01-02T15:04Z", timeExpr)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse time failed, error: %v, time expr: %v", err, timeExpr)
	}

	return t, nil
}

func ParseUTCTimeWithSecond(timeExpr string) (time.Time, error) {
	t, err := time.Parse("2006-01-02T15:04:05Z", timeExpr)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse time failed, error: %v, time expr: %v", err, timeExpr)
	}

	t = t.Add(time.Duration(t.Second()) * time.Second * -1)

	return t, nil
}

func (s *VestackEcsInvocationService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			//no failed status.
			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("InvocationStatus", demo)
			if err != nil {
				return nil, "", err
			}
			return demo, status.(string), err
		},
	}
}

func (VestackEcsInvocationService) WithResourceResponseHandlers(invocation map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return invocation, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackEcsInvocationService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "InvokeCommand",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeDefault,
			Convert: map[string]bp.RequestConvert{
				"instance_ids": {
					TargetField: "InstanceIds",
					ConvertType: bp.ConvertWithN,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.InvocationId", *resp)
				d.SetId(id.(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackEcsInvocationService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackEcsInvocationService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "StopInvocation",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeDefault,
			SdkParam: &map[string]interface{}{
				"InvocationId": resourceData.Id(),
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				status := d.Get("invocation_status")
				mode := d.Get("repeat_mode")
				if mode.(string) == "Once" || (status.(string) != "Pending" && status.(string) != "Scheduled") {
					return false, nil
				} else {
					(*call.SdkParam)["InvocationId"] = d.Id()
					return true, nil
				}
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Stopped"},
				Timeout: resourceData.Timeout(schema.TimeoutDelete),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackEcsInvocationService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"invocation_status": {
				TargetField: "InvocationStatus",
				Convert: func(data *schema.ResourceData, i interface{}) interface{} {
					var status string
					statusSet, ok := data.GetOk("invocation_status")
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
		NameField:    "InvocationName",
		IdField:      "InvocationId",
		CollectField: "invocations",
		ResponseConverts: map[string]bp.ResponseConvert{
			"InvocationId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *VestackEcsInvocationService) ReadResourceId(id string) string {
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
