package ipv6_address_bandwidth

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIpv6AddressBandwidthService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewIpv6AddressBandwidthService(c *bp.SdkClient) *VestackIpv6AddressBandwidthService {
	return &VestackIpv6AddressBandwidthService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackIpv6AddressBandwidthService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIpv6AddressBandwidthService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		action := "DescribeIpv6AddressBandwidths"
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

		results, err = bp.ObtainSdkValue("Result.Ipv6AddressBandwidths", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Ipv6AddressBandwidths is not Slice")
		}

		return data, err
	})
}

func (s *VestackIpv6AddressBandwidthService) ReadResource(resourceData *schema.ResourceData, allocationId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if allocationId == "" {
		allocationId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"AllocationIds.1": allocationId,
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
		return data, fmt.Errorf("Ipv6AddressBandwidth %s is not exist ", allocationId)
	}

	return data, err
}

func (s *VestackIpv6AddressBandwidthService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				demo   map[string]interface{}
				status interface{}
			)
			//no failed status.

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
			return demo, status.(string), err
		},
	}
}

func (VestackIpv6AddressBandwidthService) WithResourceResponseHandlers(ipv6AddressBandwidth map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return ipv6AddressBandwidth, map[string]bp.ResponseConvert{
			"ISP": {
				TargetField: "isp",
			},
			"BillingType": {
				TargetField: "billing_type",
				Convert:     billingTypeResponseConvert,
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackIpv6AddressBandwidthService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AllocateIpv6AddressBandwidth",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"billing_type": {
					TargetField: "BillingType",
					Convert:     billingTypeRequestConvert,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.AllocationId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackIpv6AddressBandwidthService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyIpv6AddressBandwidth",
			ConvertMode: bp.RequestConvertAll,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["AllocationId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackIpv6AddressBandwidthService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	removeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ReleaseIpv6AddressBandwidth",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"AllocationId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
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
							return resource.NonRetryableError(fmt.Errorf("error on  reading Ipv6AddressBandwidth on delete %q, %w", d.Id(), callErr))
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
	callbacks = append(callbacks, removeCallback)

	return callbacks
}

func (s *VestackIpv6AddressBandwidthService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "AllocationIds",
				ConvertType: bp.ConvertWithN,
			},
			"isp": {
				TargetField: "ISP",
			},
			"ipv6_addresses": {
				TargetField: "Ipv6Addresses",
				ConvertType: bp.ConvertWithN,
			},
		},
		IdField:      "AllocationId",
		CollectField: "ipv6_address_bandwidths",
		ResponseConverts: map[string]bp.ResponseConvert{
			"AllocationId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"ISP": {
				TargetField: "isp",
			},
			"BillingType": {
				TargetField: "billing_type",
				Convert:     billingTypeResponseConvert,
			},
		},
	}
}

func (s *VestackIpv6AddressBandwidthService) ReadResourceId(id string) string {
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
