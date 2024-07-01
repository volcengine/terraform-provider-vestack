package vpc

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackVpcService struct {
	Client *bp.SdkClient
}

func NewVpcService(c *bp.SdkClient) *VestackVpcService {
	return &VestackVpcService{
		Client: c,
	}
}

func (s *VestackVpcService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackVpcService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		vpcClient := s.Client.VpcClient
		action := "DescribeVpcs"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = vpcClient.DescribeVpcsCommon(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = vpcClient.DescribeVpcsCommon(&condition)
			if err != nil {
				return data, err
			}
		}
		logger.Debug(logger.RespFormat, action, resp)
		results, err = bp.ObtainSdkValue("Result.Vpcs", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Vpcs is not Slice")
		}

		return data, err
	})
}

func (s *VestackVpcService) ReadResource(resourceData *schema.ResourceData, vpcId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if vpcId == "" {
		vpcId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"VpcIds.1": vpcId,
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
		return data, fmt.Errorf("Vpc %s not exist ", vpcId)
	}

	if _, ok1 := data["AuxiliaryCidrBlocks"]; !ok1 {
		data["AuxiliaryCidrBlocks"] = []string{}
	}

	if ipv6CidrBlock, ok2 := data["Ipv6CidrBlock"]; ok2 && ipv6CidrBlock != "" {
		data["EnableIpv6"] = true
	} else {
		data["EnableIpv6"] = false
	}

	return data, err
}

func (s *VestackVpcService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
					return nil, "", fmt.Errorf("Vpc  status  error, status:%s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}

}

func (VestackVpcService) WithResourceResponseHandlers(vpc map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return vpc, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackVpcService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateVpc",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"dns_servers": {
					TargetField: "DnsServers",
					ConvertType: bp.ConvertWithN,
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertListN,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//创建vpc
				return s.Client.VpcClient.CreateVpcCommon(call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := bp.ObtainSdkValue("Result.VpcId", *resp)
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

func (s *VestackVpcService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyVpcAttributes",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"dns_servers": {
					TargetField: "DnsServers",
					ConvertType: bp.ConvertWithN,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["VpcId"] = d.Id()
				delete(*call.SdkParam, "Tags")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//修改vpc属性
				return s.Client.VpcClient.ModifyVpcAttributesCommon(call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, callback)

	// 更新Tags
	setResourceTagsCallbacks := bp.SetResourceTags(s.Client, "TagResources", "UntagResources", "vpc", resourceData, getUniversalInfo)
	callbacks = append(callbacks, setResourceTagsCallbacks...)

	return callbacks
}

func (s *VestackVpcService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteVpc",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"VpcId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//删除VPC
				return s.Client.VpcClient.DeleteVpcCommon(call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading vpc on delete %q, %w", d.Id(), callErr))
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

func (s *VestackVpcService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "VpcIds",
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
		NameField:    "VpcName",
		IdField:      "VpcId",
		CollectField: "vpcs",
		ResponseConverts: map[string]bp.ResponseConvert{
			"VpcId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *VestackVpcService) ReadResourceId(id string) string {
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

func (s *VestackVpcService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "vpc",
		ResourceType:         "vpc",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}
