package volume

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	re "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackVolumeService struct {
	Client *bp.SdkClient
}

func NewVolumeService(c *bp.SdkClient) *VestackVolumeService {
	return &VestackVolumeService{
		Client: c,
	}
}

func (s *VestackVolumeService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackVolumeService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		ebs := s.Client.EbsClient
		action := "DescribeVolumes"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = ebs.DescribeVolumesCommon(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = ebs.DescribeVolumesCommon(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = bp.ObtainSdkValue("Result.Volumes", *resp)
		if err != nil {
			return data, err
		}
		logger.Debug(logger.ReqFormat, action, results)
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Volumes is not Slice")
		}
		return data, err
	})
}

func (s *VestackVolumeService) ReadResource(resourceData *schema.ResourceData, volumeId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if volumeId == "" {
		volumeId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"VolumeIds.1": volumeId,
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
		return data, fmt.Errorf("volume %s not exist ", volumeId)
	}

	payType, ok := data["PayType"]
	if !ok {
		return data, fmt.Errorf(" PayType of volume is not exist ")
	}
	if payType.(string) == "post" {
		data["VolumeChargeType"] = "PostPaid"
	} else if payType.(string) == "pre" {
		data["VolumeChargeType"] = "PrePaid"
	}

	return data, err
}

func (s *VestackVolumeService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				ebs        map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "error")

			if err = resource.Retry(20*time.Minute, func() *resource.RetryError {
				ebs, err = s.ReadResource(resourceData, id)
				if err != nil {
					if bp.ResourceNotFoundError(err) {
						return resource.RetryableError(err)
					} else {
						return resource.NonRetryableError(err)
					}
				}
				return nil
			}); err != nil {
				return nil, "", err
			}

			status, err = bp.ObtainSdkValue("Status", ebs)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("volume status error, status:%s", status.(string))
				}
			}
			return ebs, status.(string), err
		},
	}
}

func (VestackVolumeService) WithResourceResponseHandlers(volume map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return volume, map[string]bp.ResponseConvert{
			"Size": {
				TargetField: "size",
				Convert:     sizeConvertFunc,
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackVolumeService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateVolume",
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.EbsClient.CreateVolumeCommon(call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.VolumeId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"available", "attached"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackVolumeService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	if resourceData.HasChanges("volume_name", "description", "delete_with_instance") {
		callbacks = append(callbacks, bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyVolumeAttribute",
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"volume_name": {
						TargetField: "VolumeName",
						ForceGet:    true,
					},
					"description": {
						TargetField: "Description",
					},
					"delete_with_instance": {
						TargetField: "DeleteWithInstance",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["VolumeId"] = d.Id()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					return s.Client.EbsClient.ModifyVolumeAttributeCommon(call.SdkParam)
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"available", "attached"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		})
	}

	if resourceData.HasChange("size") { // 调用新的 api
		callbacks = append(callbacks, bp.Callback{
			Call: bp.SdkCall{
				Action:      "ExtendVolume",
				ConvertMode: bp.RequestConvertIgnore,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["VolumeId"] = d.Id()
					(*call.SdkParam)["NewSize"] = d.Get("size")
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					return s.Client.EbsClient.ExtendVolumeCommon(call.SdkParam)
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"available", "attached"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		})
	}

	if resourceData.HasChange("volume_charge_type") {
		callbacks = append(callbacks, bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyVolumeChargeType",
				ConvertMode: bp.RequestConvertIgnore,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					oldV, newV := resourceData.GetChange("volume_charge_type")
					if oldV == "PrePaid" && newV == "PostPaid" {
						return false, errors.New("cannot convert PrePaid volume to PostPaid")
					}
					if d.Get("instance_id").(string) == "" {
						return false, errors.New("instance id cannot be empty")
					}

					(*call.SdkParam)["VolumeIds.1"] = d.Id()
					(*call.SdkParam)["DiskChargeType"] = "PrePaid"
					(*call.SdkParam)["AutoPay"] = true
					(*call.SdkParam)["InstanceId"] = d.Get("instance_id")
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp)
					logger.Debug(logger.RespFormat, call.Action, err)
					return resp, err
				},
				CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
					oldV, newV := resourceData.GetChange("volume_charge_type")
					if oldV == "PrePaid" && newV == "PostPaid" {
						return errors.New("cannot convert PrePaid volume to PostPaid")
					}
					if d.Get("instance_id").(string) == "" {
						return errors.New("instance id cannot be empty")
					}
					// retry modifyVolumeChargeType
					return re.Retry(15*time.Minute, func() *re.RetryError {
						data, callErr := s.ReadResource(d, d.Id())
						if callErr != nil {
							return re.NonRetryableError(fmt.Errorf("error on reading volume %q: %w", d.Id(), callErr))
						}
						// 计费方式已经转变成功
						if data["PayType"] == "pre" {
							return nil
						}
						// 计费方式还没有转换成功，尝试重新转换
						_, callErr = call.ExecuteCall(d, client, call)
						if callErr == nil {
							return nil
						}
						// 按量实例下挂载的云盘不支持按量转包年操作
						if strings.Contains(callErr.Error(), "ErrorInvalidEcsChargeType") {
							return re.NonRetryableError(callErr)
						}
						return re.RetryableError(callErr)
					})
				},
			},
		})
	}
	return callbacks
}

func (s *VestackVolumeService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteVolume",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"VolumeId": resourceData.Id(),
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				volume, err := s.ReadResource(d, d.Id())
				if err != nil {
					return false, err
				}
				status, err := bp.ObtainSdkValue("Status", volume)
				if err != nil {
					return false, err
				}
				if status != "available" {
					return false, fmt.Errorf(" Only volume with a status of `available` can be deleted. ")
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.EbsClient.DeleteVolumeCommon(call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				// 不能删除已挂载云盘
				if strings.Contains(baseErr.Error(), "Only volume with a status of `available` can be deleted.") {
					msg := fmt.Sprintf("error: %s\n msg: %s",
						baseErr.Error(),
						"For volume with a status of `attached`, please use `terraform state rm vestack_volume.resource_name` command to remove it from terraform state file and management.")
					return fmt.Errorf(msg)
				}
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading vpc on delete %q, %w", d.Id(), callErr))
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

func (s *VestackVolumeService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "VolumeIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "VolumeName",
		IdField:      "VolumeId",
		CollectField: "volumes",
		ResponseConverts: map[string]bp.ResponseConvert{
			"VolumeId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"Size": {
				TargetField: "size",
				Convert:     sizeConvertFunc,
			},
		},
	}
}

var sizeConvertFunc = func(i interface{}) interface{} {
	// Notice: the type of filed Size in openapi doc is size, but api return type is string
	size, ok := i.(string)
	if !ok {
		return i
	}
	res, err := strconv.Atoi(size)
	if err != nil {
		logger.Debug(logger.ReqFormat, "sizeConvertFunc", i)
		return i
	}
	return res
}

func (s *VestackVolumeService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "storage_ebs",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}

func (s *VestackVolumeService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "storage_ebs",
		ResourceType:         "volume",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}
