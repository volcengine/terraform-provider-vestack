package iam_service_linked_role

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

type VestackIamServiceLinkedRoleService struct {
	Client *bp.SdkClient
}

func (v *VestackIamServiceLinkedRoleService) GetClient() *bp.SdkClient {
	return v.Client
}

func (v *VestackIamServiceLinkedRoleService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageOffsetQuery(m, "Limit", "Offset", 100, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListRoles"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = v.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = v.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
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

func (v *VestackIamServiceLinkedRoleService) ReadResource(resourceData *schema.ResourceData, roleId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if roleId == "" {
		roleId = v.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(roleId, ":")
	if len(ids) != 2 {
		return nil, errors.New("error id")
	}
	req := map[string]interface{}{
		"RoleName": ids[1],
	}
	results, err = v.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, r := range results {
		if data, ok = r.(map[string]interface{}); !ok {
			return data, errors.New("value is not map")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Role %s not exist ", roleId)
	}
	return data, err
}

func (v *VestackIamServiceLinkedRoleService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, s string) *resource.StateChangeConf {
	return nil
}

func (v *VestackIamServiceLinkedRoleService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
	return []bp.ResourceResponseHandler{}
}

func (v *VestackIamServiceLinkedRoleService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	createIamRoleCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateServiceLinkedRole",
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return v.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				serviceName := d.Get("service_name").(string)
				currentRoleName := ""
				results, err := v.ReadResources(map[string]interface{}{})
				if err != nil {
					return err
				}
				for _, r := range results {
					if result, ok := r.(map[string]interface{}); ok {
						roleName := result["RoleName"].(string)
						if strings.Contains(roleName, "ServiceRoleFor") {
							roleService := strings.ToLower(roleName[14:])
							if roleService == strings.ToLower(bp.DownLineToHump(serviceName)) {
								currentRoleName = roleName
							}
						}
					}
				}
				roleName := serviceName + ":" + currentRoleName
				d.SetId(roleName)
				return nil
			},
		},
	}
	return []bp.Callback{createIamRoleCallback}
}

func (v *VestackIamServiceLinkedRoleService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}

func (v *VestackIamServiceLinkedRoleService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	ids := strings.Split(data.Id(), ":")
	deleteRoleCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteRole",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["RoleName"] = ids[1]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return v.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := v.ReadResource(d, "")
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

func (v *VestackIamServiceLinkedRoleService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (v *VestackIamServiceLinkedRoleService) ReadResourceId(s string) string {
	return s
}

func NewIamServiceLinkedRoleService(c *bp.SdkClient) *VestackIamServiceLinkedRoleService {
	return &VestackIamServiceLinkedRoleService{
		Client: c,
	}
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "iam",
		Version:     "2018-01-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}
