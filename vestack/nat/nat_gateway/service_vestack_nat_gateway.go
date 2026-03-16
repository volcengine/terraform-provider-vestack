package nat_gateway

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackNatGatewayService struct {
	Client *bp.SdkClient
}

func NewNatGatewayService(c *bp.SdkClient) *VestackNatGatewayService {
	return &VestackNatGatewayService{
		Client: c,
	}
}

func (s *VestackNatGatewayService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackNatGatewayService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		action := "DescribeNatGateways"
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
		logger.Debug(logger.RespFormat, "testDescribeNatGateways", condition, *resp)

		results, err = bp.ObtainSdkValue("Result.NatGateways", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.NatGateways is not Slice")
		}
		return data, err
	})
}

func (s *VestackNatGatewayService) ReadResource(resourceData *schema.ResourceData, natGatewayId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if natGatewayId == "" {
		natGatewayId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"NatGatewayIds.1": natGatewayId,
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
		return data, fmt.Errorf("NatGateway %s not exist ", natGatewayId)
	}

	return data, err
}

func (s *VestackNatGatewayService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      3 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				demo   map[string]interface{}
				status interface{}
			)

			if err = resource.Retry(20*time.Minute, func() *resource.RetryError {
				demo, err = s.ReadResource(resourceData, id)
				if err != nil {
					if bp.ResourceNotFoundError(err) {
						return resource.RetryableError(err)
					} else {
						return resource.NonRetryableError(err)
					}
				}
				return nil
			}); err != nil {
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

func (VestackNatGatewayService) WithResourceResponseHandlers(natGateway map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return natGateway, map[string]bp.ResponseConvert{
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
						return "PostPaid"
					case 3:
						return "PostPaidByUsage"
					}
					return fmt.Sprintf("%v", i)
				},
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackNatGatewayService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateNatGateway",
			ConvertMode: bp.RequestConvertAll,
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
			Convert: map[string]bp.RequestConvert{
				"billing_type": {
					TargetField: "BillingType",
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
						case "PostPaidByUsage":
							return 3
						}
						return 0
					},
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertListN,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				// PeriodUnit 默认传 Month
				if (*call.SdkParam)["BillingType"] == 1 {
					(*call.SdkParam)["PeriodUnit"] = "Month"
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//创建natGateway
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := bp.ObtainSdkValue("Result.NatGatewayId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}

}

func (s *VestackNatGatewayService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyNatGatewayAttributes",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"nat_gateway_name": {
					TargetField: "NatGatewayName",
				},
				"description": {
					TargetField: "Description",
				},
				"spec": {
					TargetField: "Spec",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["NatGatewayId"] = d.Id()
				delete(*call.SdkParam, "Tags")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//修改natGateway属性
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, callback)

	// 更新Tags
	setResourceTagsCallbacks := bp.SetResourceTags(s.Client, "TagResources", "UntagResources", "ngw", resourceData, getVpcUniversalInfo)
	callbacks = append(callbacks, setResourceTagsCallbacks...)

	return callbacks
}

func (s *VestackNatGatewayService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	id := resourceData.Id()
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteNatGateway",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"NatGatewayId": id,
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//删除NatGateway
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//由于异步删除问题 这里补充一个轮询查询(临时解决方案)
				return resource.Retry(3*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, id)
					//能查询成功代表还在删除中，重试
					if callErr == nil {
						return resource.RetryableError(fmt.Errorf("Nat still in remobp "))
					} else {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(callErr)
						}
					}
				})
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading nat gateway on delete %q, %w", d.Id(), callErr))
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

func (s *VestackNatGatewayService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "NatGatewayIds",
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
		NameField:    "NatGatewayName",
		IdField:      "NatGatewayId",
		CollectField: "nat_gateways",
		ResponseConverts: map[string]bp.ResponseConvert{
			"NatGatewayId": {
				TargetField: "id",
				KeepDefault: true,
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
						return "PostPaid"
					case 3:
						return "PostPaidByUsage"
					}
					return fmt.Sprintf("%v", i)
				},
			},
		},
	}
}

func (s *VestackNatGatewayService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "natgateway",
		Action:      actionName,
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
	}
}

func getVpcUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vpc",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}

func (s *VestackNatGatewayService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "natgateway",
		ResourceType:         "ngw",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}

func (s *VestackNatGatewayService) UnsubscribeInfo(resourceData *schema.ResourceData, resource *schema.Resource) (*bp.UnsubscribeInfo, error) {
	info := bp.UnsubscribeInfo{
		InstanceId: s.ReadResourceId(resourceData.Id()),
	}
	if resourceData.Get("billing_type") == "PrePaid" {
		info.Products = []string{"NAT_Gateway"}
		info.NeedUnsubscribe = true
	}
	return &info, nil
}
