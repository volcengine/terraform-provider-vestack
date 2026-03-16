package traffic_mirror_filter_rule

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

type VestackTrafficMirrorFilterRuleService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewTrafficMirrorFilterRuleService(c *bp.SdkClient) *VestackTrafficMirrorFilterRuleService {
	return &VestackTrafficMirrorFilterRuleService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackTrafficMirrorFilterRuleService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackTrafficMirrorFilterRuleService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		filters []interface{}
		ok      bool
	)
	return bp.WithNextTokenQuery(m, "MaxResults", "NextToken", 100, nil, func(condition map[string]interface{}) (data []interface{}, next string, err error) {
		action := "DescribeTrafficMirrorFilters"

		bytes, _ := json.Marshal(condition)
		logger.Debug(logger.ReqFormat, action, string(bytes))
		if condition == nil {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, next, err
			}
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, next, err
			}
		}
		respBytes, _ := json.Marshal(resp)
		logger.Debug(logger.RespFormat, action, condition, string(respBytes))
		results, err = bp.ObtainSdkValue("Result.TrafficMirrorFilters", *resp)
		if err != nil {
			return data, next, err
		}
		nextToken, err := bp.ObtainSdkValue("Result.NextToken", *resp)
		if err != nil {
			return data, next, err
		}
		next = nextToken.(string)
		if results == nil {
			results = []interface{}{}
		}
		if filters, ok = results.([]interface{}); !ok {
			return data, next, errors.New("Result.TrafficMirrorFilters is not Slice")
		}

		for _, filter := range filters {
			filterMap, ok := filter.(map[string]interface{})
			if !ok {
				return data, next, errors.New("TrafficMirrorFilter is not map")
			}
			if v, ok := filterMap["IngressFilterRules"]; ok {
				if ingressRules, ok := v.([]interface{}); ok {
					data = append(data, ingressRules...)
				}
			}
			if v, ok := filterMap["EgressFilterRules"]; ok {
				if egressRules, ok := v.([]interface{}); ok {
					data = append(data, egressRules...)
				}
			}
		}

		return data, next, err
	})
}

func (s *VestackTrafficMirrorFilterRuleService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return data, fmt.Errorf("invald traffic mirror filter rule id: %s", id)
	}

	req := map[string]interface{}{
		"TrafficMirrorFilterIds.1": ids[0],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		rule := make(map[string]interface{})
		if rule, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
		if rule["TrafficMirrorFilterRuleId"] == ids[1] {
			data = rule
			break
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("traffic_mirror_filter_rule %s not exist ", id)
	}
	return data, err
}

func (s *VestackTrafficMirrorFilterRuleService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				d          map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "Failed")
			d, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", d)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("traffic_mirror_filter_rule status error, status: %s", status.(string))
				}
			}
			return d, status.(string), err
		},
	}
}

func (VestackTrafficMirrorFilterRuleService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackTrafficMirrorFilterRuleService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateTrafficMirrorFilterRule",
			ConvertMode: bp.RequestConvertAll,
			Convert:     map[string]bp.RequestConvert{},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.TrafficMirrorFilterRuleId", *resp)
				filterId := d.Get("traffic_mirror_filter_id").(string)
				d.SetId(filterId + ":" + id.(string))
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("traffic_mirror_filter_id").(string)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackTrafficMirrorFilterRuleService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyTrafficMirrorFilterRuleAttributes",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"traffic_direction": {
					TargetField: "TrafficDirection",
				},
				"priority": {
					TargetField: "Priority",
				},
				"policy": {
					TargetField: "Policy",
				},
				"protocol": {
					TargetField: "Protocol",
				},
				"source_cidr_block": {
					TargetField: "SourceCidrBlock",
				},
				"source_port_range": {
					TargetField: "SourcePortRange",
				},
				"destination_cidr_block": {
					TargetField: "DestinationCidrBlock",
				},
				"destination_port_range": {
					TargetField: "DestinationPortRange",
				},
				"description": {
					TargetField: "Description",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) > 0 {
					ids := strings.Split(d.Id(), ":")
					if len(ids) != 2 {
						return false, fmt.Errorf("invalid traffic mirror filter rule id: %s", d.Id())
					}
					(*call.SdkParam)["TrafficMirrorFilterRuleId"] = ids[1]
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("traffic_mirror_filter_id").(string)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackTrafficMirrorFilterRuleService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteTrafficMirrorFilterRule",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ids := strings.Split(d.Id(), ":")
				if len(ids) != 2 {
					return false, fmt.Errorf("invalid traffic mirror filter rule id: %s", d.Id())
				}
				(*call.SdkParam)["TrafficMirrorFilterRuleId"] = ids[1]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("traffic_mirror_filter_id").(string)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(5*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading traffic mirror filter rule on delete %q, %w", d.Id(), callErr))
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
	return []bp.Callback{callback}
}

func (s *VestackTrafficMirrorFilterRuleService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"traffic_mirror_filter_ids": {
				TargetField: "TrafficMirrorFilterIds",
				ConvertType: bp.ConvertWithN,
			},
			"traffic_mirror_filter_names": {
				TargetField: "TrafficMirrorFilterNames",
				ConvertType: bp.ConvertWithN,
			},
			"tags": {
				TargetField: "TagFilters",
				ConvertType: bp.ConvertListN,
				NextLevelConvert: map[string]bp.RequestConvert{
					"value": {
						TargetField: "Values.1",
					},
				},
			},
		},
		IdField:      "TrafficMirrorFilterRuleId",
		CollectField: "traffic_mirror_filter_rules",
		ResponseConverts: map[string]bp.ResponseConvert{
			"TrafficMirrorFilterRuleId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *VestackTrafficMirrorFilterRuleService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vpc",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
