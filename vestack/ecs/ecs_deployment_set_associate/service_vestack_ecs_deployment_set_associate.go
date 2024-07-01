package ecs_deployment_set_associate

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

type VestackEcsDeploymentSetAssociateService struct {
	Client *bp.SdkClient
}

func NewEcsDeploymentSetAssociateService(c *bp.SdkClient) *VestackEcsDeploymentSetAssociateService {
	return &VestackEcsDeploymentSetAssociateService{
		Client: c,
	}
}

func (s *VestackEcsDeploymentSetAssociateService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackEcsDeploymentSetAssociateService) ReadResources(condition map[string]interface{}) ([]interface{}, error) {
	return nil, nil
}

func (s *VestackEcsDeploymentSetAssociateService) ReadResource(resourceData *schema.ResourceData, tmpId string) (data map[string]interface{}, err error) {
	var (
		resp             *map[string]interface{}
		results          interface{}
		ok               bool
		deploymentSetId  string
		targetInstanceId string
		ids              []string
		dep              []interface{}
		instanceIds      []interface{}
	)

	if tmpId == "" {
		tmpId = s.ReadResourceId(resourceData.Id())
	}

	ids = strings.Split(tmpId, ":")
	deploymentSetId = ids[0]
	targetInstanceId = ids[1]

	req := map[string]interface{}{
		"DeploymentSetIds.1": deploymentSetId,
	}
	client := s.Client.UniversalClient
	action := "DescribeDeploymentSets"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = client.DoCall(getUniversalInfo(action), &req)
	if err != nil {
		return data, err
	}

	results, err = bp.ObtainSdkValue("Result.DeploymentSets", *resp)
	if err != nil {
		return data, err
	}
	if dep, ok = results.([]interface{}); !ok {
		return data, errors.New("Result.DeploymentSets is not Slice")
	}
	if len(dep) == 0 {
		return data, fmt.Errorf("Ecs DeploymentSet %s not exist ", deploymentSetId)
	}
	results, err = bp.ObtainSdkValue("InstanceIds", dep[0])
	if instanceIds, ok = results.([]interface{}); !ok {
		return data, errors.New("InstanceIds is not Slice")
	}

	for _, id := range instanceIds {
		if id.(string) == targetInstanceId {
			data = make(map[string]interface{})
			data["DeploymentSetId"] = deploymentSetId
			data["InstanceId"] = targetInstanceId
			return data, err
		}
	}

	return data, err
}

func (s *VestackEcsDeploymentSetAssociateService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				data map[string]interface{}
			)

			if err = resource.Retry(3*time.Minute, func() *resource.RetryError {
				data, err = s.ReadResource(resourceData, id)
				if err != nil {
					return resource.NonRetryableError(err)
				}
				if len(data) == 0 {
					return resource.RetryableError(fmt.Errorf("Retry "))
				}
				return nil
			}); err != nil {
				return nil, "error", err
			}
			return data, "success", err
		},
	}
}

func (s *VestackEcsDeploymentSetAssociateService) WithResourceResponseHandlers(deploymentSet map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return deploymentSet, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackEcsDeploymentSetAssociateService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyInstanceDeployment",
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId((*call.SdkParam)["DeploymentSetId"].(string) + ":" + (*call.SdkParam)["InstanceId"].(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"success"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackEcsDeploymentSetAssociateService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackEcsDeploymentSetAssociateService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyInstanceDeployment",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"InstanceId":      strings.Split(resourceData.Id(), ":")[1],
				"DeploymentSetId": "",
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
							return resource.NonRetryableError(fmt.Errorf("error on  reading deployment set associate on delete %q, %w", d.Id(), callErr))
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

func (s *VestackEcsDeploymentSetAssociateService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackEcsDeploymentSetAssociateService) ReadResourceId(id string) string {
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
