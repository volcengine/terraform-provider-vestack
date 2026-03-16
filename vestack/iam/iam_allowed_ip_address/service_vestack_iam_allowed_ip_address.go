package iam_allowed_ip_address

import (
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamAllowedIpAddressService struct {
	Client *bp.SdkClient
}

func NewIamAllowedIpAddressService(c *bp.SdkClient) *VestackIamAllowedIpAddressService {
	return &VestackIamAllowedIpAddressService{
		Client: c,
	}
}

func (s *VestackIamAllowedIpAddressService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamAllowedIpAddressService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	action := "GetAllowedIPAddresses"
	logger.Debug(logger.ReqFormat, action, m)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &m)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp)

	// 根据返回结构，Result 是顶级字段
	result, err := bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return data, err
	}

	if resMap, ok := result.(map[string]interface{}); ok {
		data = append(data, resMap)
	} else {
		return data, errors.New("Result is not map")
	}

	return data, nil
}

func (s *VestackIamAllowedIpAddressService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	results, err := s.ReadResources(nil)
	if err != nil {
		return data, err
	}
	if len(results) > 0 {
		if resMap, ok := results[0].(map[string]interface{}); ok {
			return resMap, nil
		}
		return data, errors.New("results[0] is not map")
	}
	return data, errors.New("Allowed IP Address not found")
}

func (s *VestackIamAllowedIpAddressService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (VestackIamAllowedIpAddressService) WithResourceResponseHandlers(v map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return v, map[string]bp.ResponseConvert{
			"IPList": {
				TargetField: "ip_list",
				Convert: func(i interface{}) interface{} {
					if i == nil {
						return nil
					}
					var list []interface{}
					if rawList, ok := i.([]interface{}); ok {
						for _, item := range rawList {
							if m, ok := item.(map[string]interface{}); ok {
								list = append(list, map[string]interface{}{
									"ip":          m["IP"],
									"description": m["Description"],
								})
							}
						}
					}
					return list
				},
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackIamAllowedIpAddressService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateAllowedIPAddresses",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"enable_ip_list": {
					TargetField: "EnableIPList",
				},
				"ip_list": {
					TargetField: "IPList",
					ConvertType: bp.ConvertJsonObjectArray,
					NextLevelConvert: map[string]bp.RequestConvert{
						"ip": {
							TargetField: "IP",
						},
						"description": {
							TargetField: "Description",
						},
					},
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId("iam_allowed_ip_address")
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackIamAllowedIpAddressService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return s.CreateResource(resourceData, resource)
}

func (s *VestackIamAllowedIpAddressService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateAllowedIPAddresses",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["EnableIPList"] = false
				(*call.SdkParam)["IPList"] = []interface{}{}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackIamAllowedIpAddressService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		CollectField: "allowed_ip_addresses",
		ResponseConverts: map[string]bp.ResponseConvert{
			"IPList": {
				TargetField: "ip_list",
				Convert: func(i interface{}) interface{} {
					if i == nil {
						return nil
					}
					var list []interface{}
					if rawList, ok := i.([]interface{}); ok {
						for _, item := range rawList {
							if m, ok := item.(map[string]interface{}); ok {
								list = append(list, map[string]interface{}{
									"ip":          m["IP"],
									"description": m["Description"],
								})
							}
						}
					}
					return list
				},
			},
		},
	}
}

func (s *VestackIamAllowedIpAddressService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	method := bp.GET
	contentType := bp.Default
	if actionName == "UpdateAllowedIPAddresses" {
		method = bp.POST
		contentType = bp.ApplicationJSON
	}
	return bp.UniversalInfo{
		ServiceName: "iam",
		Action:      actionName,
		Version:     "2018-01-01",
		HttpMethod:  method,
		ContentType: contentType,
		RegionType:  bp.Global,
	}
}
