package security_group

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/vpc"
)

type VestackSecurityGroupService struct {
	Client *bp.SdkClient
}

func NewSecurityGroupService(c *bp.SdkClient) *VestackSecurityGroupService {
	return &VestackSecurityGroupService{
		Client: c,
	}
}

func (s *VestackSecurityGroupService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackSecurityGroupService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		vpcClient := s.Client.VpcClient
		action := "DescribeSecurityGroups"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = vpcClient.DescribeSecurityGroupsCommon(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = vpcClient.DescribeSecurityGroupsCommon(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = bp.ObtainSdkValue("Result.SecurityGroups", *resp)

		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.SecurityGroups is not Slice")
		}
		logger.Debug(logger.ReqFormat, "", data)

		return data, err
	})
}

func (s *VestackSecurityGroupService) ReadResource(resourceData *schema.ResourceData, securityGroupId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if securityGroupId == "" {
		securityGroupId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"SecurityGroupIds.1": securityGroupId,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("value is not map")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("SecurityGroup %s not exist ", securityGroupId)
	}
	return data, err
}

func (s *VestackSecurityGroupService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
					return nil, "", fmt.Errorf("subnet status error, status:%s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}

}

func (VestackSecurityGroupService) WithResourceResponseHandlers(subnet map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return subnet, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackSecurityGroupService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateSecurityGroup",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertListN,
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.VpcClient.CreateSecurityGroupCommon(call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.SecurityGroupId", *resp)
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
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackSecurityGroupService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifySecurityGroupAttributes",
			ConvertMode: bp.RequestConvertAll,
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["SecurityGroupId"] = d.Id()
				delete(*call.SdkParam, "Tags")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.VpcClient.ModifySecurityGroupAttributesCommon(call.SdkParam)
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
		},
	}
	callbacks = append(callbacks, callback)

	// 更新Tags
	setResourceTagsCallbacks := bp.SetResourceTags(s.Client, "TagResources", "UntagResources", "securitygroup", resourceData, getUniversalInfo)
	callbacks = append(callbacks, setResourceTagsCallbacks...)

	return callbacks
}

func (s *VestackSecurityGroupService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteSecurityGroup",
			ConvertMode: bp.RequestConvertIgnore,
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
			SdkParam: &map[string]interface{}{
				"SecurityGroupId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.VpcClient.DeleteSecurityGroupCommon(call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading security group on delete %q, %w", d.Id(), callErr))
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

func (s *VestackSecurityGroupService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "SecurityGroupIds",
				ConvertType: bp.ConvertWithN,
			},
			"security_group_names": {
				TargetField: "SecurityGroupNames",
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
		NameField:    "SecurityGroupName",
		IdField:      "SecurityGroupId",
		CollectField: "security_groups",
		ResponseConverts: map[string]bp.ResponseConvert{
			"SecurityGroupId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *VestackSecurityGroupService) ReadResourceId(id string) string {
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

func (s *VestackSecurityGroupService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "vpc",
		ResourceType:         "securitygroup",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}
