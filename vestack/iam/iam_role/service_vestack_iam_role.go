package iam_role

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamRoleService struct {
	Client *bp.SdkClient
}

func NewIamRoleService(c *bp.SdkClient) *VestackIamRoleService {
	return &VestackIamRoleService{
		Client: c,
	}
}

func (s *VestackIamRoleService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamRoleService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageOffsetQuery(m, "Limit", "Offset", 100, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListRoles"
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

		results, err = bp.ObtainSdkValue("Result.RoleMetadata", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.RoleMetadata is not Slice")
		}
		return data, err
	})
}

func (s *VestackIamRoleService) ReadResource(resourceData *schema.ResourceData, roleId string) (data map[string]interface{}, err error) {
	var (
		result interface{}
		ok     bool
	)
	if roleId == "" {
		roleId = s.ReadResourceId(resourceData.Id())
	}
	condition := map[string]interface{}{
		"RoleName": roleId,
	}
	action := "GetRole"
	logger.Debug(logger.ReqFormat, action, condition)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, condition, *resp)

	result, err = bp.ObtainSdkValue("Result.Role", *resp)
	if err != nil {
		return data, err
	}
	if data, ok = result.(map[string]interface{}); !ok {
		return data, errors.New("value is not map")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Role %s not exist ", roleId)
	}
	return data, err
}

func (s *VestackIamRoleService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *VestackIamRoleService) WithResourceResponseHandlers(role map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		role["Id"] = role["RoleName"]
		return role, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackIamRoleService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	createIamRoleCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateRole",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertListN,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				roleName, err := bp.ObtainSdkValue("Result.Role.RoleName", *resp)
				if err != nil {
					return err
				}
				if nameStr, ok := roleName.(string); ok && nameStr != "" {
					d.SetId(nameStr)
				}
				return nil
			},
		},
	}
	return []bp.Callback{createIamRoleCallback}
}

func (s *VestackIamRoleService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	updateRoleCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateRole",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"role_name": {
					TargetField: "NewRoleName",
					ConvertType: bp.ConvertDefault,
				},
				"display_name": {
					TargetField: "NewDisplayName",
					ConvertType: bp.ConvertDefault,
				},
				"description": {
					TargetField: "NewDescription",
					ConvertType: bp.ConvertDefault,
				},
				"max_session_duration": {
					TargetField: "MaxSessionDuration",
					ConvertType: bp.ConvertDefault,
				},
				"tags": {
					Ignore: true,
				},
			},
			RequestIdField: "RoleName",
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks := []bp.Callback{updateRoleCallback}
	setResourceTagsCallbacks := s.setResourceTags(data, "Role", callbacks)
	return setResourceTagsCallbacks
}

func (s *VestackIamRoleService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	deleteRoleCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteRole",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["RoleName"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
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
							return resource.NonRetryableError(fmt.Errorf("error on reading iam role on delete %q, %w", d.Id(), callErr))
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
	return []bp.Callback{deleteRoleCallback}
}

func (s *VestackIamRoleService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"query": {
				TargetField: "Query",
			},
		},
		ResponseConverts: map[string]bp.ResponseConvert{
			"RoleName": {
				TargetField: "role_name",
			},
			"RoleId": {
				TargetField: "role_id",
			},
			"IsServiceLinkedRole": {
				TargetField: "is_service_linked_role",
			},
			"DisplayName": {
				TargetField: "display_name",
			},
			"MaxSessionDuration": {
				TargetField: "max_session_duration",
			},
			"Tags": {
				TargetField: "tags",
			},
			"Trn": {
				TargetField: "trn",
			},
			"Description": {
				TargetField: "description",
			},
			"TrustPolicyDocument": {
				TargetField: "trust_policy_document",
			},
			"CreateDate": {
				TargetField: "create_date",
			},
			"UpdateDate": {
				TargetField: "update_date",
			},
		},
		NameField:    "RoleName",
		IdField:      "RoleName",
		CollectField: "roles",
	}
}

func (s *VestackIamRoleService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "iam",
		Action:      actionName,
		Version:     "2018-01-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		RegionType:  bp.Global,
	}
}

func (s *VestackIamRoleService) setResourceTags(resourceData *schema.ResourceData, resourceType string, callbacks []bp.Callback) []bp.Callback {
	addedTags, removedTags, _, _ := bp.GetSetDifference("tags", resourceData, bp.TagsHash, false)

	removeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UntagResources",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if removedTags != nil && len(removedTags.List()) > 0 {
					(*call.SdkParam)["ResourceNames.1"] = resourceData.Id()
					(*call.SdkParam)["ResourceType"] = resourceType
					for index, tag := range removedTags.List() {
						tm, ok := tag.(map[string]interface{})
						if !ok {
							return false, errors.New("tag item is not map")
						}
						key, ok := tm["key"].(string)
						if !ok {
							return false, errors.New("tag key is not string")
						}
						(*call.SdkParam)["TagKeys."+strconv.Itoa(index+1)] = key
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
					(*call.SdkParam)["ResourceNames.1"] = resourceData.Id()
					(*call.SdkParam)["ResourceType"] = resourceType
					for index, tag := range addedTags.List() {
						tm, ok := tag.(map[string]interface{})
						if !ok {
							return false, errors.New("tag item is not map")
						}
						key, ok := tm["key"].(string)
						if !ok {
							return false, errors.New("tag key is not string")
						}
						value, ok := tm["value"].(string)
						if !ok {
							return false, errors.New("tag value is not string")
						}
						(*call.SdkParam)["Tags."+strconv.Itoa(index+1)+".Key"] = key
						(*call.SdkParam)["Tags."+strconv.Itoa(index+1)+".Value"] = value
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
