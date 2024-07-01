package network_acl

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/vpc"
)

type VestackNetworkAclService struct {
	Client *bp.SdkClient
}

func NewNetworkAclService(c *bp.SdkClient) *VestackNetworkAclService {
	return &VestackNetworkAclService{
		Client: c,
	}
}

func (s *VestackNetworkAclService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackNetworkAclService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		action := "DescribeNetworkAcls"
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

		results, err = bp.ObtainSdkValue("Result.NetworkAcls", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.NetworkAcls is not Slice")
		}

		return data, err
	})
}

func (s *VestackNetworkAclService) ReadResource(resourceData *schema.ResourceData, networkAclId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if networkAclId == "" {
		networkAclId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"NetworkAclIds.1": networkAclId,
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
		return data, fmt.Errorf("network acl %s is not exist ", networkAclId)
	}

	// 删除默认创建的拒绝规则
	if ingressAclEntries, ok := data["IngressAclEntries"]; ok {
		var tempEntries []interface{}
		for _, entry := range ingressAclEntries.([]interface{}) {
			if priority, ok := entry.(map[string]interface{})["Priority"]; ok && priority.(float64) < 100 {
				tempEntries = append(tempEntries, entry)
			}
		}
		data["IngressAclEntries"] = tempEntries
	}
	if egressAclEntries, ok := data["EgressAclEntries"]; ok {
		var tempEntries []interface{}
		for _, entry := range egressAclEntries.([]interface{}) {
			if priority, ok := entry.(map[string]interface{})["Priority"]; ok && priority.(float64) < 100 {
				tempEntries = append(tempEntries, entry)
			}
		}
		data["EgressAclEntries"] = tempEntries
	}

	return data, err
}

func (s *VestackNetworkAclService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				demo   map[string]interface{}
				status interface{}
			)
			//no failed status.
			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", demo)
			if err != nil {
				return nil, "", err
			}
			return demo, status.(string), err
		},
	}
}

func (VestackNetworkAclService) WithResourceResponseHandlers(acl map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return acl, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackNetworkAclService) CreateResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateNetworkAcl",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"ingress_acl_entries": {
					Ignore: true,
				},
				"egress_acl_entries": {
					Ignore: true,
				},
				"resources": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.NetworkAclId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				vpc.NewVpcService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("vpc_id").(string),
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
		},
	}
	callbacks = append(callbacks, callback)

	// 规则创建
	entryCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateNetworkAclEntries",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"ingress_acl_entries": {
					ConvertType: bp.ConvertListN,
				},
				"egress_acl_entries": {
					ConvertType: bp.ConvertListN,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				if len(*call.SdkParam) > 0 {
					(*call.SdkParam)["NetworkAclId"] = d.Id()
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					if _, ok := d.GetOk("ingress_acl_entries"); ok {
						(*call.SdkParam)["UpdateIngressAclEntries"] = true
					}
					if _, ok := d.GetOk("egress_acl_entries"); ok {
						(*call.SdkParam)["UpdateEgressAclEntries"] = true
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, entryCallback)

	return callbacks
}

func (s *VestackNetworkAclService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyNetworkAclAttributes",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"ingress_acl_entries": {
					Ignore: true,
				},
				"egress_acl_entries": {
					Ignore: true,
				},
				"resources": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["NetworkAclId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				vpc.NewVpcService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("vpc_id").(string),
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
		},
	}
	callbacks = append(callbacks, callback)

	// 规则修改
	if resourceData.HasChange("ingress_acl_entries") {
		ingressUpdateCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "UpdateNetworkAclEntries",
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"ingress_acl_entries": {
						ConvertType: bp.ConvertListN,
						ForceGet:    true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) > 0 {
						(*call.SdkParam)["NetworkAclId"] = d.Id()
						(*call.SdkParam)["ClientToken"] = uuid.New().String()
						(*call.SdkParam)["UpdateIngressAclEntries"] = true
						for index, entry := range d.Get("ingress_acl_entries").([]interface{}) {
							(*call.SdkParam)["IngressAclEntries."+strconv.Itoa(index+1)+".NetworkAclEntryId"] = entry.(map[string]interface{})["network_acl_entry_id"].(string)
						}
						return true, nil
					}
					return false, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Available"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
				LockId: func(d *schema.ResourceData) string {
					return d.Get("vpc_id").(string)
				},
				ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
					vpc.NewVpcService(s.Client): {
						Target:     []string{"Available"},
						Timeout:    resourceData.Timeout(schema.TimeoutCreate),
						ResourceId: resourceData.Get("vpc_id").(string),
					},
				},
			},
		}
		callbacks = append(callbacks, ingressUpdateCallback)
	}
	if resourceData.HasChange("egress_acl_entries") {
		ingressUpdateCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "UpdateNetworkAclEntries",
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"egress_acl_entries": {
						ConvertType: bp.ConvertListN,
						ForceGet:    true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) > 0 {
						(*call.SdkParam)["NetworkAclId"] = d.Id()
						(*call.SdkParam)["ClientToken"] = uuid.New().String()
						(*call.SdkParam)["UpdateEgressAclEntries"] = true
						for index, entry := range d.Get("egress_acl_entries").([]interface{}) {
							(*call.SdkParam)["EgressAclEntries."+strconv.Itoa(index+1)+".NetworkAclEntryId"] = entry.(map[string]interface{})["network_acl_entry_id"].(string)
						}
						return true, nil
					}
					return false, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Available"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
			},
		}
		callbacks = append(callbacks, ingressUpdateCallback)
	}

	return callbacks
}

func (s *VestackNetworkAclService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	removeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteNetworkAcl",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"NetworkAclId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
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
							return resource.NonRetryableError(fmt.Errorf("error on  reading network acl on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				vpc.NewVpcService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("vpc_id").(string),
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
		},
	}
	callbacks = append(callbacks, removeCallback)

	return callbacks
}

func (s *VestackNetworkAclService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "NetworkAclIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "NetworkAclName",
		IdField:      "NetworkAclId",
		CollectField: "network_acls",
		ResponseConverts: map[string]bp.ResponseConvert{
			"NetworkAclId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *VestackNetworkAclService) ReadResourceId(id string) string {
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

func (s *VestackNetworkAclService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "vpc",
		ResourceType:         "networkacl",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}
