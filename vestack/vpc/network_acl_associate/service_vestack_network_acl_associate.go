package network_acl_associate

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/network_acl"
)

type VestackNetworkAclAssociateService struct {
	Client *bp.SdkClient
}

func NewNetworkAclAssociateService(c *bp.SdkClient) *VestackNetworkAclAssociateService {
	return &VestackNetworkAclAssociateService{
		Client: c,
	}
}

func (s *VestackNetworkAclAssociateService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackNetworkAclAssociateService) ReadResources(condition map[string]interface{}) ([]interface{}, error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
		err     error
	)

	return bp.WithSimpleQuery(condition, func(m map[string]interface{}) ([]interface{}, error) {
		action := "DescribeNetworkAcls"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return []interface{}{}, err
			}
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return []interface{}{}, err
			}
		}
		logger.Debug(logger.RespFormat, action, condition, *resp)

		results, err = bp.ObtainSdkValue("Result.NetworkAcls", *resp)
		if err != nil {
			return []interface{}{}, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if _, ok = results.([]interface{}); !ok {
			return []interface{}{}, errors.New("Result.NetworkAcls is not Slice")
		}
		return results.([]interface{}), err
	})
}

func (s *VestackNetworkAclAssociateService) ReadResource(resourceData *schema.ResourceData, associateId string) (data map[string]interface{}, err error) {
	if associateId == "" {
		associateId = resourceData.Id()
	}

	ids := strings.Split(associateId, ":")
	if len(ids) != 2 {
		return map[string]interface{}{}, fmt.Errorf("invalid acl associateId: %s", associateId)
	}

	networkAclId := ids[0]
	resourceId := ids[1]
	req := map[string]interface{}{
		"NetworkAclIds.1": networkAclId,
	}

	networkAcls, err := s.ReadResources(req)
	if err != nil {
		return nil, err
	}
	if len(networkAcls) == 0 {
		return map[string]interface{}{}, fmt.Errorf("network acl %s not exist ", networkAclId)
	}
	for _, v := range networkAcls {
		if _, ok := v.(map[string]interface{}); !ok {
			return map[string]interface{}{}, errors.New("Value is not map ")
		}
	}

	aclResources := networkAcls[0].(map[string]interface{})["Resources"]
	if len(aclResources.([]interface{})) == 0 {
		return map[string]interface{}{}, fmt.Errorf("network acl resource %s:%s not exist ", networkAclId, resourceId)
	}
	for _, v := range aclResources.([]interface{}) {
		if _, ok := v.(map[string]interface{}); !ok {
			return map[string]interface{}{}, errors.New("Value is not map ")
		}
		if v.(map[string]interface{})["ResourceId"] == resourceId {
			data = v.(map[string]interface{})
		}
	}

	return data, err
}

func (s *VestackNetworkAclAssociateService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (VestackNetworkAclAssociateService) WithResourceResponseHandlers(aclEntry map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return aclEntry, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackNetworkAclAssociateService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AssociateNetworkAcl",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["NetworkAclId"] = d.Get("network_acl_id")
				(*call.SdkParam)["Resource.1.ResourceId"] = d.Get("resource_id")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				// ResourceData中，network_acl_associate的Id形式为'network_acl_id:resource_id'
				id := fmt.Sprintf("%s:%s", d.Get("network_acl_id"), d.Get("resource_id"))
				d.SetId(id)
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("network_acl_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				network_acl.NewNetworkAclService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("network_acl_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackNetworkAclAssociateService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackNetworkAclAssociateService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DisassociateNetworkAcl",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				aclAssociateId := d.Id()
				ids := strings.Split(aclAssociateId, ":")
				if len(ids) != 2 {
					return false, fmt.Errorf("error network acl associate id: %s", aclAssociateId)
				}
				(*call.SdkParam)["NetworkAclId"] = ids[0]
				(*call.SdkParam)["Resource.1.ResourceId"] = ids[1]
				return true, nil
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
							return resource.NonRetryableError(fmt.Errorf("error on  reading acl entry on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("network_acl_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				network_acl.NewNetworkAclService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutDelete),
					ResourceId: resourceData.Get("network_acl_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackNetworkAclAssociateService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackNetworkAclAssociateService) ReadResourceId(id string) string {
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
