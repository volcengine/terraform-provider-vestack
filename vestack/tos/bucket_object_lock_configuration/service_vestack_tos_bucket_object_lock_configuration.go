package tos_bucket_object_lock_configuration

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackTosBucketObjectLockConfigurationService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewTosBucketObjectLockConfigurationService(c *bp.SdkClient) *VestackTosBucketObjectLockConfigurationService {
	return &VestackTosBucketObjectLockConfigurationService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackTosBucketObjectLockConfigurationService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackTosBucketObjectLockConfigurationService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return data, err
}

func (s *VestackTosBucketObjectLockConfigurationService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	tos := s.Client.BypassSvcClient
	var (
		ok bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	action := "GetBucketObjectLockConfiguration"
	logger.Debug(logger.ReqFormat, action, id)
	resp, err := tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     id,
		UrlParam: map[string]string{
			"object-lock": "",
		},
	}, nil)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp, err)
	if data, ok = (*resp)[bp.BypassResponse].(map[string]interface{}); !ok {
		return data, errors.New("GetBucketObjectLockConfiguration Resp is not map")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("tos_bucket_object_lock_configuration %s not exist ", id)
	}

	data["BucketName"] = id
	if v, exist := data["Rule"]; exist {
		if rule, ok := v.(map[string]interface{}); ok {
			if v1, exist1 := rule["DefaultRetention"]; exist1 {
				if defaultRetention, ok1 := v1.(map[string]interface{}); ok1 {
					rule["DefaultRetention"] = []interface{}{defaultRetention}
				}
			}
		}
	}

	return data, err
}

func (s *VestackTosBucketObjectLockConfigurationService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *VestackTosBucketObjectLockConfigurationService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"ObjectLockConfiguration": {
				TargetField: "object_lock_configuration",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackTosBucketObjectLockConfigurationService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateObjectLockConfiguration(resourceData, resource, false)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketObjectLockConfigurationService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateObjectLockConfiguration(resourceData, resource, true)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketObjectLockConfigurationService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	logger.Debug(logger.ReqFormat, "RemoveResource", "Remove VestackTosBucketObjectLockConfigurationService dose not have callback")
	return []bp.Callback{}
}

func (s *VestackTosBucketObjectLockConfigurationService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackTosBucketObjectLockConfigurationService) createOrUpdateObjectLockConfiguration(resourceData *schema.ResourceData, resource *schema.Resource, isUpdate bool) bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "PutBucketObjectLockConfiguration",
			ConvertMode:     bp.RequestConvertInConvert,
			ContentType:     bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"bucket_name": {
					ConvertType: bp.ConvertDefault,
					TargetField: "BucketName",
					SpecialParam: &bp.SpecialParam{
						Type: bp.DomainParam,
					},
					ForceGet: isUpdate,
				},
				"rule": {
					ConvertType: bp.ConvertJsonObject,
					TargetField: "Rule",
					ForceGet:    isUpdate,
					NextLevelConvert: map[string]bp.RequestConvert{
						"default_retention": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "DefaultRetention",
							ForceGet:    isUpdate,
							NextLevelConvert: map[string]bp.RequestConvert{
								"mode": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Mode",
									ForceGet:    isUpdate,
								},
								"days": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Days",
									ForceGet:    isUpdate,
								},
								"years": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Years",
									ForceGet:    isUpdate,
								},
							},
						},
					},
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				var sourceParam map[string]interface{}
				sourceParam, err := bp.SortAndStartTransJson((*call.SdkParam)[bp.BypassParam].(map[string]interface{}))
				if err != nil {
					return false, err
				}
				(*call.SdkParam)[bp.BypassParam] = sourceParam

				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				param := (*call.SdkParam)[bp.BypassParam].(map[string]interface{})
				param["ObjectLockEnabled"] = "Enabled"
				resp, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					ContentType: bp.ApplicationJSON,
					HttpMethod:  bp.PUT,
					Domain:      (*call.SdkParam)[bp.BypassDomain].(string),
					UrlParam: map[string]string{
						"object-lock": "",
					},
				}, &param)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId((*call.SdkParam)[bp.BypassDomain].(string))
				return nil
			},
		},
	}

	return callback
}

func (s *VestackTosBucketObjectLockConfigurationService) ReadResourceId(id string) string {
	return id
}
