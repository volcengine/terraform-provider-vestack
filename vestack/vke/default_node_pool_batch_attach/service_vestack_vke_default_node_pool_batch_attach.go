package default_node_pool_batch_attach

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/default_node_pool"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/node_pool"
)

type VestackVkeDefaultNodePoolBatchAttachService struct {
	Client                 *bp.SdkClient
	defaultNodePoolService *default_node_pool.VestackDefaultNodePoolService
}

func NewVestackVkeDefaultNodePoolBatchAttachService(c *bp.SdkClient) *VestackVkeDefaultNodePoolBatchAttachService {
	return &VestackVkeDefaultNodePoolBatchAttachService{
		Client:                 c,
		defaultNodePoolService: default_node_pool.NewDefaultNodePoolService(c),
	}
}

func (s *VestackVkeDefaultNodePoolBatchAttachService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackVkeDefaultNodePoolBatchAttachService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return data, err
}

func (s *VestackVkeDefaultNodePoolBatchAttachService) ReadResource(resourceData *schema.ResourceData, nodePoolId string) (data map[string]interface{}, err error) {
	if nodePoolId == "" {
		nodePoolId = s.ReadResourceId(resourceData.Id())
	}
	data, err = s.defaultNodePoolService.ReadResource(resourceData, nodePoolId)
	// 节点池和节点均有kubernetes config
	// 相同key以节点池的config为准，不同key时一起生效，因此需删除config保持不触发变更
	delete(data, "KubernetesConfig")
	logger.Debug(logger.ReqFormat, "VestackVkeDefaultNodePoolBatchAttachService ReadResource", "data", data)
	return data, err
}

func (s *VestackVkeDefaultNodePoolBatchAttachService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return s.defaultNodePoolService.RefreshResourceState(data, strings, duration, id)
}

func (s *VestackVkeDefaultNodePoolBatchAttachService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
	return s.defaultNodePoolService.WithResourceResponseHandlers(m)
}

func (s *VestackVkeDefaultNodePoolBatchAttachService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var calls []bp.Callback
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateDefaultNodePool",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				resp := make(map[string]interface{})
				resp["Id"] = (*call.SdkParam)["DefaultNodePoolId"]
				return &resp, nil
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Id", *resp)
				d.SetId(id.(string))
				return nil
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				node_pool.NewNodePoolService(s.Client): {
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
			},
		},
	}

	calls = append(calls, callback)

	calls = s.defaultNodePoolService.ProcessNodeInstances(resourceData, calls)

	return calls
}

func (s *VestackVkeDefaultNodePoolBatchAttachService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var calls []bp.Callback
	calls = s.defaultNodePoolService.ProcessNodeInstances(resourceData, calls)
	return calls
}

func (s *VestackVkeDefaultNodePoolBatchAttachService) RemoveResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var calls []bp.Callback
	var delNode []string
	nv := resourceData.Get("instances")
	if nv == nil {
		nv = new(schema.Set)
	}
	remove := nv.(*schema.Set)

	for _, v := range remove.List() {
		m := v.(map[string]interface{})
		delNode = append(delNode, m["id"].(string))
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

	return calls
}

func (s *VestackVkeDefaultNodePoolBatchAttachService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackVkeDefaultNodePoolBatchAttachService) ReadResourceId(id string) string {
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
