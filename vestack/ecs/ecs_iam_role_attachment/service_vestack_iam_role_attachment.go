package ecs_iam_role_attachment

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamRoleAttachmentService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewIamRoleAttachmentService(c *bp.SdkClient) *VestackIamRoleAttachmentService {
	return &VestackIamRoleAttachmentService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackIamRoleAttachmentService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamRoleAttachmentService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return nil, nil
}

func (s *VestackIamRoleAttachmentService) ReadResource(resourceData *schema.ResourceData, id string) (map[string]interface{}, error) {
	var (
		results interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	ids := strings.Split(id, ":")
	iamRoleName := ids[0]
	instanceId := ids[1]

	data := make(map[string]interface{})
	action := "DescribeInstancesIamRoles"
	req := map[string]interface{}{
		"InstanceIds.1": instanceId,
	}
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, req, resp)

	results, err = bp.ObtainSdkValue("Result.InstancesIamRoles", *resp)
	if err != nil {
		return data, err
	}
	iamRoles, ok := results.([]interface{})
	if !ok {
		return data, fmt.Errorf("Result.InstancesIamRoles is not slice ")
	}
	if len(iamRoles) == 0 {
		return data, fmt.Errorf("iam_role_attachment %s not exist ", id)
	}
	iamRoleMap, ok := iamRoles[0].(map[string]interface{})
	if !ok {
		return data, fmt.Errorf("Value is not map ")
	}
	roleNames, ok := iamRoleMap["RoleNames"].([]interface{})
	if !ok {
		return data, fmt.Errorf("RoleNames is not slice ")
	}
	for _, v := range roleNames {
		if v.(string) == iamRoleName {
			data["IamRoleName"] = v.(string)
			data["InstanceId"] = iamRoleMap["InstanceId"]
		}
	}

	if len(data) == 0 {
		return data, fmt.Errorf("iam role %v does not attach to the instance %v", iamRoleName, instanceId)
	}
	return data, err
}

func (s *VestackIamRoleAttachmentService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (VestackIamRoleAttachmentService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackIamRoleAttachmentService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AssociateInstancesIamRole",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"iam_role_name": {
					TargetField: "IamRoleName",
				},
				"instance_id": {
					TargetField: "InstanceIds.1",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				if err != nil {
					return resp, err
				}

				operationDetails, err := bp.ObtainSdkValue("Result.OperationDetails", *resp)
				if err != nil {
					return resp, err
				}
				if _, ok := operationDetails.([]interface{}); !ok || len(operationDetails.([]interface{})) == 0 {
					return resp, fmt.Errorf("Result.OperationDetails is not slice")
				}
				operationDetail, ok := operationDetails.([]interface{})[0].(map[string]interface{})
				if !ok {
					return resp, fmt.Errorf("operation detail is not a map")
				}
				errStruct, exist := operationDetail["Error"]
				if !exist || errStruct == nil {
					return resp, nil
				}
				errMap, ok := errStruct.(map[string]interface{})
				if !ok {
					return resp, fmt.Errorf("operation detail Error is not a map")
				}
				code := errMap["Code"].(string)
				message := errMap["Message"].(string)
				return resp, fmt.Errorf(code + ": " + message)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId((*call.SdkParam)["IamRoleName"].(string) + ":" + (*call.SdkParam)["InstanceIds.1"].(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackIamRoleAttachmentService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackIamRoleAttachmentService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	ids := strings.Split(s.ReadResourceId(resourceData.Id()), ":")
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DisassociateInstancesIamRole",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"IamRoleName":   ids[0],
				"InstanceIds.1": ids[1],
				"ClientToken":   uuid.New().String(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				if err != nil {
					return resp, err
				}

				operationDetails, err := bp.ObtainSdkValue("Result.OperationDetails", *resp)
				if err != nil {
					return resp, err
				}
				if _, ok := operationDetails.([]interface{}); !ok || len(operationDetails.([]interface{})) == 0 {
					return resp, fmt.Errorf("Result.OperationDetails is not slice")
				}
				operationDetail, ok := operationDetails.([]interface{})[0].(map[string]interface{})
				if !ok {
					return resp, fmt.Errorf("operation detail is not a map")
				}
				errStruct, exist := operationDetail["Error"]
				if !exist || errStruct == nil {
					return resp, nil
				}
				errMap, ok := errStruct.(map[string]interface{})
				if !ok {
					return resp, fmt.Errorf("operation detail Error is not a map")
				}
				code := errMap["Code"].(string)
				message := errMap["Message"].(string)
				return resp, fmt.Errorf(code + ": " + message)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading iam role attachment on delete %q, %w", d.Id(), callErr))
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

func (s *VestackIamRoleAttachmentService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackIamRoleAttachmentService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "ecs",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
