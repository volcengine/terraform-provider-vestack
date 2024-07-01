package default_node_pool

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/node"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/node_pool"
)

type VestackDefaultNodePoolService struct {
	Client          *bp.SdkClient
	nodePoolService *node_pool.VestackNodePoolService
	nodeService     *node.VestackVkeNodeService
}

func NewDefaultNodePoolService(c *bp.SdkClient) *VestackDefaultNodePoolService {
	return &VestackDefaultNodePoolService{
		Client:          c,
		nodePoolService: node_pool.NewNodePoolService(c),
		nodeService:     node.NewVestackVkeNodeService(c),
	}
}

const (
	DefaultNodePoolName = "vke-default-nodepool"
)

func (s *VestackDefaultNodePoolService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackDefaultNodePoolService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return data, err
}

func (s *VestackDefaultNodePoolService) ReadResource(resourceData *schema.ResourceData, nodePoolId string) (data map[string]interface{}, err error) {
	var (
		nodes []interface{}
	)
	if nodePoolId == "" {
		nodePoolId = s.ReadResourceId(resourceData.Id())
	}
	clusterId := ""
	if c, ok := resourceData.Get("cluster_id").(string); ok {
		clusterId = c
	} else {
		err = errors.New("cluster_id is empty")
		return nil, err
	}
	data, err = s.nodePoolService.ReadResource(resourceData, nodePoolId)
	if err != nil {
		return data, err
	}

	// 只能导入默认节点池，不是默认节点池直接报错
	if data["Name"].(string) != DefaultNodePoolName {
		return nil, fmt.Errorf("only the default node pool is supported")
	}

	data["NodeConfig"].(map[string]interface{})["Security"].(map[string]interface{})["Login"].(map[string]interface{})["Password"] =
		resourceData.Get("node_config.0.security.0.login.0.password")

	nodes, err = s.nodeService.ReadResources(map[string]interface{}{
		"Filter": map[string]interface{}{
			"NodePoolIds": []string{nodePoolId},
			"ClusterIds":  []string{clusterId},
		},
	})
	if err != nil {
		return data, err
	}

	instanceMap := make(map[string]string)
	instances := resourceData.Get("instances").(*schema.Set)
	var ins []interface{}
	if resourceData.Get("is_import").(bool) {
		// 第一次导入的时候由于读不到tf文件的数据，在这里默认导入所有的节点
		for _, n := range nodes {
			n.(map[string]interface{})["Phase"], _ = bp.ObtainSdkValue("Status.Phase", n)
			ins = append(ins, n)
		}
		_ = resourceData.Set("is_import", false)
	} else {
		for _, ins := range instances.List() {
			instancesId, _ := bp.ObtainSdkValue("instance_id", ins)
			imageId, _ := bp.ObtainSdkValue("image_id", ins)
			instanceMap[instancesId.(string)] = imageId.(string)
		}
		for _, n := range nodes {
			instancesId, _ := bp.ObtainSdkValue("InstanceId", n)
			if v, ok := instanceMap[instancesId.(string)]; ok {
				if v == "" {
					n.(map[string]interface{})["ImageId"] = ""
				}
				n.(map[string]interface{})["Phase"], _ = bp.ObtainSdkValue("Status.Phase", n)
				ins = append(ins, n)
			}
		}
	}

	data["Instances"] = ins
	logger.Debug(logger.RespFormat, "VestackDefaultNodePoolService result ", data)

	return data, err
}

func (s *VestackDefaultNodePoolService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				nodes      []interface{}
				status     interface{}
				statuses   []string
				failStates []string
			)
			instanceMap := make(map[string]bool)
			instances := resourceData.Get("instances").(*schema.Set)
			clusterId := ""
			if c, ok := resourceData.Get("cluster_id").(string); ok {
				clusterId = c
			} else {
				err = errors.New("cluster_id is empty")
				return nil, "", err
			}
			for _, ins := range instances.List() {
				instancesId, _ := bp.ObtainSdkValue("instance_id", ins)
				instanceMap[instancesId.(string)] = true
			}

			failStates = []string{
				"Failed+InitializeFailed",
				"Failed+Unknown",
				"Failed+ResourceCleanupFailed",
			}
			nodes, err = s.nodeService.ReadResources(map[string]interface{}{
				"Filter": map[string]interface{}{
					"NodePoolIds": []string{id},
					"ClusterIds":  []string{clusterId},
				},
			})
			if err != nil {
				return nil, "", err
			}
			for _, n := range nodes {
				var (
					instancesId interface{}
					sts         interface{}
					conditions  interface{}
					ss          []string
				)
				instancesId, err = bp.ObtainSdkValue("InstanceId", n)
				if err != nil {
					return nil, "", err
				}
				sts, err = bp.ObtainSdkValue("Status.Phase", n)
				if err != nil {
					return nil, "", err
				}
				conditions, _ = bp.ObtainSdkValue("Status.Conditions", n)
				if err != nil {
					return nil, "", err
				}
				ss, err = vke.BinaryJudgment(sts.(string), conditions.([]interface{}), []string{"Failed"})
				if err != nil {
					return nil, "", err
				}
				if _, ok := instanceMap[instancesId.(string)]; ok {
					statuses = append(statuses, ss...)
				}
			}
			for _, v := range failStates {
				for _, status = range statuses {
					if v == status {
						return nil, "", fmt.Errorf("node status error, status:%s", status.(string))
					}
				}
			}

			for _, v := range statuses {
				if v != "Running" && v != "Stopped" && v != "Failed" {
					return nodes, v, err
				}
			}

			return nodes, "Running", err
		},
	}
}

func (s *VestackDefaultNodePoolService) WithResourceResponseHandlers(nodePool map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		var (
			security = make([]interface{}, 0)
			login    = make([]interface{}, 0)
		)

		priSecurity := nodePool["NodeConfig"].(map[string]interface{})["Security"]
		priLogin := priSecurity.(map[string]interface{})["Login"]
		delete(nodePool, "Login")
		login = append(login, priLogin)
		priSecurity.(map[string]interface{})["Login"] = login
		security = append(security, priSecurity)

		delete(nodePool, "Security")
		nodePool["NodeConfig"].(map[string]interface{})["Security"] = security

		return nodePool, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackDefaultNodePoolService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var calls []bp.Callback
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateDefaultNodePool",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"instances": {
					Ignore: true,
				},
				"kubernetes_config": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"labels": {
							ConvertType: bp.ConvertJsonArray,
						},
						"taints": {
							ConvertType: bp.ConvertJsonArray,
						},
					},
				},
				"node_config": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"security": {
							ConvertType: bp.ConvertJsonObject,
							NextLevelConvert: map[string]bp.RequestConvert{
								"login": {
									ConvertType: bp.ConvertJsonObject,
								},
								"security_group_ids": {
									ConvertType: bp.ConvertJsonArray,
								},
								"security_strategies": {
									ConvertType: bp.ConvertJsonArray,
								},
							},
						},
						"name_prefix": {
							ConvertType: bp.ConvertJsonObject,
						},
						"ecs_tags": {
							TargetField: "Tags",
							ConvertType: bp.ConvertJsonObjectArray,
						},
					},
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertJsonObjectArray,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				(*call.SdkParam)["ClientToken"] = "default-nodeService-pool-" + uuid.New().String()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.Id", *resp)
				d.SetId(id.(string))
				return nil
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				node_pool.NewNodePoolService(s.Client): {
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("cluster_id").(string)
			},
		},
	}

	calls = append(calls, callback)

	if _, ok := resourceData.GetOk("instances"); ok {
		calls = s.ProcessNodeInstances(resourceData, calls)
	}

	return calls
}

func (s *VestackDefaultNodePoolService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var calls []bp.Callback
	// 先修改节点池配置
	calls = append(calls, bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateNodePoolConfig",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"node_config": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"security": {
							ConvertType: bp.ConvertJsonObject,
							NextLevelConvert: map[string]bp.RequestConvert{
								"login": {
									ConvertType: bp.ConvertJsonObject,
								},
								"security_group_ids": {
									ConvertType: bp.ConvertJsonArray,
								},
								"security_strategies": {
									ConvertType: bp.ConvertJsonArray,
								},
							},
						},
						"initialize_script": {
							ConvertType: bp.ConvertJsonObject,
						},
						"name_prefix": {
							ConvertType: bp.ConvertJsonObject,
						},
						"ecs_tags": {
							TargetField: "Tags",
							ConvertType: bp.ConvertJsonObjectArray,
						},
					},
				},
				"kubernetes_config": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"labels": {
							ConvertType: bp.ConvertJsonObjectArray,
							ForceGet:    true,
						},
						"taints": {
							ConvertType: bp.ConvertJsonObjectArray,
							ForceGet:    true,
						},
						"cordon": {
							ConvertType: bp.ConvertJsonObject,
						},
					},
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Id"] = d.Id()
				(*call.SdkParam)["ClusterId"] = d.Get("cluster_id")

				delete(*call.SdkParam, "Tags")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				//adapt vke api
				nodeconfig := (*call.SdkParam)["NodeConfig"]
				if nodeconfig != nil {
					security := nodeconfig.(map[string]interface{})["Security"]
					if security != nil {
						login := security.(map[string]interface{})["Login"]
						if login != nil && login.(map[string]interface{})["SshKeyPairName"] != nil && login.(map[string]interface{})["SshKeyPairName"].(string) == "" {
							delete((*call.SdkParam)["NodeConfig"].(map[string]interface{})["Security"].(map[string]interface{})["Login"].(map[string]interface{}), "SshKeyPairName")
						}
						_, ok := security.(map[string]interface{})["SecurityStrategies"]
						if !ok && d.HasChange("node_config.0.security.0.security_strategies") {
							(*call.SdkParam)["NodeConfig"].(map[string]interface{})["Security"].(map[string]interface{})["SecurityStrategies"] = []interface{}{}
						}
					}
				}

				// 当列表被删除时，入参添加空列表来置空
				bp.DefaultMapValue(call.SdkParam, "KubernetesConfig", map[string]interface{}{
					"Labels": []interface{}{},
					"Taints": []interface{}{},
				})

				if d.HasChange("node_config.0.ecs_tags") {
					bp.DefaultMapValue(call.SdkParam, "NodeConfig", map[string]interface{}{
						"Tags": []interface{}{},
					})
				}

				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Running"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	})

	// 更新Tags
	calls = s.setResourceTags(resourceData, "NodePool", calls)

	//修改实例
	if resourceData.HasChange("instances") {
		calls = s.ProcessNodeInstances(resourceData, calls)
	}
	return calls
}

func (s *VestackDefaultNodePoolService) RemoveResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteNodePool",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"Id":        resourceData.Id(),
				"ClusterId": resourceData.Get("cluster_id"),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
		},
	}
	return []bp.Callback{callback}
}

func (VestackDefaultNodePoolService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackDefaultNodePoolService) ReadResourceId(id string) string {
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

func (s *VestackDefaultNodePoolService) ProcessNodeInstances(resourceData *schema.ResourceData, calls []bp.Callback) []bp.Callback {
	add, remove, _, _ := bp.GetSetDifference("instances", resourceData, defaultNodePoolNodeHash, false)
	logger.Debug(logger.RespFormat, "processNodeInstancesAdd", add)
	logger.Debug(logger.RespFormat, "processNodeInstancesRemove", remove)
	newNode := make(map[string][]string)
	var delNode []string
	if add != nil {
		for _, v := range add.List() {
			m := v.(map[string]interface{})
			if m["instance_id"] == nil || len(m["instance_id"].(string)) == 0 {
				continue
			}
			key := strconv.FormatBool(m["keep_instance_name"].(bool)) + ":" + strconv.FormatBool(m["additional_container_storage_enabled"].(bool)) + ":" +
				m["image_id"].(string) + ":" + m["container_storage_path"].(string)
			if _, ok1 := newNode[key]; !ok1 {
				newNode[key] = []string{}
			}
			newNode[key] = append(newNode[key], m["instance_id"].(string))
		}
	}
	if remove != nil {
		for _, v := range remove.List() {
			m := v.(map[string]interface{})
			delNode = append(delNode, m["id"].(string))
		}
	}

	// 删除节点
	for i := 0; i < len(delNode)/100+1; i++ {
		start := i * 100
		end := (i + 1) * 100
		if end > len(delNode) {
			end = len(delNode)
		}
		if end <= start {
			break
		}
		calls = append(calls, func(nodeIds []string, clusterId, nodePoolId string) bp.Callback {
			return bp.Callback{
				Call: bp.SdkCall{
					Action:      "DeleteNodes",
					ConvertMode: bp.RequestConvertIgnore,
					ContentType: bp.ContentTypeJson,
					SdkParam: &map[string]interface{}{
						"ClusterId":  clusterId,
						"NodePoolId": nodePoolId,
					},
					BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
						if len(nodeIds) < 1 {
							return false, nil
						}
						for index, id := range nodeIds {
							(*call.SdkParam)[fmt.Sprintf("Ids.%d", index+1)] = id
						}
						return true, nil
					},
					ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
						logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
						resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
						logger.Debug(logger.RespFormat, call.Action, resp, err)
						return resp, err
					},
					Refresh: &bp.StateRefresh{
						Target:  []string{"Running"},
						Timeout: resourceData.Timeout(schema.TimeoutCreate),
					},
					ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
						node_pool.NewNodePoolService(s.Client): {
							Target:  []string{"Running"},
							Timeout: resourceData.Timeout(schema.TimeoutCreate),
						},
					},
					LockId: func(d *schema.ResourceData) string {
						return d.Get("cluster_id").(string)
					},
				},
			}
		}(delNode[start:end], resourceData.Get("cluster_id").(string), resourceData.Id()))
	}

	// 新增加节点
	for k, v := range newNode {
		nodeCall := bp.Callback{
			Call: bp.SdkCall{
				Action:      "CreateNodes",
				ConvertMode: bp.RequestConvertIgnore,
				ContentType: bp.ContentTypeJson,
				SdkParam: &map[string]interface{}{
					"Key":   k,
					"Value": v,
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["ClientToken"] = "default-nodeService-pool-" + uuid.New().String()
					(*call.SdkParam)["ClusterId"] = d.Get("cluster_id")
					if (*call.SdkParam)["Value"] == nil || len((*call.SdkParam)["Value"].([]string)) == 0 {
						return false, nil
					}
					for i, v1 := range (*call.SdkParam)["Value"].([]string) {
						(*call.SdkParam)["InstanceIds."+strconv.Itoa(i+1)] = v1
					}
					ks := strings.Split((*call.SdkParam)["Key"].(string), ":")
					keepInstanceName, _ := strconv.ParseBool(ks[0])
					additionalContainerStorageEnabled, _ := strconv.ParseBool(ks[1])
					(*call.SdkParam)["KeepInstanceName"] = keepInstanceName
					(*call.SdkParam)["AdditionalContainerStorageEnabled"] = additionalContainerStorageEnabled
					if ks[2] != "" {
						(*call.SdkParam)["ImageId"] = ks[2]
					}
					if len(ks) == 4 && ks[3] != "" {
						(*call.SdkParam)["ContainerStoragePath"] = ks[3]
					} else if len(ks) > 4 {
						(*call.SdkParam)["ContainerStoragePath"] = (*call.SdkParam)["Key"].(string)[len(ks[0])+1+len(ks[1])+1+len(ks[2])+1:]
					}
					delete(*call.SdkParam, "Key")
					delete(*call.SdkParam, "Value")
					kubernetesConfig := vke.TransKubernetesConfig(d)
					if kubernetesConfig != nil {
						(*call.SdkParam)["KubernetesConfig"] = kubernetesConfig
					}
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
				ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
					node_pool.NewNodePoolService(s.Client): {
						Target:  []string{"Running"},
						Timeout: resourceData.Timeout(schema.TimeoutCreate),
					},
				},
				LockId: func(d *schema.ResourceData) string {
					return d.Get("cluster_id").(string)
				},
			},
		}
		calls = append(calls, nodeCall)
	}
	return calls
}

func (s *VestackDefaultNodePoolService) setResourceTags(resourceData *schema.ResourceData, resourceType string, callbacks []bp.Callback) []bp.Callback {
	addedTags, removedTags, _, _ := bp.GetSetDifference("tags", resourceData, bp.TagsHash, false)

	removeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UntagResources",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if removedTags != nil && len(removedTags.List()) > 0 {
					(*call.SdkParam)["ResourceIds"] = []string{resourceData.Id()}
					(*call.SdkParam)["ResourceType"] = resourceType
					(*call.SdkParam)["TagKeys"] = make([]string, 0)
					for _, tag := range removedTags.List() {
						(*call.SdkParam)["TagKeys"] = append((*call.SdkParam)["TagKeys"].([]string), tag.(map[string]interface{})["key"].(string))
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, removeCallback)

	addCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "TagResources",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if addedTags != nil && len(addedTags.List()) > 0 {
					(*call.SdkParam)["ResourceIds"] = []string{resourceData.Id()}
					(*call.SdkParam)["ResourceType"] = resourceType
					(*call.SdkParam)["Tags"] = make([]map[string]interface{}, 0)
					for _, tag := range addedTags.List() {
						(*call.SdkParam)["Tags"] = append((*call.SdkParam)["Tags"].([]map[string]interface{}), tag.(map[string]interface{}))
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, addCallback)

	return callbacks
}
