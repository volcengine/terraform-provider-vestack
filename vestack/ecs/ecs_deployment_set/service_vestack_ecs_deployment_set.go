package ecs_deployment_set

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackEcsDeploymentSetService struct {
	Client *bp.SdkClient
}

func NewEcsDeploymentSetService(c *bp.SdkClient) *VestackEcsDeploymentSetService {
	return &VestackEcsDeploymentSetService{
		Client: c,
	}
}

func (s *VestackEcsDeploymentSetService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackEcsDeploymentSetService) ReadResources(condition map[string]interface{}) ([]interface{}, error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithNextTokenQuery(condition, "MaxResults", "NextToken", 20, nil, func(m map[string]interface{}) (data []interface{}, next string, err error) {
		client := s.Client.UniversalClient
		action := "DescribeDeploymentSets"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = client.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, next, err
			}
		} else {
			resp, err = client.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, next, err
			}
		}
		logger.Debug(logger.RespFormat, action, condition, *resp)

		results, err = bp.ObtainSdkValue("Result.DeploymentSets", *resp)
		if err != nil {
			return data, next, err
		}
		nextToken, err := bp.ObtainSdkValue("Result.NextToken", *resp)
		if err != nil {
			return data, next, err
		}
		next = nextToken.(string)
		if results == nil {
			results = []interface{}{}
		}

		if data, ok = results.([]interface{}); !ok {
			return data, next, errors.New("Result.DeploymentSets is not Slice")
		}
		return data, next, err
	})
}

func (s *VestackEcsDeploymentSetService) ReadResource(resourceData *schema.ResourceData, instanceId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if instanceId == "" {
		instanceId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"DeploymentSetIds.1": instanceId,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, fmt.Errorf("Value is not map ")
		}
	}

	if len(data) == 0 {
		return data, fmt.Errorf("Ecs DeploymentSet %s not exist ", instanceId)
	}
	return data, nil
}

func (s *VestackEcsDeploymentSetService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *VestackEcsDeploymentSetService) WithResourceResponseHandlers(deploymentSet map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return deploymentSet, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackEcsDeploymentSetService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateDeploymentSet",
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//创建DeploymentSet
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := bp.ObtainSdkValue("Result.DeploymentSetId", *resp)
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, id)
				d.SetId(id.(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackEcsDeploymentSetService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDeploymentSetAttribute",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"deployment_set_name": {
					ConvertType: bp.ConvertDefault,
				},
				"description": {
					ConvertType: bp.ConvertDefault,
				},
			},
			RequestIdField: "deployment_set_id",
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) > 1 {
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//修改部署集属性
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackEcsDeploymentSetService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteDeploymentSet",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"DeploymentSetId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//删除部署集
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
							return resource.NonRetryableError(fmt.Errorf("error on  reading deployment set on delete %q, %w", d.Id(), callErr))
						}
					}

					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					if strings.Contains(callErr.Error(), "InvalidDeploymentSet.InUse") {
						return resource.NonRetryableError(fmt.Errorf("there are instances in the deployment set, " +
							"please remove the instances in the deployment set before deleting the deployment set"))
					}
					return resource.RetryableError(callErr)
				})
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackEcsDeploymentSetService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "DeploymentSetIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "DeploymentSetName",
		IdField:      "DeploymentSetId",
		CollectField: "deployment_sets",
	}
}

func (VestackEcsDeploymentSetService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "ecs",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}
