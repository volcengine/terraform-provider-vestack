package listener

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/clb"
)

type VestackListenerService struct {
	Client *bp.SdkClient
}

func NewListenerService(c *bp.SdkClient) *VestackListenerService {
	return &VestackListenerService{
		Client: c,
	}
}

func (s *VestackListenerService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackListenerService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		action := "DescribeListeners"
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

		results, err = bp.ObtainSdkValue("Result.Listeners", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Listeners is not Slice")
		}
		data, err = removeSystemTags(data)
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

	action := "DescribeListenerAttributes"
	listenerResp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &map[string]interface{}{
		"ListenerId": listenerId,
	})
	if err != nil {
		return nil, err
	}

	listenerAttrMap := make(map[string]interface{})

	timeout, err := bp.ObtainSdkValue("Result.EstablishedTimeout", *listenerResp)
	if err != nil {
		return nil, err
	}
	desc, err := bp.ObtainSdkValue("Result.Description", *listenerResp)
	if err != nil {
		return nil, err
	}
	loadBalancerId, err := bp.ObtainSdkValue("Result.LoadBalancerId", *listenerResp)
	if err != nil {
		return nil, err
	}
	scheduler, err := bp.ObtainSdkValue("Result.Scheduler", *listenerResp)
	if err != nil {
		return nil, err
	}

	listenerAttrMap["EstablishedTimeout"] = timeout
	listenerAttrMap["Description"] = desc
	listenerAttrMap["LoadBalancerId"] = loadBalancerId
	listenerAttrMap["Scheduler"] = scheduler

	for k, v := range listenerAttrMap {
		data[k] = v
	}

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
			status, err = bp.ObtainSdkValue("Status", demo)
			if err != nil {
				return nil, "", err
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}

}

func (*VestackListenerService) WithResourceResponseHandlers(listener map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return listener, map[string]bp.ResponseConvert{
			"CAEnabled": {
				TargetField: "ca_enabled",
			},
			"CACertificateId": {
				TargetField: "ca_certificate_id",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackListenerService) refreshAclStatus() bp.CallFunc {
	return func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) error {
		var aclIds []string
		for k, v := range *call.SdkParam {
			if strings.HasPrefix(k, "AclIds.") {
				aclIds = append(aclIds, v.(string))
			}
		}
		if len(aclIds) > 0 {
			if err := s.checkAcl(aclIds); err != nil {
				return err
			}
		}
		return nil
	}
}

func (s *VestackListenerService) checkAcl(aclIds []string) error {
	return resource.Retry(20*time.Minute, func() *resource.RetryError {
		action := "DescribeAcls"
		req := make(map[string]interface{})
		for index, id := range aclIds {
			req["AclIds."+strconv.Itoa(index+1)] = id
		}
		logger.Debug(logger.ReqFormat, "DescribeAcls", aclIds)
		// create 的时候上限为5个，无需翻页
		resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		logger.Debug(logger.RespFormat, "DescribeAcls", aclIds, *resp)

		statusOK := true
		aclIdMap := make(map[string]bool)
		results, err := bp.ObtainSdkValue("Result.Acls", *resp)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		acls, ok := results.([]interface{})
		if !ok {
			return resource.NonRetryableError(fmt.Errorf("checkAcl Result.Acls is not Slice"))
		}
		for _, element := range acls {
			aclMap, ok := element.(map[string]interface{})
			if !ok {
				return resource.NonRetryableError(fmt.Errorf("checkAcl Acl is not map"))
			}
			aclIdMap[aclMap["AclId"].(string)] = true
			if aclMap["Status"] == "Deleting" {
				return resource.NonRetryableError(fmt.Errorf("acl is in deleting status"))
			} else if aclMap["Status"] != "Active" { // Creating / Configuring
				statusOK = false
				break
			}
		}
		if !statusOK {
			return resource.RetryableError(fmt.Errorf("acl still in waiting status"))
		}

		for _, aclId := range aclIds {
			if _, exist := aclIdMap[aclId]; !exist {
				return resource.NonRetryableError(fmt.Errorf("Cannot find acl: %s ", aclId))
			}
		}
		return nil
	})
}

func (s *VestackListenerService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateListener",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"acl_ids": {
					ConvertType: bp.ConvertWithN,
				},
				"health_check": {
					ConvertType: bp.ConvertListUnique,
					NextLevelConvert: map[string]bp.RequestConvert{
						"un_healthy_threshold": {
							TargetField: "UnhealthyThreshold",
						},
						"uri": {
							TargetField: "URI",
						},
					},
				},
				"port": {
					TargetField: "Port",
					ForceGet:    true,
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertListN,
				},
				"ca_enabled": {
					TargetField: "CAEnabled",
				},
				"ca_certificate_id": {
					TargetField: "CACertificateId",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				protocol := (*call.SdkParam)["Protocol"].(string)
				// 1. established_timeout
				if protocol == "HTTP" || protocol == "HTTPS" {
					// not allow establish_timeout
					if _, ok := (*call.SdkParam)["EstablishedTimeout"]; ok {
						return false, errors.New("established_timeout is not allowed for HTTP or HTTPS")
					}
				}

				// 2. certificate_id
				if protocol != "HTTPS" && (*call.SdkParam)["CertificateId"] != nil {
					return false, errors.New("certificate_id is only allowed for HTTPS")
				}

				// 3. connection_drain_timeout
				if d.Get("connection_drain_enabled") == "on" {
					timeout := d.Get("connection_drain_timeout").(int)
					if timeout == 0 {
						(*call.SdkParam)["ConnectionDrainTimeout"] = 0
					}
				}

				// 4. Only Https and Http support
				if protocol != "HTTP" && protocol != "HTTPS" {
					httpSpecificParams := []string{
						"ClientHeaderTimeout",
						"ClientBodyTimeout",
						"KeepaliveTimeout",
						"ProxyConnectTimeout",
						"ProxySendTimeout",
						"ProxyReadTimeout",
						"SendTimeout",
					}
					for _, param := range httpSpecificParams {
						if _, ok := (*call.SdkParam)[param]; ok {
							return false, fmt.Errorf("%s is only allowed for HTTP or HTTPS protocols", param)
						}
					}
				}

				// 5. Only Https support
				if protocol != "HTTPS" {
					httpsSpecificParams := []string{
						"SecurityPolicyId",
						"Http2Enabled",
					}
					for _, param := range httpsSpecificParams {
						if _, ok := (*call.SdkParam)[param]; ok {
							return false, fmt.Errorf("%s is only allowed for HTTPS protocols", param)
						}
					}
				}
				return true, nil
			},
			AfterLocked: s.refreshAclStatus(),
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//创建listener
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := bp.ObtainSdkValue("Result.ListenerId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active", "Disabled"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				clb.NewClbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("load_balancer_id").(string),
				},
			},
			AfterRefresh: s.refreshAclStatus(),
			LockId: func(d *schema.ResourceData) string {
				return resourceData.Get("load_balancer_id").(string)
			},
		},
	}
	return []bp.Callback{callback}

}

func (s *VestackListenerService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	clbId, err := s.queryLoadBalancerId(resourceData.Id())
	if err != nil {
		return []bp.Callback{{
			Err: err,
		}}
	}
	var callbacks []bp.Callback
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyListenerAttributes",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"acl_ids": {
					ConvertType: bp.ConvertWithN,
				},
				"health_check": {
					ConvertType: bp.ConvertListUnique,
					NextLevelConvert: map[string]bp.RequestConvert{
						"un_healthy_threshold": {
							TargetField: "UnhealthyThreshold",
						},
						"uri": {
							TargetField: "URI",
						},
					},
				},
				"ca_enabled": {
					TargetField: "CAEnabled",
				},
				"ca_certificate_id": {
					TargetField: "CACertificateId",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				protocol := d.Get("protocol").(string)
				// 1. established_timeout
				if protocol == "HTTP" || protocol == "HTTPS" {
					// not allow establish_timeout
					if _, ok := (*call.SdkParam)["EstablishedTimeout"]; ok {
						return false, errors.New("established_timeout is not allowed for HTTP or HTTPS")
					}
				}

				// 2. certificate_id
				if protocol != "HTTPS" && (*call.SdkParam)["CertificateId"] != nil {
					return false, errors.New("certificate_id is only allowed for HTTPS")
				}

				// 3. Only Https and Http support
				if protocol != "HTTP" && protocol != "HTTPS" {
					httpSpecificParams := []string{
						"ClientHeaderTimeout",
						"ClientBodyTimeout",
						"KeepaliveTimeout",
						"ProxyConnectTimeout",
						"ProxySendTimeout",
						"ProxyReadTimeout",
						"SendTimeout",
					}
					for _, param := range httpSpecificParams {
						if _, ok := (*call.SdkParam)[param]; ok {
							return false, fmt.Errorf("%s is only allowed for HTTP or HTTPS protocols", param)
						}
					}
				}

				// 4. Only Https support
				if protocol != "HTTPS" {
					httpsSpecificParams := []string{
						"SecurityPolicyId",
						"Http2Enabled",
					}
					for _, param := range httpsSpecificParams {
						if _, ok := (*call.SdkParam)[param]; ok {
							return false, fmt.Errorf("%s is only allowed for HTTPS protocols", param)
						}
					}
				}

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
			AfterLocked: s.refreshAclStatus(),
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//修改 listener 属性
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active", "Disabled"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				clb.NewClbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: clbId,
				},
			},
			AfterRefresh: s.refreshAclStatus(),
			LockId: func(d *schema.ResourceData) string {
				return clbId
			},
		},
	}
	callbacks = append(callbacks, callback)
	// 更新tags
	setResourceTagsCallbacks := bp.SetResourceTags(s.Client, "TagResources", "UntagResources", "listener", resourceData, getUniversalInfo)
	callbacks = append(callbacks, setResourceTagsCallbacks...)

	return callbacks
}

func (s *VestackListenerService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	clbId, err := s.queryLoadBalancerId(resourceData.Id())
	if err != nil {
		return []bp.Callback{{
			Err: err,
		}}
	}

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteListener",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"ListenerId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//删除 Listener
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
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
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
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
	return []bp.Callback{callback}
}

func (s *VestackListenerService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "ListenerIds",
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
		NameField:    "ListenerName",
		IdField:      "ListenerId",
		CollectField: "listeners",
		ResponseConverts: map[string]bp.ResponseConvert{
			"ListenerId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"CAEnabled": {
				TargetField: "ca_enabled",
			},
			"CACertificateId": {
				TargetField: "ca_certificate_id",
			},
			"HealthCheck.Port": {
				TargetField: "helth_check_port",
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
			"HealthCheck.UdpRequest": {
				TargetField: "health_check_udp_request",
			},
			"HealthCheck.UdpExpect": {
				TargetField: "health_check_udp_expect",
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

	// 查询 LoadBalancerId
	action := "DescribeListenerAttributes"
	serverGroupResp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &map[string]interface{}{
		"ListenerId": listenerId,
	})
	if err != nil {
		return "", err
	}
	clbId, err := bp.ObtainSdkValue("Result.LoadBalancerId", *serverGroupResp)
	if err != nil {
		return "", err
	}
	return clbId.(string), nil
}

func removeSystemTags(data []interface{}) ([]interface{}, error) {
	var (
		ok      bool
		result  map[string]interface{}
		results []interface{}
		tags    []interface{}
	)
	for _, d := range data {
		if result, ok = d.(map[string]interface{}); !ok {
			return results, errors.New("The elements in data are not map ")
		}
		tags, ok = result["Tags"].([]interface{})
		if ok {
			tags = bp.FilterSystemTags(tags)
			result["Tags"] = tags
		}
		results = append(results, result)
	}
	return results, nil
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
