package network_interface

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackNetworkInterfaceService struct {
	Client *bp.SdkClient
}

func NewNetworkInterfaceService(c *bp.SdkClient) *VestackNetworkInterfaceService {
	return &VestackNetworkInterfaceService{
		Client: c,
	}
}

func (s *VestackNetworkInterfaceService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackNetworkInterfaceService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		vpcClient := s.Client.VpcClient
		action := "DescribeNetworkInterfaces"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = vpcClient.DescribeNetworkInterfacesCommon(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = vpcClient.DescribeNetworkInterfacesCommon(&condition)
			if err != nil {
				return data, err
			}
		}
		logger.Debug(logger.RespFormat, action, *resp)
		results, err = bp.ObtainSdkValue("Result.NetworkInterfaceSets", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.NetworkInterfaceSets is not Slice")
		}
		return data, err
	})
}

func (s *VestackNetworkInterfaceService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"NetworkInterfaceIds.1": id,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("value is not map")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("network_interface %s not exist ", id)
	}
	privateIpAddress := make([]string, 0)
	if privateIpMap, ok := data["PrivateIpSets"].(map[string]interface{}); ok {
		if privateIpSets, ok := privateIpMap["PrivateIpSet"].([]interface{}); ok {
			for _, p := range privateIpSets {
				if pMap, ok := p.(map[string]interface{}); ok {
					isPrimary := pMap["Primary"].(bool)
					ip := pMap["PrivateIpAddress"].(string)
					if !isPrimary {
						privateIpAddress = append(privateIpAddress, ip)
					}
				}
			}
		}
	}
	data["PrivateIpAddress"] = privateIpAddress
	data["SecondaryPrivateIpAddressCount"] = len(privateIpAddress)
	return data, err
}

func (s *VestackNetworkInterfaceService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			failStates = append(failStates, "Error")
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
					return nil, "", fmt.Errorf("network_interface status error, status:%s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}
}

func (VestackNetworkInterfaceService) WithResourceResponseHandlers(networkInterface map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return networkInterface, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackNetworkInterfaceService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateNetworkInterface",
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.VpcClient.CreateNetworkInterfaceCommon(call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.NetworkInterfaceId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			Convert: map[string]bp.RequestConvert{
				"security_group_ids": {
					TargetField: "SecurityGroupIds",
					ConvertType: bp.ConvertWithN,
				},
				"private_ip_address": {
					TargetField: "PrivateIpAddress",
					ConvertType: bp.ConvertWithN,
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertListN,
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackNetworkInterfaceService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyNetworkInterfaceAttributes",
			ConvertMode: bp.RequestConvertAll,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["NetworkInterfaceId"] = d.Id()
				delete(*call.SdkParam, "Tags")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.VpcClient.ModifyNetworkInterfaceAttributesCommon(call.SdkParam)
			},
			Convert: map[string]bp.RequestConvert{
				"security_group_ids": {
					TargetField: "SecurityGroupIds",
					ConvertType: bp.ConvertWithN,
				},
				"private_ip_address": {
					Ignore: true,
				},
				"secondary_private_ip_address_count": {
					Ignore: true,
				},
			},
		},
	}
	callbacks = append(callbacks, callback)

	// 检查private_ip_address改变
	if resourceData.HasChange("private_ip_address") {
		add, remove, _, _ := bp.GetSetDifference("private_ip_address", resourceData, schema.HashString, false)
		if remove.Len() > 0 {
			callback = bp.Callback{
				Call: bp.SdkCall{
					Action:      "UnassignPrivateIpAddresses",
					ConvertMode: bp.RequestConvertInConvert,
					BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
						(*call.SdkParam)["NetworkInterfaceId"] = d.Id()
						for index, r := range remove.List() {
							(*call.SdkParam)["PrivateIpAddress."+strconv.Itoa(index+1)] = r
						}
						return true, nil
					},
					Convert: map[string]bp.RequestConvert{
						"private_ip_address": {
							Ignore: true,
						},
						"secondary_private_ip_address_count": {
							Ignore: true,
						},
					},
					ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
						logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
						return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					},
				},
			}
			callbacks = append(callbacks, callback)
		}
		if add.Len() > 0 {
			callback = bp.Callback{
				Call: bp.SdkCall{
					Action:      "AssignPrivateIpAddresses",
					ConvertMode: bp.RequestConvertInConvert,
					BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
						(*call.SdkParam)["NetworkInterfaceId"] = d.Id()
						for index, r := range add.List() {
							(*call.SdkParam)["PrivateIpAddress."+strconv.Itoa(index+1)] = r
						}
						return true, nil
					},
					Convert: map[string]bp.RequestConvert{
						"private_ip_address": {
							Ignore: true,
						},
						"secondary_private_ip_address_count": {
							Ignore: true,
						},
					},
					ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
						logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
						return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					},
				},
			}
			callbacks = append(callbacks, callback)
		}
	}
	// 检查secondary_private_ip_address_count改变
	if resourceData.HasChange("secondary_private_ip_address_count") {
		privateIpAddress := resourceData.Get("private_ip_address").(*schema.Set).List()
		oldCount, newCount := resourceData.GetChange("secondary_private_ip_address_count")
		if oldCount != nil && newCount != nil && newCount != len(privateIpAddress) {
			diff := newCount.(int) - oldCount.(int)
			if diff > 0 {
				callback = bp.Callback{
					Call: bp.SdkCall{
						Action:      "AssignPrivateIpAddresses",
						ConvertMode: bp.RequestConvertInConvert,
						BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
							(*call.SdkParam)["NetworkInterfaceId"] = d.Id()
							(*call.SdkParam)["SecondaryPrivateIpAddressCount"] = diff
							return true, nil
						},
						Convert: map[string]bp.RequestConvert{
							"private_ip_address": {
								Ignore: true,
							},
							"secondary_private_ip_address_count": {
								Ignore: true,
							},
						},
						ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
							logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
							return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
						},
					},
				}
				callbacks = append(callbacks, callback)
			} else {
				diff *= -1
				removeIpAddress := privateIpAddress[:diff]
				callback = bp.Callback{
					Call: bp.SdkCall{
						Action:      "UnassignPrivateIpAddresses",
						ConvertMode: bp.RequestConvertInConvert,
						BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
							(*call.SdkParam)["NetworkInterfaceId"] = d.Id()
							for index, r := range removeIpAddress {
								(*call.SdkParam)["PrivateIpAddress."+strconv.Itoa(index+1)] = r
							}
							return true, nil
						},
						Convert: map[string]bp.RequestConvert{
							"private_ip_address": {
								Ignore: true,
							},
							"secondary_private_ip_address_count": {
								Ignore: true,
							},
						},
						ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
							logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
							return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
						},
					},
				}
				callbacks = append(callbacks, callback)
			}
		}
	}

	// 更新Tags
	setResourceTagsCallbacks := bp.SetResourceTags(s.Client, "TagResources", "UntagResources", "eni", resourceData, getUniversalInfo)
	callbacks = append(callbacks, setResourceTagsCallbacks...)

	return callbacks
}

func (s *VestackNetworkInterfaceService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteNetworkInterface",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"NetworkInterfaceId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.VpcClient.DeleteNetworkInterfaceCommon(call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading network interface on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 3*time.Minute)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackNetworkInterfaceService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "NetworkInterfaceIds",
				ConvertType: bp.ConvertWithN,
			},
			"primary_ip_addresses": {
				TargetField: "PrimaryIpAddresses",
				ConvertType: bp.ConvertWithN,
			},
			"private_ip_addresses": {
				TargetField: "PrivateIpAddresses",
				ConvertType: bp.ConvertWithN,
			},
			"network_interface_ids": {
				TargetField: "NetworkInterfaceIds",
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
		NameField:    "NetworkInterfaceName",
		IdField:      "NetworkInterfaceId",
		CollectField: "network_interfaces",
		ResponseConverts: map[string]bp.ResponseConvert{
			"NetworkInterfaceId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"AssociatedElasticIp.AllocationId": {
				TargetField: "associated_elastic_ip_id",
			},
			"AssociatedElasticIp.EipAddress": {
				TargetField: "associated_elastic_ip_address",
			},
		},
	}
}

func (s *VestackNetworkInterfaceService) ReadResourceId(id string) string {
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

func (s *VestackNetworkInterfaceService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "vpc",
		ResourceType:         "eni",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}
