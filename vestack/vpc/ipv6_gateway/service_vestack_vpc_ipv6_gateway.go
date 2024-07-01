package ipv6_gateway

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

type VestackIpv6GatewayService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewIpv6GatewayService(c *bp.SdkClient) *VestackIpv6GatewayService {
	return &VestackIpv6GatewayService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackIpv6GatewayService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIpv6GatewayService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		action := "DescribeIpv6Gateways"
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

		results, err = bp.ObtainSdkValue("Result.Ipv6Gateways", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Ipv6Gateways is not Slice")
		}

		return data, err
	})
}

func (s *VestackIpv6GatewayService) ReadResource(resourceData *schema.ResourceData, ipv6GatewayId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if ipv6GatewayId == "" {
		ipv6GatewayId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"Ipv6GatewayIds.1": ipv6GatewayId,
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
		return data, fmt.Errorf("Ipv6Gateway %s is not exist ", ipv6GatewayId)
	}

	return data, err
}

func (s *VestackIpv6GatewayService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			status, err = bp.ObtainSdkValue("Status", demo)
			if err != nil {
				return nil, "", err
			}
			return demo, status.(string), err
		},
	}
}

func (VestackIpv6GatewayService) WithResourceResponseHandlers(ipv6Gateway map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return ipv6Gateway, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackIpv6GatewayService) CreateResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateIpv6Gateway",
			ConvertMode: bp.RequestConvertAll,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.Ipv6GatewayId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
		},
	}
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackIpv6GatewayService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyIpv6GatewayAttribute",
			ConvertMode: bp.RequestConvertAll,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Ipv6GatewayId"] = d.Id()
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
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
		},
	}
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackIpv6GatewayService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	removeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteIpv6Gateway",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"Ipv6GatewayId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading Ipv6Gateway on delete %q, %w", d.Id(), callErr))
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

func (s *VestackIpv6GatewayService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "Ipv6GatewayIds",
				ConvertType: bp.ConvertWithN,
			},
			"vpc_ids": {
				TargetField: "VpcIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "Name",
		IdField:      "Ipv6GatewayId",
		CollectField: "ipv6_gateways",
		ResponseConverts: map[string]bp.ResponseConvert{
			"Ipv6GatewayId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *VestackIpv6GatewayService) ReadResourceId(id string) string {
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
