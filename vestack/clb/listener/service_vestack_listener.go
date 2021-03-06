package listener

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ve "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/clb"
)

type VestackListenerService struct {
	Client     *ve.SdkClient
	Dispatcher *ve.Dispatcher
}

func NewListenerService(c *ve.SdkClient) *VestackListenerService {
	return &VestackListenerService{
		Client:     c,
		Dispatcher: &ve.Dispatcher{},
	}
}

func (s *VestackListenerService) GetClient() *ve.SdkClient {
	return s.Client
}

func (s *VestackListenerService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return ve.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		clbClient := s.Client.ClbClient
		action := "DescribeListeners"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = clbClient.DescribeListenersCommon(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = clbClient.DescribeListenersCommon(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = ve.ObtainSdkValue("Result.Listeners", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Listeners is not Slice")
		}
		return data, err
	})
}

func (s *VestackListenerService) ReadResource(resourceData *schema.ResourceData, listenerId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if listenerId == "" {
		listenerId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"ListenerIds.1": listenerId,
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
		return data, fmt.Errorf("Listener %s not exist ", listenerId)
	}

	clbClient := s.Client.ClbClient

	listenerAttr, err := clbClient.DescribeListenerAttributesCommon(&map[string]interface{}{
		"ListenerId": listenerId,
	})
	if err != nil {
		return nil, err
	}
	timeout, err := ve.ObtainSdkValue("Result.EstablishedTimeout", *listenerAttr)
	if err != nil {
		return nil, err
	}
	desc, err := ve.ObtainSdkValue("Result.Description", *listenerAttr)
	if err != nil {
		return nil, err
	}
	loadBalancerId, err := ve.ObtainSdkValue("Result.LoadBalancerId", *listenerAttr)
	if err != nil {
		return nil, err
	}
	scheduler, err := ve.ObtainSdkValue("Result.Scheduler", *listenerAttr)
	if err != nil {
		return nil, err
	}
	_ = resourceData.Set("established_timeout", int(timeout.(float64)))
	_ = resourceData.Set("description", desc.(string))
	_ = resourceData.Set("load_balancer_id", loadBalancerId.(string))
	_ = resourceData.Set("scheduler", scheduler.(string))

	return data, err
}

func (s *VestackListenerService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = ve.ObtainSdkValue("Status", demo)
			if err != nil {
				return nil, "", err
			}
			//?????? ???????????????????????????????????? ????????????????????????
			return demo, status.(string), err
		},
	}

}

func (VestackListenerService) WithResourceResponseHandlers(listener map[string]interface{}) []ve.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]ve.ResponseConvert, error) {
		return listener, nil, nil
	}
	return []ve.ResourceResponseHandler{handler}

}

func (s *VestackListenerService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []ve.Callback {
	callback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "CreateListener",
			ConvertMode: ve.RequestConvertAll,
			Convert: map[string]ve.RequestConvert{
				"acl_ids": {
					ConvertType: ve.ConvertWithN,
				},
				"health_check": {
					ConvertType: ve.ConvertListUnique,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//??????listener
				return s.Client.ClbClient.CreateListenerCommon(call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *ve.SdkClient, resp *map[string]interface{}, call ve.SdkCall) error {
				//?????? ???????????? ??????????????????????????? ???????????????
				id, _ := ve.ObtainSdkValue("Result.ListenerId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &ve.StateRefresh{
				Target:  []string{"Active", "Disabled"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			ExtraRefresh: map[ve.ResourceService]*ve.StateRefresh{
				clb.NewClbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("load_balancer_id").(string),
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return resourceData.Get("load_balancer_id").(string)
			},
		},
	}
	return []ve.Callback{callback}

}

func (s *VestackListenerService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []ve.Callback {
	clbId, err := s.queryLoadBalancerId(resourceData.Id())
	if err != nil {
		return []ve.Callback{{
			Err: err,
		}}
	}

	callback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "ModifyListenerAttributes",
			ConvertMode: ve.RequestConvertAll,
			Convert: map[string]ve.RequestConvert{
				"acl_ids": {
					ConvertType: ve.ConvertWithN,
				},
				"health_check": {
					ConvertType: ve.ConvertListUnique,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (bool, error) {
				(*call.SdkParam)["ListenerId"] = d.Id()
				aclStatus := d.Get("acl_status")
				if aclStatus, ok := aclStatus.(string); ok && aclStatus == "on" {
					(*call.SdkParam)["AclStatus"] = d.Get("acl_status").(string)
					(*call.SdkParam)["AclType"] = d.Get("acl_type").(string)
					if !d.HasChange("acl_ids") {
						if m, ok := d.Get("acl_ids").(*schema.Set); ok {
							aclIds := m.List()
							for i, aclId := range aclIds {
								k := fmt.Sprintf("AclIds.%d", i+1)
								(*call.SdkParam)[k] = aclId
							}
						}
					}
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//?????? listener ??????
				return s.Client.ClbClient.ModifyListenerAttributesCommon(call.SdkParam)
			},
			Refresh: &ve.StateRefresh{
				Target:  []string{"Active", "Disabled"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			ExtraRefresh: map[ve.ResourceService]*ve.StateRefresh{
				clb.NewClbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: clbId,
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return clbId
			},
		},
	}
	return []ve.Callback{callback}
}

func (s *VestackListenerService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []ve.Callback {
	clbId, err := s.queryLoadBalancerId(resourceData.Id())
	if err != nil {
		return []ve.Callback{{
			Err: err,
		}}
	}

	callback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "DeleteListener",
			ConvertMode: ve.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"ListenerId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//?????? Listener
				return s.Client.ClbClient.DeleteListenerCommon(call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall, baseErr error) error {
				//?????????????????????
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if ve.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading listener on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			ExtraRefresh: map[ve.ResourceService]*ve.StateRefresh{
				clb.NewClbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: clbId,
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return clbId
			},
		},
	}
	return []ve.Callback{callback}
}

func (s *VestackListenerService) DatasourceResources(*schema.ResourceData, *schema.Resource) ve.DataSourceInfo {
	return ve.DataSourceInfo{
		RequestConverts: map[string]ve.RequestConvert{
			"ids": {
				TargetField: "ListenerIds",
				ConvertType: ve.ConvertWithN,
			},
		},
		NameField:    "ListenerName",
		IdField:      "ListenerId",
		CollectField: "listeners",
		ResponseConverts: map[string]ve.ResponseConvert{
			"ListenerId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"HealthCheck.Enabled": {
				TargetField: "health_check_enabled",
			},
			"HealthCheck.Interval": {
				TargetField: "health_check_interval",
			},
			"HealthCheck.HealthyThreshold": {
				TargetField: "health_check_healthy_threshold",
			},
			"HealthCheck.UnHealthyThreshold": {
				TargetField: "health_check_un_healthy_threshold",
			},
			"HealthCheck.Timeout": {
				TargetField: "health_check_timeout",
			},
			"HealthCheck.Method": {
				TargetField: "health_check_method",
			},
			"HealthCheck.Uri": {
				TargetField: "health_check_uri",
			},
			"HealthCheck.Domain": {
				TargetField: "health_check_domain",
			},
			"HealthCheck.HttpCode": {
				TargetField: "health_check_http_code",
			},
		},
	}
}

func (s *VestackListenerService) ReadResourceId(id string) string {
	return id
}

func (s *VestackListenerService) queryLoadBalancerId(listenerId string) (string, error) {
	if listenerId == "" {
		return "", fmt.Errorf("listener ID cannot be empty")
	}

	// ?????? LoadBalancerId
	serverGroupResp, err := s.Client.ClbClient.DescribeListenerAttributesCommon(&map[string]interface{}{
		"ListenerId": listenerId,
	})
	if err != nil {
		return "", err
	}
	clbId, err := ve.ObtainSdkValue("Result.LoadBalancerId", *serverGroupResp)
	if err != nil {
		return "", err
	}
	return clbId.(string), nil
}
