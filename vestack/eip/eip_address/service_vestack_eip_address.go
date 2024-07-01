package eip_address

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackEipAddressService struct {
	Client *bp.SdkClient
}

func NewEipAddressService(c *bp.SdkClient) *VestackEipAddressService {
	return &VestackEipAddressService{
		Client: c,
	}
}

func (s *VestackEipAddressService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackEipAddressService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		vpcClient := s.Client.VpcClient
		action := "DescribeEipAddresses"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = vpcClient.DescribeEipAddressesCommon(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = vpcClient.DescribeEipAddressesCommon(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = bp.ObtainSdkValue("Result.EipAddresses", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.EipAddresses is not Slice")
		}
		data, err = RemoveSystemTags(data)
		return data, err
	})
}

func (s *VestackEipAddressService) ReadResource(resourceData *schema.ResourceData, allocationId string) (data map[string]interface{}, err error) {
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
			return data, errors.New("value is not map")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("eip address %s not exist ", allocationId)
	}
	return data, err
}

func (s *VestackEipAddressService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      5 * time.Second,
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
					return nil, "", fmt.Errorf("eip address status error, status:%s", status.(string))
				}
			}
			project, err := bp.ObtainSdkValue("ProjectName", demo)
			if err != nil {
				return nil, "", err
			}
			if resourceData.Get("project_name") != nil && resourceData.Get("project_name").(string) != "" {
				if project != resourceData.Get("project_name") {
					return demo, "", err
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}
}

func (VestackEipAddressService) WithResourceResponseHandlers(eip map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return eip, map[string]bp.ResponseConvert{
			"BillingType": {
				TargetField: "billing_type",
				Convert:     billingTypeResponseConvert,
			},
			"ISP": {
				TargetField: "isp",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackEipAddressService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AllocateEipAddress",
			ConvertMode: bp.RequestConvertAll,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				//periodUnit, ok := (*call.SdkParam)["PeriodUnit"]
				//if !ok {
				//	return true, nil
				//}
				//(*call.SdkParam)["PeriodUnit"] = periodUnitRequestConvert(periodUnit)

				// PeriodUnit 默认传 1(Month)
				if (*call.SdkParam)["BillingType"] == 1 {
					(*call.SdkParam)["PeriodUnit"] = 1
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.VpcClient.AllocateEipAddressCommon(call.SdkParam)
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
			Convert: map[string]bp.RequestConvert{
				"billing_type": {
					TargetField: "BillingType",
					Convert:     billingTypeRequestConvert,
				},
				"isp": {
					TargetField: "ISP",
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

func (s *VestackEipAddressService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyEipAddressAttributes",
			ConvertMode: bp.RequestConvertAll,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) > 0 {
					(*call.SdkParam)["AllocationId"] = d.Id()
					delete(*call.SdkParam, "Tags")
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.VpcClient.ModifyEipAddressAttributesCommon(call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available", "Attached"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
			Convert: map[string]bp.RequestConvert{
				"billing_type": {
					Ignore: true,
				},
				"isp": {
					Ignore: true,
				},
			},
		},
	}

	callbacks = append(callbacks, callback)

	if resourceData.HasChange("billing_type") {
		chargeTypeCall := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ConvertEipAddressBillingType",
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"billing_type": {
						Convert: billingTypeRequestConvert,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) > 0 {
						(*call.SdkParam)["AllocationId"] = d.Id()
						if (*call.SdkParam)["BillingType"] == 1 {
							//periodUnit, ok := d.GetOk("period_unit")
							//if !ok {
							//	return false, fmt.Errorf("PeriodUnit is not exist")
							//}
							//(*call.SdkParam)["PeriodUnit"] = periodUnitRequestConvert(periodUnit)

							// PeriodUnit 默认传 1(Month)
							(*call.SdkParam)["PeriodUnit"] = 1
							(*call.SdkParam)["Period"] = d.Get("period")
						}
						return true, nil
					}
					return false, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					if d.Get("billing_type").(string) != "PrePaid" {
						_ = d.Set("period", nil)
						//d.Set("period_unit", nil)
					}
					return nil
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Available", "Attached"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, chargeTypeCall)
	}

	// 更新Tags
	setResourceTagsCallbacks := bp.SetResourceTags(s.Client, "TagResources", "UntagResources", "eip", resourceData, getUniversalInfo)
	callbacks = append(callbacks, setResourceTagsCallbacks...)

	return callbacks
}

func (s *VestackEipAddressService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ReleaseEipAddress",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"AllocationId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.VpcClient.ReleaseEipAddressCommon(call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading eip address on delete %q, %w", d.Id(), callErr))
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

func (s *VestackEipAddressService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "AllocationIds",
				ConvertType: bp.ConvertWithN,
			},
			"eip_addresses": {
				TargetField: "EipAddresses",
				ConvertType: bp.ConvertWithN,
			},
			"isp": {
				TargetField: "ISP",
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
		NameField:    "Name",
		IdField:      "AllocationId",
		CollectField: "addresses",
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

func (s *VestackEipAddressService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vpc",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}

func (s *VestackEipAddressService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "vpc",
		ResourceType:         "eip",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}

func (s *VestackEipAddressService) UnsubscribeInfo(resourceData *schema.ResourceData, resource *schema.Resource) (*bp.UnsubscribeInfo, error) {
	info := bp.UnsubscribeInfo{
		InstanceId: s.ReadResourceId(resourceData.Id()),
	}
	if resourceData.Get("billing_type") == "PrePaid" {
		info.Products = []string{"EIP"}
		info.NeedUnsubscribe = true
	}
	return &info, nil
}
