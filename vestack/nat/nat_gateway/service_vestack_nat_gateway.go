package nat_gateway

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ve "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackNatGatewayService struct {
	Client     *ve.SdkClient
	Dispatcher *ve.Dispatcher
}

func NewNatGatewayService(c *ve.SdkClient) *VestackNatGatewayService {
	return &VestackNatGatewayService{
		Client:     c,
		Dispatcher: &ve.Dispatcher{},
	}
}

func (s *VestackNatGatewayService) GetClient() *ve.SdkClient {
	return s.Client
}

func (s *VestackNatGatewayService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return ve.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		natGateway := s.Client.NatClient
		action := "DescribeNatGateways"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = natGateway.DescribeNatGatewaysCommon(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = natGateway.DescribeNatGatewaysCommon(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = ve.ObtainSdkValue("Result.NatGateways", *resp)
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
		results     []interface{}
		ok          bool
		billingType interface{}
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

	billingType, err = ve.ObtainSdkValue("BillingType", data)
	if err != nil {
		return data, err
	}

	if billingType, ok = billingType.(float64); !ok {
		return data, errors.New("BillingType is invalid")
	}

	if billingType.(float64) == 1 {
		// prepaid not support
		return data, fmt.Errorf("Prepaid nat gateway %s is not supported ", natGatewayId)
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

func (VestackNatGatewayService) WithResourceResponseHandlers(natGateway map[string]interface{}) []ve.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]ve.ResponseConvert, error) {
		return natGateway, map[string]ve.ResponseConvert{
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
					}
					return i
				},
			},
		}, nil
	}
	return []ve.ResourceResponseHandler{handler}

}

func (s *VestackNatGatewayService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []ve.Callback {
	callback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "CreateNatGateway",
			ConvertMode: ve.RequestConvertAll,
			Convert: map[string]ve.RequestConvert{
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
						}
						return i
					},
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//??????natGateway
				return s.Client.NatClient.CreateNatGatewayCommon(call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *ve.SdkClient, resp *map[string]interface{}, call ve.SdkCall) error {
				//?????? ???????????? ??????????????????????????? ???????????????
				id, _ := ve.ObtainSdkValue("Result.NatGatewayId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &ve.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []ve.Callback{callback}

}

func (s *VestackNatGatewayService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []ve.Callback {
	callback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "ModifyNatGatewayAttributes",
			ConvertMode: ve.RequestConvertAll,
			Convert: map[string]ve.RequestConvert{
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
						}
						return i
					},
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (bool, error) {
				(*call.SdkParam)["NatGatewayId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//??????natGateway??????
				return s.Client.NatClient.ModifyNatGatewayAttributesCommon(call.SdkParam)
			},
			Refresh: &ve.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []ve.Callback{callback}
}

func (s *VestackNatGatewayService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []ve.Callback {
	id := resourceData.Id()
	callback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "DeleteNatGateway",
			ConvertMode: ve.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"NatGatewayId": id,
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//??????NatGateway
				return s.Client.NatClient.DeleteNatGatewayCommon(call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *ve.SdkClient, resp *map[string]interface{}, call ve.SdkCall) error {
				//???????????????????????? ??????????????????????????????(??????????????????)
				return resource.Retry(3*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, id)
					//?????????????????????????????????????????????
					if callErr == nil {
						return resource.RetryableError(fmt.Errorf("Nat still in remove "))
					} else {
						if ve.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(callErr)
						}
					}
				})
			},
			CallError: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall, baseErr error) error {
				//?????????????????????
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if ve.ResourceNotFoundError(callErr) {
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
	return []ve.Callback{callback}
}

func (s *VestackNatGatewayService) DatasourceResources(*schema.ResourceData, *schema.Resource) ve.DataSourceInfo {
	return ve.DataSourceInfo{
		RequestConverts: map[string]ve.RequestConvert{
			"ids": {
				TargetField: "NatGatewayIds",
				ConvertType: ve.ConvertWithN,
			},
		},
		NameField:    "NatGatewayName",
		IdField:      "NatGatewayId",
		CollectField: "nat_gateways",
		ResponseConverts: map[string]ve.ResponseConvert{
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
					}
					return i
				},
			},
		},
	}
}

func (s *VestackNatGatewayService) ReadResourceId(id string) string {
	return id
}
