package server_group

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ve "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/clb"
)

type VestackServerGroupService struct {
	Client *ve.SdkClient
}

func NewServerGroupService(c *ve.SdkClient) *VestackServerGroupService {
	return &VestackServerGroupService{
		Client: c,
	}
}

func (s *VestackServerGroupService) GetClient() *ve.SdkClient {
	return s.Client
}

func (s *VestackServerGroupService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return ve.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		action := "DescribeServerGroups"
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

		results, err = ve.ObtainSdkValue("Result.ServerGroups", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.ServerGroups is not Slice")
		}
		for index, serverGroup := range data {
			if serverGroupMap, ok := serverGroup.(map[string]interface{}); ok {
				id := serverGroupMap["ServerGroupId"].(string)
				clbId, err := s.queryLoadBalancerId(id)
				if err != nil {
					return data, err
				}
				serverGroupMap["LoadBalancerId"] = clbId
				data[index] = serverGroupMap
			}
		}
		return data, err
	})
}

func (s *VestackServerGroupService) ReadResource(resourceData *schema.ResourceData, serverGroupId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if serverGroupId == "" {
		serverGroupId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"ServerGroupIds.1": serverGroupId,
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
		return data, fmt.Errorf("ServerGroup %s not exist ", serverGroupId)
	}
	return data, err
}

func (s *VestackServerGroupService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (VestackServerGroupService) WithResourceResponseHandlers(serverGroup map[string]interface{}) []ve.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]ve.ResponseConvert, error) {
		return serverGroup, nil, nil
	}
	return []ve.ResourceResponseHandler{handler}

}

func (s *VestackServerGroupService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []ve.Callback {
	callback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "CreateServerGroup",
			ConvertMode: ve.RequestConvertAll,
			Convert: map[string]ve.RequestConvert{
				"servers": {
					ConvertType: ve.ConvertListN,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				// 创建 server group
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *ve.SdkClient, resp *map[string]interface{}, call ve.SdkCall) error {
				// 注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := ve.ObtainSdkValue("Result.ServerGroupId", *resp)
				d.SetId(id.(string))
				return nil
			},
			ExtraRefresh: map[ve.ResourceService]*ve.StateRefresh{
				clb.NewClbService(s.Client): {
					Target:     []string{"Active", "Inactive"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("load_balancer_id").(string),
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return resourceData.Get("load_balancer_id").(string)
			},
		},
	}
	return []ve.Callback{callback}

}

func (s *VestackServerGroupService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []ve.Callback {
	clbId, err := s.queryLoadBalancerId(resourceData.Id())
	if err != nil {
		return []ve.Callback{{
			Err: err,
		}}
	}

	callback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "ModifyServerGroupAttributes",
			ConvertMode: ve.RequestConvertAll,
			BeforeCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (bool, error) {
				(*call.SdkParam)["ServerGroupId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				// 修改 server group 属性
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			ExtraRefresh: map[ve.ResourceService]*ve.StateRefresh{
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
	return []ve.Callback{callback}
}

func (s *VestackServerGroupService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []ve.Callback {
	clbId, err := s.queryLoadBalancerId(resourceData.Id())
	if err != nil {
		return []ve.Callback{{
			Err: err,
		}}
	}

	callback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "DeleteServerGroup",
			ConvertMode: ve.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"ServerGroupId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//删除 Server Group
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if ve.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading server group on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			ExtraRefresh: map[ve.ResourceService]*ve.StateRefresh{
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
	return []ve.Callback{callback}
}

func (s *VestackServerGroupService) DatasourceResources(*schema.ResourceData, *schema.Resource) ve.DataSourceInfo {
	return ve.DataSourceInfo{
		RequestConverts: map[string]ve.RequestConvert{
			"ids": {
				TargetField: "ServerGroupIds",
				ConvertType: ve.ConvertWithN,
			},
		},
		NameField:    "ServerGroupName",
		IdField:      "ServerGroupId",
		CollectField: "groups",
		ResponseConverts: map[string]ve.ResponseConvert{
			"ServerGroupId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *VestackServerGroupService) ReadResourceId(id string) string {
	return id
}

func (s *VestackServerGroupService) queryLoadBalancerId(serverGroupId string) (string, error) {
	if serverGroupId == "" {
		return "", fmt.Errorf("server_group_id cannot be empty")
	}

	// 查询 LoadBalancerId
	action := "DescribeServerGroupAttributes"
	serverGroupResp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &map[string]interface{}{
		"ServerGroupId": serverGroupId,
	})
	if err != nil {
		return "", err
	}
	clbId, err := ve.ObtainSdkValue("Result.LoadBalancerId", *serverGroupResp)
	if err != nil {
		return "", err
	}
	return clbId.(string), nil
}

func getUniversalInfo(actionName string) ve.UniversalInfo {
	return ve.UniversalInfo{
		ServiceName: "clb",
		Version:     "2020-04-01",
		HttpMethod:  ve.GET,
		ContentType: ve.Default,
		Action:      actionName,
	}
}
