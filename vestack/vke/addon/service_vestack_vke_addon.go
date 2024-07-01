package addon

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackVkeAddonService struct {
	Client *bp.SdkClient
}

func NewVkeAddonService(c *bp.SdkClient) *VestackVkeAddonService {
	return &VestackVkeAddonService{
		Client: c,
	}
}

func (s *VestackVkeAddonService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackVkeAddonService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	data, err = bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		action := "ListAddons"
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

		results, err = bp.ObtainSdkValue("Result.Items", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Items is not Slice")
		}
		return data, err
	})
	return data, err
}

func (s *VestackVkeAddonService) ReadResource(resourceData *schema.ResourceData, resourceId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if resourceId == "" {
		resourceId = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(resourceId, ":")
	if len(ids) != 2 {
		return map[string]interface{}{}, fmt.Errorf("invalid addon id")
	}

	clusterId := ids[0]
	name := ids[1]

	req := map[string]interface{}{
		"Filter": map[string]interface{}{
			"ClusterIds": []string{clusterId},
			"Names":      []string{name},
		},
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
		return data, fmt.Errorf("Vke Addon %s:%s not exist ", clusterId, name)
	}
	if cfg, ok := resourceData.GetOkExists("config"); ok {
		// 返回的 config 可能会添加默认参数，这里始终使用创建的
		data["Config"] = cfg
	}
	return data, err
}

func (s *VestackVkeAddonService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	failStateTimes := 0 // 处于Failed状态的组件可能经过短暂的时间（自愈）状态变成 Running
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				demo       map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "Failed")
			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status.Phase", demo)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					failStateTimes++
					if failStateTimes > 10 { // 硬编码：检查10次还是错误状态
						return nil, "", fmt.Errorf(" Vke addon status error, status:%s", status.(string))
					}
					return demo, "", nil
				}
			}
			return demo, status.(string), err
		},
	}

}

func (VestackVkeAddonService) WithResourceResponseHandlers(data map[string]interface{}) []bp.ResourceResponseHandler {
	return []bp.ResourceResponseHandler{}
}

func (s *VestackVkeAddonService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateAddon",
			ContentType: bp.ContentTypeJson,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id := fmt.Sprintf("%s:%s", d.Get("cluster_id"), d.Get("name"))
				d.SetId(id)
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Running"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}

}

func (s *VestackVkeAddonService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateAddonConfig",
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				databaseId := d.Id()
				ids := strings.Split(databaseId, ":")
				if len(ids) != 2 {
					return false, fmt.Errorf("invalid addon id")
				}
				(*call.SdkParam)["ClusterId"] = ids[0]
				(*call.SdkParam)["Name"] = ids[1]

				if ids[1] == "ingress-nginx" {
					return false, fmt.Errorf("ingress-nginx addon prohibits updating config")
				}

				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				req, err := json.Marshal(*call.SdkParam)
				if err != nil {
					return nil, err
				}
				logger.Debug(logger.ReqFormat, call.Action, string(req))
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Running"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackVkeAddonService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteAddon",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				databaseId := d.Id()
				ids := strings.Split(databaseId, ":")
				if len(ids) != 2 {
					return false, fmt.Errorf("invalid addon id")
				}
				(*call.SdkParam)["ClusterId"] = ids[0]
				(*call.SdkParam)["Name"] = ids[1]
				(*call.SdkParam)["CascadingDeleteResources"] = []string{"Crd"}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 10*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				if strings.Contains(baseErr.Error(), "forbidden to delete required addon") { // 一些组件禁止删除，直接返回
					msg := fmt.Sprintf("error: %s. msg: %s",
						baseErr.Error(),
						"If you want to remove it form terraform state, "+
							"please use `terraform state rm vestack_vke_addon.resource_name` command ")
					return fmt.Errorf(msg)
				}
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading addon on delete %q, %w", d.Id(), callErr))
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

func (s *VestackVkeAddonService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"cluster_ids": {
				TargetField: "Filter.ClusterIds",
				ConvertType: bp.ConvertJsonArray,
			},
			"names": {
				TargetField: "Filter.Names",
				ConvertType: bp.ConvertJsonArray,
			},
			"deploy_mode": {
				TargetField: "Filter.DeployModes",
				ConvertType: bp.ConvertJsonArray,
			},
			"deploy_node_types": {
				TargetField: "Filter.DeployNodeTypes",
				ConvertType: bp.ConvertJsonArray,
			},
			"statuses": {
				TargetField: "Filter.Statuses",
				ConvertType: bp.ConvertJsonObjectArray,
				NextLevelConvert: map[string]bp.RequestConvert{
					"phase": {
						TargetField: "Phase",
					},
					"conditions_type": {
						TargetField: "ConditionsType",
					},
				},
			},
			"create_client_token": {
				TargetField: "Filter.CreateClientToken",
			},
			"update_client_token": {
				TargetField: "Filter.UpdateClientToken",
			},
		},
		ContentType:  bp.ContentTypeJson,
		NameField:    "Name",
		CollectField: "addons",
	}
}

func (s *VestackVkeAddonService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vke",
		Version:     "2022-05-12",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
