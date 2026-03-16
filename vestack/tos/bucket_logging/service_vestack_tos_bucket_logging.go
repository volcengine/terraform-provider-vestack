package tos_bucket_logging

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackTosBucketLoggingService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewTosBucketLoggingService(c *bp.SdkClient) *VestackTosBucketLoggingService {
	return &VestackTosBucketLoggingService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackTosBucketLoggingService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackTosBucketLoggingService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return data, err
}

func (s *VestackTosBucketLoggingService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	tos := s.Client.BypassSvcClient
	var (
		ok bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	action := "GetBucketLogging"
	logger.Debug(logger.ReqFormat, action, id)
	resp, err := tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     id,
		UrlParam: map[string]string{
			"logging": "",
		},
	}, nil)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp, err)
	if data, ok = (*resp)[bp.BypassResponse].(map[string]interface{}); !ok {
		return data, errors.New("GetBucketLogging Resp is not map")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("tos_bucket_logging %s not exist ", id)
	}

	data["BucketName"] = id
	return data, err
}

func (s *VestackTosBucketLoggingService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *VestackTosBucketLoggingService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"LoggingEnabled": {
				TargetField: "logging_enabled",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackTosBucketLoggingService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateLogging(resourceData, resource, false)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketLoggingService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateLogging(resourceData, resource, true)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketLoggingService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "PutBucketLogging",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["BucketName"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				params := map[string]interface{}{}
				resp, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					ContentType: bp.ApplicationJSON,
					HttpMethod:  bp.PUT,
					Domain:      (*call.SdkParam)["BucketName"].(string),
					UrlParam: map[string]string{
						"logging": "",
					},
				}, &params)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				return resource.Retry(5*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading tos bucket logging on delete %q, %w", s.ReadResourceId(d.Id()), callErr))
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
				return d.Get("bucket_name").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackTosBucketLoggingService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackTosBucketLoggingService) createOrUpdateLogging(resourceData *schema.ResourceData, resource *schema.Resource, isUpdate bool) bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "PutBucketLogging",
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
				"logging_enabled": {
					ConvertType: bp.ConvertJsonObject,
					TargetField: "LoggingEnabled",
					ForceGet:    isUpdate,
					NextLevelConvert: map[string]bp.RequestConvert{
						"target_bucket": {
							ConvertType: bp.ConvertDefault,
							TargetField: "TargetBucket",
							ForceGet:    isUpdate,
						},
						"target_prefix": {
							ConvertType: bp.ConvertDefault,
							TargetField: "TargetPrefix",
							ForceGet:    isUpdate,
						},
						"role": {
							ConvertType: bp.ConvertDefault,
							TargetField: "Role",
							ForceGet:    isUpdate,
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
				resp, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					ContentType: bp.ApplicationJSON,
					HttpMethod:  bp.PUT,
					Domain:      (*call.SdkParam)[bp.BypassDomain].(string),
					UrlParam: map[string]string{
						"logging": "",
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

func (s *VestackTosBucketLoggingService) ReadResourceId(id string) string {
	return id
}
