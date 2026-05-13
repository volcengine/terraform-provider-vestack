package clb

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackClbService struct {
	Client *bp.SdkClient
}

func NewClbService(c *bp.SdkClient) *VestackClbService {
	return &VestackClbService{
		Client: c,
	}
}

func (s *VestackClbService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackClbService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	// 使用分页查询获取所有CLB实例
	data, err = bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		action := "DescribeLoadBalancers"
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
		logger.Debug(logger.RespFormat, action, condition, *resp)

		results, err = bp.ObtainSdkValue("Result.LoadBalancers", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.LoadBalancers is not Slice")
		}
		// 移除系统标签
		data, err = removeSystemTags(data)
		return data, err
	})
	if err != nil {
		return data, err
	}

	// 为每个CLB实例获取详细信息和计费信息
	for _, value := range data {
		clb, ok := value.(map[string]interface{})
		if !ok {
			return data, fmt.Errorf(" Clb is not map ")
		}

		// 获取CLB详细信息，包括EIP配置和IPv6带宽信息
		eipAction := "DescribeLoadBalancerAttributes"
		eipReq := map[string]interface{}{
			"LoadBalancerId": clb["LoadBalancerId"],
		}
		logger.Debug(logger.ReqFormat, eipAction, eipReq)
		eipResp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(eipAction), &eipReq)
		if err != nil {
			return data, err
		}
		logger.Debug(logger.RespFormat, eipAction, *eipResp)

		eipConfig, err := bp.ObtainSdkValue("Result.Eip", *eipResp)
		if err != nil {
			return data, err
		}
		clb["EipBillingConfig"] = eipConfig

		ipv6EipConfig, err := bp.ObtainSdkValue("Result.Ipv6AddressBandwidth", *eipResp)
		if err != nil {
			return data, err
		}
		clb["Ipv6AddressBandwidth"] = ipv6EipConfig

		logTopicId, err := bp.ObtainSdkValue("Result.LogTopicId", *eipResp)
		if err != nil {
			return data, err
		}
		clb["LogTopicId"] = logTopicId

		enabled, err := bp.ObtainSdkValue("Result.Enabled", *eipResp)
		if err != nil {
			return data, err
		}
		clb["Enabled"] = enabled

		listeners, err := bp.ObtainSdkValue("Result.Listeners", *eipResp)
		if err != nil {
			return data, err
		}
		clb["Listeners"] = listeners

		serverGroups, err := bp.ObtainSdkValue("Result.ServerGroups", *eipResp)
		if err != nil {
			return data, err
		}
		clb["ServerGroups"] = serverGroups

		accessLog, err := bp.ObtainSdkValue("Result.AccessLog", *eipResp)
		if err != nil {
			return data, err
		}
		clb["AccessLog"] = accessLog

		// `PostPaid` 实例不需查询续费相关信息
		if billingType := clb["LoadBalancerBillingType"]; billingType == 2.0 {
			continue
		}
		// 获取计费信息（仅对PrePaid类型）
		billingAction := "DescribeLoadBalancersBilling"
		billingReq := map[string]interface{}{
			"LoadBalancerIds.1": clb["LoadBalancerId"],
		}
		logger.Debug(logger.ReqFormat, billingAction, billingReq)
		billingResp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(billingAction), &billingReq)
		if err != nil {
			return data, err
		}
		logger.Debug(logger.RespFormat, billingAction, *billingResp)

		billingConfigs, err := bp.ObtainSdkValue("Result.LoadBalancerBillingConfigs", *billingResp)
		if err != nil {
			return data, err
		}
		if billingConfigs == nil {
			return data, fmt.Errorf(" DescribeLoadBalancersBilling error ")
		}
		configs, ok := billingConfigs.([]interface{})
		if !ok {
			return data, fmt.Errorf(" Result.LoadBalancerBillingConfigs is not slice ")
		}
		if len(configs) == 0 {
			return data, fmt.Errorf("LoadBalancerBilling of the clb instance %s is not exist ", clb["LoadBalancerId"])
		}
		config, ok := configs[0].(map[string]interface{})
		if !ok {
			return data, fmt.Errorf(" BillingConfigs is not map ")
		}
		for k, v := range config {
			clb[k] = v
		}
	}

	return data, err
}

func (s *VestackClbService) ReadResource(resourceData *schema.ResourceData, clbId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if clbId == "" {
		clbId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"LoadBalancerIds.1": clbId,
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
		return data, fmt.Errorf("Clb %s not exist ", clbId)
	}

	data["RegionId"] = s.Client.Region

	return data, err
}

func (s *VestackClbService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				demo       map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "CreateFailed")
			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", demo)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("Clb  status  error, status:%s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}

}

func (VestackClbService) WithResourceResponseHandlers(clb map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return clb, map[string]bp.ResponseConvert{
			"LoadBalancerBillingType": {
				TargetField: "load_balancer_billing_type",
				Convert: func(i interface{}) interface{} {
					if i == nil {
						return nil
					}
					billingType := i.(float64)
					switch billingType {
					case 1:
						return "PrePaid"
					case 2:
						return "PostPaid"
					case 3:
						return "PostPaidByLCU"
					}
					return i
				},
			},
			"RenewType": {
				TargetField: "renew_type",
				Convert: func(i interface{}) interface{} {
					if i == nil {
						return nil
					}
					renewType := i.(float64)
					switch renewType {
					case 1:
						return "ManualRenew"
					case 2:
						return "AutoRenew"
					case 3:
						return "NoneRenew"
					}
					return i
				},
			},
			"EipID": {
				TargetField: "eip_id",
			},
			"ISP": {
				TargetField: "isp",
			},
			"EipBillingType": {
				TargetField: "eip_billing_type",
				Convert: func(i interface{}) interface{} {
					if i == nil {
						return nil
					}
					billingType := i.(float64)
					switch billingType {
					case 1:
						return "PrePaid"
					case 2:
						return "PostPaidByBandwidth"
					case 3:
						return "PostPaidByTraffic"
					}
					return ""
				},
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackClbService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateLoadBalancer",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"load_balancer_billing_type": {
					TargetField: "LoadBalancerBillingType",
					Convert: func(data *schema.ResourceData, i interface{}) interface{} {
						if i == nil {
							return nil
						}
						billingType := i.(string)
						switch billingType {
						case "PrePaid":
							return 1
						case "PostPaid":
							return 2
						case "PostPaidByLCU":
							return 3
						}
						return i
					},
				},
				"eip_billing_config": {
					TargetField: "EipBillingConfig",
					ConvertType: bp.ConvertListUnique,
					NextLevelConvert: map[string]bp.RequestConvert{
						"isp": {
							TargetField: "ISP",
						},
						"security_protection_types": {
							TargetField: "SecurityProtectionTypes",
							ConvertType: bp.ConvertWithN,
						},
					},
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertListN,
				},
				"renew_type": {
					Ignore: true,
				},
				"renew_period_times": {
					Ignore: true,
				},
				"remain_renew_times": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if regionId, ok := (*call.SdkParam)["RegionId"]; !ok {
					(*call.SdkParam)["RegionId"] = s.Client.Region
				} else if regionId.(string) != s.Client.Region {
					return false, fmt.Errorf("region_id is not equal to provider region config(%s)", s.Client.Region)
				}

				// private 类型不传 eip_billing_config
				if (*call.SdkParam)["Type"] == "private" {
					delete(*call.SdkParam, "EipBillingConfig.ISP")
					delete(*call.SdkParam, "EipBillingConfig.EipBillingType")
					delete(*call.SdkParam, "EipBillingConfig.Bandwidth")
					delete(*call.SdkParam, "EipBillingConfig.BandwidthPackageId")
					delete(*call.SdkParam, "EipBillingConfig.SecurityProtectionTypes")
					delete(*call.SdkParam, "EipBillingConfig.SecurityProtectionInstanceId")
				}
				if eipBillingType, exist := (*call.SdkParam)["EipBillingConfig.EipBillingType"]; exist {
					ty := 0
					switch eipBillingType.(string) {
					case "PrePaid":
						ty = 1
					case "PostPaidByBandwidth":
						ty = 2
					case "PostPaidByTraffic":
						ty = 3
					}
					(*call.SdkParam)["EipBillingConfig.EipBillingType"] = ty
				}

				// PeriodUnit 默认传 Month
				if (*call.SdkParam)["LoadBalancerBillingType"] == 1 {
					(*call.SdkParam)["PeriodUnit"] = "Month"
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				//创建clb
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := bp.ObtainSdkValue("Result.LoadBalancerId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, callback)

	// 只有在创建 PrePaid 类型 CLB 且设置了有效的 renew_type 时才需要设置续费类型
	if billingType, ok := resourceData.GetOk("load_balancer_billing_type"); ok {
		if billingType.(string) == "PrePaid" {
			if renewType, ok := resourceData.GetOk("renew_type"); ok && renewType.(string) != "" {
				renewCallback := s.setLoadBalancerRenewal(resourceData)
				callbacks = append(callbacks, renewCallback...)
			}
		}
	}
	return callbacks

}
func (s *VestackClbService) setLoadBalancerRenewal(resourceData *schema.ResourceData) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "SetLoadBalancerRenewal",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"renew_type": {
					TargetField: "RenewType",
					ForceGet:    true,
					Convert: func(data *schema.ResourceData, i interface{}) interface{} {
						if i == nil {
							return nil
						}
						renewType := i.(string)
						switch renewType {
						case "ManualRenew":
							return 1
						case "AutoRenew":
							return 2
						}
						return i
					},
				},
				"renew_period_times": {
					TargetField: "RenewPeriodTimes",
					ForceGet:    true,
				},
				"remain_renew_times": {
					TargetField: "RemainRenewTimes",
					ForceGet:    true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) > 0 {
					(*call.SdkParam)["LoadBalancerId"] = d.Id()
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
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackClbService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	// 修改CLB基本属性
	attributesCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyLoadBalancerAttributes",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"load_balancer_name": {
					TargetField: "LoadBalancerName",
				},
				"description": {
					TargetField: "Description",
				},
				"modification_protection_status": {
					TargetField: "ModificationProtectionStatus",
				},
				"modification_protection_reason": {
					TargetField: "ModificationProtectionReason",
				},
				"load_balancer_spec": {
					TargetField: "LoadBalancerSpec",
				},
				"address_ip_version": {
					TargetField: "AddressIpVersion",
				},
				"eni_ipv6_address": {
					TargetField: "EniIpv6Address",
				},
				"bypass_security_group_enabled": {
					TargetField: "BypassSecurityGroupEnabled",
				},
				"timestamp_remove_enabled": {
					TargetField: "TimestampRemoveEnabled",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				// PostPaidByLCU 类型不支持指定规格
				oldType, _ := d.GetChange("load_balancer_billing_type")
				if oldType == "PostPaidByLCU" {
					delete(*call.SdkParam, "LoadBalancerSpec")
				}
				if len(*call.SdkParam) > 0 {
					(*call.SdkParam)["LoadBalancerId"] = d.Id()
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//修改clb属性
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	callbacks = append(callbacks, attributesCallback)

	// 处理计费类型变更
	// 处理 renew_type 变更
	if resourceData.HasChanges("renew_type", "renew_period_times", "remain_renew_times") && resourceData.Get("renew_type").(string) != "" {
		renewalCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "SetLoadBalancerRenewal",
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"renew_type": {
						TargetField: "RenewType",
						ForceGet:    true,
						Convert: func(data *schema.ResourceData, i interface{}) interface{} {
							if i == nil {
								return nil
							}
							renewType := i.(string)
							switch renewType {
							case "ManualRenew":
								return 1
							case "AutoRenew":
								return 2
							}
							return i
						},
					},
					"renew_period_times": {
						TargetField: "RenewPeriodTimes",
						ForceGet:    true,
					},
					"remain_renew_times": {
						TargetField: "RemainRenewTimes",
						ForceGet:    true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if billType, ok := d.GetOk("load_balancer_billing_type"); ok && billType.(string) != "PrePaid" {
						return false, fmt.Errorf("renew_type can only be set when load_balancer_billing_type is PrePaid")
					}
					if len(*call.SdkParam) > 0 {
						(*call.SdkParam)["LoadBalancerId"] = d.Id()
						return true, nil
					}
					return false, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					return nil
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Active"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, renewalCallback)
	}
	if resourceData.HasChange("load_balancer_billing_type") {
		billingTypeCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ConvertLoadBalancerBillingType",
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"load_balancer_billing_type": {
						TargetField: "LoadBalancerBillingType",
						Convert: func(data *schema.ResourceData, i interface{}) interface{} {
							if i == nil {
								return nil
							}
							billingType := i.(string)
							switch billingType {
							case "PrePaid":
								return 1
							case "PostPaid":
								return 2
							case "PostPaidByLCU":
								return 3
							}
							return i
						},
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) > 0 {
						(*call.SdkParam)["LoadBalancerId"] = d.Id()
						oldType, newType := d.GetChange("load_balancer_billing_type")
						if oldType == "PostPaidByLCU" && newType == "PostPaid" {
							(*call.SdkParam)["LoadBalancerSpec"] = d.Get("load_balancer_spec")
						}

						if (*call.SdkParam)["LoadBalancerBillingType"].(int) == 1 {
							// PeriodUnit 默认传 Month
							(*call.SdkParam)["PeriodUnit"] = "Month"
							(*call.SdkParam)["Period"] = d.Get("period")
						}
						return true, nil
					}
					return false, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					//修改 clb 计费类型
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					time.Sleep(10 * time.Second)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Active"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, billingTypeCallback)
	} else if resourceData.Get("renew_type").(string) == "ManualRenew" && resourceData.HasChange("period") {
		renewCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "RenewLoadBalancer",
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"period": {
						TargetField: "Period",
						Convert: func(data *schema.ResourceData, i interface{}) interface{} {
							o, n := data.GetChange("period")
							return n.(int) - o.(int)
						},
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) > 0 {
						if (*call.SdkParam)["Period"].(int) <= 0 {
							return false, fmt.Errorf("period can only be enlarged ")
						}

						// PeriodUnit 默认传 Month
						(*call.SdkParam)["PeriodUnit"] = "Month"
						(*call.SdkParam)["LoadBalancerId"] = d.Id()
						return true, nil
					}
					return false, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					return nil
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Active"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, renewCallback)
	}

	// 更新Tags
	setResourceTagsCallbacks := bp.SetResourceTags(s.Client, "TagResources", "UntagResources", "CLB", resourceData, getUniversalInfo)
	callbacks = append(callbacks, setResourceTagsCallbacks...)

	return callbacks
}

func (s *VestackClbService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteLoadBalancer",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"LoadBalancerId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//删除Clb
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
							return resource.NonRetryableError(fmt.Errorf("error on  reading clb on delete %q, %w", d.Id(), callErr))
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

func (s *VestackClbService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "LoadBalancerIds",
				ConvertType: bp.ConvertWithN,
			},
			"instance_ids": {
				TargetField: "InstanceIds",
				ConvertType: bp.ConvertWithN,
			},
			"instance_ips": {
				TargetField: "InstanceIps",
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
		NameField:    "LoadBalancerName",
		IdField:      "LoadBalancerId",
		CollectField: "clbs",
		ResponseConverts: map[string]bp.ResponseConvert{
			"LoadBalancerId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"EipID": {
				TargetField: "eip_id",
			},
			"EniID": {
				TargetField: "eni_id",
			},
			"LoadBalancerBillingType": {
				TargetField: "load_balancer_billing_type",
				Convert: func(i interface{}) interface{} {
					if i == nil {
						return nil
					}
					billingType := i.(float64)
					switch billingType {
					case 1:
						return "PrePaid"
					case 2:
						return "PostPaid"
					case 3:
						return "PostPaidByLCU"
					}
					return i
				},
			},
			"RenewType": {
				TargetField: "renew_type",
				Convert: func(i interface{}) interface{} {
					if i == nil {
						return nil
					}
					renewType := i.(float64)
					switch renewType {
					case 1:
						return "ManualRenew"
					case 2:
						return "AutoRenew"
					case 3:
						return "NoneRenew"
					}
					return i
				},
			},
			"ISP": {
				TargetField: "isp",
			},
			"EipBillingType": {
				TargetField: "eip_billing_type",
				Convert: func(i interface{}) interface{} {
					if i == nil {
						return nil
					}
					billingType := i.(float64)
					switch billingType {
					case 1:
						return "PrePaid"
					case 2:
						return "PostPaidByBandwidth"
					case 3:
						return "PostPaidByTraffic"
					}
					return ""
				},
			},
			"BillingType": {
				TargetField: "billing_type",
				Convert: func(i interface{}) interface{} {
					if i == nil {
						return nil
					}
					billingType := i.(float64)
					switch billingType {
					case 1:
						return "PrePaid"
					case 2:
						return "PostPaidByBandwidth"
					case 3:
						return "PostPaidByTraffic"
					}
					return ""
				},
			},
		},
	}
}

func (s *VestackClbService) ReadResourceId(id string) string {
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

func (s *VestackClbService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "clb",
		ResourceType:         "clb",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
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

func (s *VestackClbService) UnsubscribeInfo(resourceData *schema.ResourceData, resource *schema.Resource) (*bp.UnsubscribeInfo, error) {
	info := bp.UnsubscribeInfo{
		InstanceId: s.ReadResourceId(resourceData.Id()),
	}
	if resourceData.Get("load_balancer_billing_type") == "PrePaid" {
		info.Products = []string{"CLB"}
		info.NeedUnsubscribe = true
	}
	return &info, nil
}
