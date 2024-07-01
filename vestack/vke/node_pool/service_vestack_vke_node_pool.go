package node_pool

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/security_group"
)

type VestackNodePoolService struct {
	Client               *bp.SdkClient
	securityGroupService *security_group.VestackSecurityGroupService
}

func NewNodePoolService(c *bp.SdkClient) *VestackNodePoolService {
	return &VestackNodePoolService{
		Client:               c,
		securityGroupService: security_group.NewSecurityGroupService(c),
	}
}

func (s *VestackNodePoolService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackNodePoolService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)

	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		// adapt vke api
		if enabled, exist := condition["AutoScalingEnabled"]; exist {
			if _, filterExist := condition["Filter"]; !filterExist {
				condition["Filter"] = make(map[string]interface{})
			}
			condition["Filter"].(map[string]interface{})["AutoScaling.Enabled"] = enabled
			delete(condition, "AutoScalingEnabled")
		}

		// 单独适配 ClusterId 字段，将 ClusterId 加入 Filter.ClusterIds
		if filter, filterExist := condition["Filter"]; filterExist {
			if clusterId, clusterIdExist := filter.(map[string]interface{})["ClusterId"]; clusterIdExist {
				if clusterIds, clusterIdsExist := filter.(map[string]interface{})["ClusterIds"]; clusterIdsExist {
					appendFlag := true
					for _, id := range clusterIds.([]interface{}) {
						if id == clusterId {
							appendFlag = false
						}
					}
					if appendFlag {
						condition["Filter"].(map[string]interface{})["ClusterIds"] = append(condition["Filter"].(map[string]interface{})["ClusterIds"].([]interface{}), clusterId)
					}
				} else {
					condition["Filter"].(map[string]interface{})["ClusterIds"] = []interface{}{clusterId}
				}
				delete(condition["Filter"].(map[string]interface{}), "ClusterId")
			}
		}

		action := "ListNodePools"
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

		respBytes, _ := json.Marshal(resp)
		logger.Debug(logger.RespFormat, action, condition, string(respBytes))
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
}

func (s *VestackNodePoolService) ReadResource(resourceData *schema.ResourceData, nodePoolId string) (data map[string]interface{}, err error) {
	var (
		results interface{}
		resp    *map[string]interface{}
		result  map[string]interface{}
		temp    []interface{}
		ok      bool
	)
	if nodePoolId == "" {
		nodePoolId = s.ReadResourceId(resourceData.Id())
	}

	action := "ListNodePools"
	nodeId := []string{nodePoolId}
	condition := make(map[string]interface{}, 0)
	condition["Filter"] = map[string]interface{}{
		"Ids": nodeId,
	}

	logger.Debug(logger.RespFormat, "ReadResource ", condition)
	resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
	logger.Debug(logger.RespFormat, "ReadResource ", resp)

	if err != nil {
		return data, err
	}
	if resp == nil {
		return data, fmt.Errorf("NodePool %s not exist ", nodePoolId)
	}

	results, err = bp.ObtainSdkValue("Result.Items", *resp)
	if err != nil {
		return data, err
	}
	if results == nil {
		results = []interface{}{}
	}

	if temp, ok = results.([]interface{}); !ok {
		return data, errors.New("Result.Items is not Slice")
	}

	if len(temp) == 0 {
		return data, fmt.Errorf("NodePool %s not exist ", nodePoolId)
	}

	result = temp[0].(map[string]interface{})
	result["NodeConfig"].(map[string]interface{})["Security"].(map[string]interface{})["Login"].(map[string]interface{})["Password"] =
		resourceData.Get("node_config.0.security.0.login.0.password")

	// 安全组过滤默认安全组
	tmpSecurityGroupIds := result["NodeConfig"].(map[string]interface{})["Security"].(map[string]interface{})["SecurityGroupIds"].([]interface{})
	if len(tmpSecurityGroupIds) > 0 {
		// 查询安全组
		securityGroupIdMap := make(map[string]interface{})
		for i, securityGroupId := range tmpSecurityGroupIds {
			securityGroupIdMap[fmt.Sprintf("SecurityGroupIds.%d", i+1)] = securityGroupId
		}
		securityGroups, err := s.securityGroupService.ReadResources(securityGroupIdMap)
		logger.Debug(logger.RespFormat, "DescribeSecurityGroups", securityGroupIdMap, securityGroups)
		if err != nil {
			return nil, err
		}

		// 每个节点池有个默认安全组，名称是${cluster_id}-common, 如果没有配置默认安全组，在这里过滤一下默认安全组
		defaultSecurityGroupName := fmt.Sprintf("%v-common", result["ClusterId"])
		nameMap := make(map[string]string)
		filteredSecurityGroupIds := make([]interface{}, 0)
		defaultCount := 0
		defaultSecurityGroupId := ""
		for _, securityGroup := range securityGroups {
			nameMap[securityGroup.(map[string]interface{})["SecurityGroupId"].(string)] = securityGroup.(map[string]interface{})["SecurityGroupName"].(string)
		}
		for _, securityGroupId := range tmpSecurityGroupIds {
			if nameMap[securityGroupId.(string)] == defaultSecurityGroupName {
				defaultCount++
				defaultSecurityGroupId = securityGroupId.(string)
				continue
			}
			filteredSecurityGroupIds = append(filteredSecurityGroupIds, securityGroupId)
		}
		if defaultCount > 1 {
			return nil, fmt.Errorf("default security group is not unique")
		}

		// 如果用户传了默认安全组id，不需要过滤
		oldSecurityGroupIds := resourceData.Get("node_config.0.security.0.security_group_ids").([]interface{})
		useDefaultSecurityGroupId := false
		for _, securityGroupId := range oldSecurityGroupIds {
			if securityGroupId.(string) == defaultSecurityGroupId {
				useDefaultSecurityGroupId = true
			}
		}
		if !useDefaultSecurityGroupId {
			result["NodeConfig"].(map[string]interface{})["Security"].(map[string]interface{})["SecurityGroupIds"] = filteredSecurityGroupIds
		}

		logger.Debug(logger.RespFormat, "filteredSecurityGroupIds", tmpSecurityGroupIds, filteredSecurityGroupIds)
	}

	if ecsTags, ok := result["NodeConfig"].(map[string]interface{})["Tags"]; ok {
		result["NodeConfig"].(map[string]interface{})["EcsTags"] = ecsTags
		delete(result["NodeConfig"].(map[string]interface{}), "Tags")
	}

	logger.Debug(logger.RespFormat, "result of ReadResource ", result)
	return result, err
}

func (s *VestackNodePoolService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			failStates = append(failStates, "Failed")
			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			//这里vke是Status是一个Object，取Phase字段判断是否失败
			status = demo["Status"].(map[string]interface{})["Phase"]
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("node pool status error, status:%s", status.(string))
				}
			}
			return demo, status.(string), err
		},
	}

}

func (VestackNodePoolService) WithResourceResponseHandlers(nodePool map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		var (
			security     = make([]interface{}, 0)
			systemVolume = make([]interface{}, 0)
			login        = make([]interface{}, 0)
		)

		priSecurity := nodePool["NodeConfig"].(map[string]interface{})["Security"]
		priLogin := priSecurity.(map[string]interface{})["Login"]
		delete(nodePool, "Login")
		login = append(login, priLogin)
		priSecurity.(map[string]interface{})["Login"] = login
		security = append(security, priSecurity)

		delete(nodePool, "Security")
		nodePool["NodeConfig"].(map[string]interface{})["Security"] = security

		priSystemVolume := nodePool["NodeConfig"].(map[string]interface{})["SystemVolume"]
		systemVolume = append(systemVolume, priSystemVolume)
		delete(nodePool, "SystemVolume")
		nodePool["NodeConfig"].(map[string]interface{})["SystemVolume"] = systemVolume

		return nodePool, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackNodePoolService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateNodePool",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"cluster_id": {
					TargetField: "ClusterId",
				},
				"client_token": {
					TargetField: "ClientToken",
				},
				"name": {
					TargetField: "Name",
				},
				"node_config": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"instance_type_ids": {
							ConvertType: bp.ConvertJsonArray,
						},
						"subnet_ids": {
							ConvertType: bp.ConvertJsonArray,
						},
						"security": {
							ConvertType: bp.ConvertJsonObject,
							NextLevelConvert: map[string]bp.RequestConvert{
								"login": {
									ConvertType: bp.ConvertJsonObject,
									NextLevelConvert: map[string]bp.RequestConvert{
										"password": {
											ConvertType: bp.ConvertJsonObject,
										},
										"ssh_key_pair_name": {
											ConvertType: bp.ConvertJsonObject,
										},
									},
								},
								"security_group_ids": {
									ConvertType: bp.ConvertJsonArray,
								},
								"security_strategies": {
									ConvertType: bp.ConvertJsonArray,
								},
							},
						},
						"system_volume": {
							ConvertType: bp.ConvertJsonObject,
							NextLevelConvert: map[string]bp.RequestConvert{
								"type": {
									ConvertType: bp.ConvertJsonObject,
								},
								"size": {
									ConvertType: bp.ConvertJsonObject,
								},
							},
						},
						"data_volumes": {
							ConvertType: bp.ConvertJsonObjectArray,
						},
						"initialize_script": {
							ConvertType: bp.ConvertJsonObject,
						},
						"additional_container_storage_enabled": {
							ConvertType: bp.ConvertJsonObject,
						},
						"image_id": {
							ConvertType: bp.ConvertJsonObject,
						},
						"instance_charge_type": {
							ConvertType: bp.ConvertJsonObject,
						},
						"period": {
							ConvertType: bp.ConvertJsonObject,
						},
						"auto_renew": {
							ForceGet:    true,
							TargetField: "AutoRenew",
						},
						"auto_renew_period": {
							ConvertType: bp.ConvertJsonObject,
						},
						"name_prefix": {
							ConvertType: bp.ConvertJsonObject,
						},
						"ecs_tags": {
							TargetField: "Tags",
							ConvertType: bp.ConvertJsonObjectArray,
						},
						"hpc_cluster_ids": {
							ConvertType: bp.ConvertJsonArray,
						},
					},
				},
				"kubernetes_config": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"labels": {
							ConvertType: bp.ConvertJsonObjectArray,
						},
						"taints": {
							ConvertType: bp.ConvertJsonObjectArray,
						},
						"cordon": {
							ConvertType: bp.ConvertJsonObject,
						},
					},
				},
				"auto_scaling": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"enabled": {
							TargetField: "Enabled",
						},
						"max_replicas": {
							TargetField: "MaxReplicas",
						},
						"min_replicas": {
							TargetField: "MinReplicas",
						},
						"desired_replicas": {
							TargetField: "DesiredReplicas",
						},
						"priority": {
							TargetField: "Priority",
						},
						"subnet_policy": {
							TargetField: "SubnetPolicy",
						},
					},
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertJsonObjectArray,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if chargeType, ok := (*call.SdkParam)["NodeConfig.InstanceChargeType"]; ok {
					if autoScalingEnabled, ok := (*call.SdkParam)["AutoScaling.Enabled"]; ok {
						if chargeType.(string) == "PrePaid" && autoScalingEnabled.(bool) {
							return false, fmt.Errorf("PrePaid charge type cannot support auto scaling")
						}
					}
				}
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
			Refresh: &bp.StateRefresh{
				Target:  []string{"Running"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackNodePoolService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateNodePoolConfig",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"cluster_id": {
					TargetField: "ClusterId",
				},
				"client_token": {
					TargetField: "ClientToken",
				},
				"name": {
					TargetField: "Name",
				},
				"node_config": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"security": {
							ConvertType: bp.ConvertJsonObject,
							NextLevelConvert: map[string]bp.RequestConvert{
								"login": {
									ConvertType: bp.ConvertJsonObject,
									NextLevelConvert: map[string]bp.RequestConvert{
										"password": {
											ConvertType: bp.ConvertJsonObject,
										},
										"ssh_key_pair_name": {
											ConvertType: bp.ConvertJsonObject,
										},
									},
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
						"subnet_ids": {
							ConvertType: bp.ConvertJsonArray,
						},
						"period": {
							ConvertType: bp.ConvertJsonObject,
						},
						"auto_renew": {
							ForceGet:    true,
							TargetField: "AutoRenew",
						},
						"auto_renew_period": {
							ConvertType: bp.ConvertJsonObject,
						},
						"name_prefix": {
							ConvertType: bp.ConvertJsonObject,
						},
						"ecs_tags": {
							TargetField: "Tags",
							ConvertType: bp.ConvertJsonObjectArray,
						},
						"instance_type_ids": {
							ConvertType: bp.ConvertJsonArray,
						},
						"hpc_cluster_ids": {
							ConvertType: bp.ConvertJsonArray,
						},
						"image_id": {
							ConvertType: bp.ConvertJsonObject,
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
				"auto_scaling": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"enabled": {
							ForceGet:    true,
							TargetField: "Enabled",
						},
						"max_replicas": {
							ForceGet:    true,
							TargetField: "MaxReplicas",
						},
						"min_replicas": {
							ForceGet:    true,
							TargetField: "MinReplicas",
						},
						"desired_replicas": {
							ForceGet:    true,
							TargetField: "DesiredReplicas",
						},
						"priority": {
							ForceGet:    true,
							TargetField: "Priority",
						},
						"subnet_policy": {
							ForceGet:    true,
							TargetField: "SubnetPolicy",
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
						_, exist := security.(map[string]interface{})["SecurityStrategies"]
						if !exist && d.HasChange("node_config.0.security.0.security_strategies") {
							security.(map[string]interface{})["SecurityStrategies"] = []interface{}{}
						}
					}

					if _, ok1 := nodeconfig.(map[string]interface{})["HpcClusterIds"]; ok1 {
						if _, ok2 := nodeconfig.(map[string]interface{})["InstanceTypeIds"]; !ok2 {
							(*call.SdkParam)["NodeConfig"].(map[string]interface{})["InstanceTypeIds"] = make([]interface{}, 0)
							instanceTypeIds := d.Get("node_config.0.instance_type_ids")
							for _, instanceTypeId := range instanceTypeIds.([]interface{}) {
								(*call.SdkParam)["NodeConfig"].(map[string]interface{})["InstanceTypeIds"] = append((*call.SdkParam)["NodeConfig"].(map[string]interface{})["InstanceTypeIds"].([]interface{}), instanceTypeId.(string))
							}
						}
					}
					if d.HasChange("node_config.0.hpc_cluster_ids") {
						bp.DefaultMapValue(call.SdkParam, "NodeConfig", map[string]interface{}{
							"HpcClusterIds": []interface{}{},
						})
					}
				}

				instanceChargeType := d.Get("node_config").([]interface{})[0].(map[string]interface{})["instance_charge_type"].(string)
				if instanceChargeType != "PrePaid" {
					if nodeCfg, ok := (*call.SdkParam)["NodeConfig"]; ok {
						if _, ok := nodeCfg.(map[string]interface{})["AutoRenew"]; ok {
							delete((*call.SdkParam)["NodeConfig"].(map[string]interface{}), "AutoRenew")
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

				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Running"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	callbacks = append(callbacks, callback)

	// 更新Tags
	callbacks = s.setResourceTags(resourceData, "NodePool", callbacks)

	return callbacks
}

func (s *VestackNodePoolService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteNodePool",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"Id":                       resourceData.Id(),
				"ClusterId":                resourceData.Get("cluster_id"),
				"CascadingDeleteResources": [1]string{"Ecs"},
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

func (s *VestackNodePoolService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "Filter.Ids",
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
			"cluster_id": {
				TargetField: "Filter.ClusterId",
			},
			"cluster_ids": {
				TargetField: "Filter.ClusterIds",
				ConvertType: bp.ConvertJsonArray,
			},
			"name": {
				TargetField: "Filter.Name",
			},
			"create_client_token": {
				TargetField: "Filter.CreateClientToken",
			},
			"update_client_token": {
				TargetField: "Filter.UpdateClientToken",
			},
			"tags": {
				TargetField: "Tags",
				ConvertType: bp.ConvertJsonObjectArray,
			},
		},
		NameField:    "Name",
		IdField:      "Id",
		CollectField: "node_pools",
		ContentType:  bp.ContentTypeJson,
		ResponseConverts: map[string]bp.ResponseConvert{
			"Status.Phase": {
				TargetField: "phase",
			},
			"Status.Conditions": {
				TargetField: "condition_types",
				Convert: func(i interface{}) interface{} {
					var results []interface{}
					if dd, ok := i.([]interface{}); ok {
						for _, _data := range dd {
							results = append(results, _data.(map[string]interface{})["Type"])
						}
					}
					return results
				},
			},
			"AutoScaling.Enabled": {
				TargetField: "enabled",
			},
			"AutoScaling.DesiredReplicas": {
				TargetField: "desired_replicas",
			},
			"AutoScaling.MinReplicas": {
				TargetField: "min_replicas",
			},
			"AutoScaling.MaxReplicas": {
				TargetField: "max_replicas",
			},
			"AutoScaling.Priority": {
				TargetField: "priority",
			},
			"AutoScaling.SubnetPolicy": {
				TargetField: "subnet_policy",
			},
			"KubernetesConfig.Cordon": {
				TargetField: "cordon",
			},
			"KubernetesConfig.Labels": {
				TargetField: "label_content",
				Convert: func(i interface{}) interface{} {
					var results []interface{}
					if dd, ok := i.([]interface{}); ok {
						for _, _data := range dd {
							label := make(map[string]string, 0)
							label["key"] = _data.(map[string]interface{})["Key"].(string)
							label["value"] = _data.(map[string]interface{})["Value"].(string)
							results = append(results, label)
						}
					}
					return results
				},
			},
			"KubernetesConfig.Taints": {
				TargetField: "taint_content",
				Convert: func(i interface{}) interface{} {
					var results []interface{}
					if dd, ok := i.([]interface{}); ok {
						for _, _data := range dd {
							label := make(map[string]string, 0)
							label["key"] = _data.(map[string]interface{})["Key"].(string)
							label["value"] = _data.(map[string]interface{})["Value"].(string)
							label["effect"] = _data.(map[string]interface{})["Effect"].(string)
							results = append(results, label)
						}
					}
					return results
				},
			},
			"NodeConfig.InitializeScript": {
				TargetField: "initialize_script",
			},
			"NodeConfig.AdditionalContainerStorageEnabled": {
				TargetField: "additional_container_storage_enabled",
			},
			"NodeConfig.InstanceTypeIds": {
				TargetField: "instance_type_ids",
				Convert: func(i interface{}) interface{} {
					var results []interface{}
					if dd, ok := i.([]interface{}); ok {
						results = dd
					}
					return results
				},
			},
			"NodeConfig.SubnetIds": {
				TargetField: "subnet_ids",
				Convert: func(i interface{}) interface{} {
					var results []interface{}
					if dd, ok := i.([]interface{}); ok {
						results = dd
					}
					return results
				},
			},
			"NodeConfig.ImageId": {
				TargetField: "image_id",
			},
			"NodeConfig.SystemVolume": {
				TargetField: "system_volume",
				Convert: func(i interface{}) interface{} {
					var results []interface{}
					if i.(map[string]interface{})["Type"] == nil || i.(map[string]interface{})["Size"] == nil {
						return results
					}
					volume := make(map[string]interface{}, 0)
					volume["type"] = i.(map[string]interface{})["Type"].(string)
					volume["size"] = strconv.FormatFloat(i.(map[string]interface{})["Size"].(float64), 'g', 5, 32)
					results = append(results, volume)
					return results
				},
			},
			"NodeConfig.DataVolumes": {
				TargetField: "data_volumes",
				Convert: func(i interface{}) interface{} {
					var results []interface{}
					if dd, ok := i.([]interface{}); ok {
						for _, _data := range dd {
							volume := make(map[string]interface{}, 0)
							volume["size"] = strconv.FormatFloat(_data.(map[string]interface{})["Size"].(float64), 'g', 5, 32)
							volume["type"] = _data.(map[string]interface{})["Type"].(string)
							if p, ok := _data.(map[string]interface{})["MountPoint"]; ok { // 可能不存在
								volume["mount_point"] = p.(string)
							}
							results = append(results, volume)
						}
					}
					return results
				},
			},
			"NodeConfig.Security.SecurityGroupIds": {
				TargetField: "security_group_ids",
				Convert: func(i interface{}) interface{} {
					var results []interface{}
					if dd, ok := i.([]interface{}); ok {
						results = dd
					}
					return results
				},
			},
			"NodeConfig.Security.SecurityStrategyEnabled": {
				TargetField: "security_strategy_enabled",
			},
			"NodeConfig.Security.SecurityStrategies": {
				TargetField: "security_strategies",
				Convert: func(i interface{}) interface{} {
					var results []interface{}
					if dd, ok := i.([]interface{}); ok {
						results = dd
					}
					return results
				},
			},
			"NodeConfig.Security.Login.Type": {
				TargetField: "login_type",
			},
			"NodeConfig.Security.Login.SshKeyPairName": {
				TargetField: "login_key_pair_name",
			},
			"NodeConfig.InstanceChargeType": {
				TargetField: "instance_charge_type",
			},
			"NodeConfig.Period": {
				TargetField: "period",
			},
			"NodeConfig.AutoRenew": {
				TargetField: "auto_renew",
			},
			"NodeConfig.AutoRenewPeriod": {
				TargetField: "auto_renew_period",
			},
			"NodeConfig.NamePrefix": {
				TargetField: "name_prefix",
			},
			"NodeConfig.HpcClusterIds": {
				TargetField: "hpc_cluster_ids",
			},
			"NodeConfig.Tags": {
				TargetField: "ecs_tags",
				Convert: func(i interface{}) interface{} {
					var results []interface{}
					if dd, ok := i.([]interface{}); ok {
						for _, data := range dd {
							tag := make(map[string]interface{}, 0)
							tag["key"] = data.(map[string]interface{})["Key"].(string)
							tag["value"] = data.(map[string]interface{})["Value"].(string)
							results = append(results, tag)
						}
					}
					return results
				},
			},
			"NodeStatistics": {
				TargetField: "node_statistics",
				Convert: func(i interface{}) interface{} {
					label := make(map[string]interface{}, 0)
					label["total_count"] = int(i.(map[string]interface{})["TotalCount"].(float64))
					label["creating_count"] = int(i.(map[string]interface{})["CreatingCount"].(float64))
					label["running_count"] = int(i.(map[string]interface{})["RunningCount"].(float64))
					label["updating_count"] = int(i.(map[string]interface{})["UpdatingCount"].(float64))
					label["deleting_count"] = int(i.(map[string]interface{})["DeletingCount"].(float64))
					label["failed_count"] = int(i.(map[string]interface{})["FailedCount"].(float64))
					label["stopped_count"] = int(i.(map[string]interface{})["StoppedCount"].(float64))
					label["stopping_count"] = int(i.(map[string]interface{})["StoppingCount"].(float64))
					label["starting_count"] = int(i.(map[string]interface{})["StartingCount"].(float64))
					return label
				},
			},
		},
	}
}

func (s *VestackNodePoolService) ReadResourceId(id string) string {
	return id
}

func (s *VestackNodePoolService) setResourceTags(resourceData *schema.ResourceData, resourceType string, callbacks []bp.Callback) []bp.Callback {
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

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vke",
		Version:     "2022-05-12",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
