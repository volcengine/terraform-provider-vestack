package tos_bucket_lifecycle

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackTosBucketLifecycleService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewTosBucketLifecycleService(c *bp.SdkClient) *VestackTosBucketLifecycleService {
	return &VestackTosBucketLifecycleService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackTosBucketLifecycleService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackTosBucketLifecycleService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return data, err
}

func (s *VestackTosBucketLifecycleService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	tos := s.Client.BypassSvcClient
	var (
		ok bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	action := "GetBucketLifecycle"
	logger.Debug(logger.ReqFormat, action, id)
	resp, err := tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     id,
		UrlParam: map[string]string{
			"lifecycle": "",
		},
	}, nil)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp, err)
	if data, ok = (*resp)[bp.BypassResponse].(map[string]interface{}); !ok {
		return data, errors.New("GetBucketLifecycle Resp is not map")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("tos_bucket_lifecycle %s not exist ", id)
	}

	data["BucketName"] = id
	if v, exist := data["Rules"]; exist {
		if rules, ok := v.([]interface{}); ok {
			for _, rule := range rules {
				if rule1, ok := rule.(map[string]interface{}); ok {
					if v1, exist1 := rule1["Filter"]; exist1 {
						if filter, ok1 := v1.(map[string]interface{}); ok1 {
							rule1["Filter"] = []interface{}{filter}
						}
					}
					if v1, exist1 := rule1["Expiration"]; exist1 {
						if expiration, ok1 := v1.(map[string]interface{}); ok1 {
							rule1["Expiration"] = []interface{}{expiration}
						}
					}
					if v1, exist1 := rule1["NoncurrentVersionExpiration"]; exist1 {
						if noncurrentVersionExpiration, ok1 := v1.(map[string]interface{}); ok1 {
							rule1["NoncurrentVersionExpiration"] = []interface{}{noncurrentVersionExpiration}
						}
					}
					if v1, exist1 := rule1["Transitions"]; exist1 {
						if transitions, ok1 := v1.(map[string]interface{}); ok1 {
							rule1["Transitions"] = []interface{}{transitions}
						}
					}
					if v1, exist1 := rule1["AbortIncompleteMultipartUpload"]; exist1 {
						if abortIncompleteMultipartUpload, ok1 := v1.(map[string]interface{}); ok1 {
							rule1["AbortIncompleteMultipartUpload"] = []interface{}{abortIncompleteMultipartUpload}
						}
					}
				}
			}
		}
	}

	return data, err
}

func (s *VestackTosBucketLifecycleService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *VestackTosBucketLifecycleService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"Rules": {
				TargetField: "rules",
			},
			"ID": {
				TargetField: "id",
			},
			"NoncurrentVersionExpiration": {
				TargetField: "non_current_version_expiration",
			},
			"NoncurrentVersionTransitions": {
				TargetField: "non_current_version_transitions",
			},
			"NoncurrentDays": {
				TargetField: "non_current_days",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackTosBucketLifecycleService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateLifecycle(resourceData, resource, false)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketLifecycleService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateLifecycle(resourceData, resource, true)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketLifecycleService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteBucketLifecycle",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["BucketName"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					ContentType: bp.ApplicationJSON,
					HttpMethod:  bp.DELETE,
					Domain:      (*call.SdkParam)["BucketName"].(string),
					UrlParam: map[string]string{
						"lifecycle": "",
					},
				}, nil)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				return resource.Retry(5*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading tos bucket lifecycle on delete %q, %w", s.ReadResourceId(d.Id()), callErr))
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

func (s *VestackTosBucketLifecycleService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackTosBucketLifecycleService) createOrUpdateLifecycle(resourceData *schema.ResourceData, resource *schema.Resource, isUpdate bool) bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "PutBucketLifecycle",
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
				"rules": {
					ConvertType: bp.ConvertJsonObjectArray,
					TargetField: "Rules",
					ForceGet:    isUpdate,
					NextLevelConvert: map[string]bp.RequestConvert{
						"id": {
							ConvertType: bp.ConvertDefault,
							TargetField: "ID",
							ForceGet:    isUpdate,
						},
						"prefix": {
							ConvertType: bp.ConvertDefault,
							TargetField: "Prefix",
							ForceGet:    isUpdate,
						},
						"status": {
							ConvertType: bp.ConvertDefault,
							TargetField: "Status",
							ForceGet:    isUpdate,
						},
						"expiration": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "Expiration",
							ForceGet:    isUpdate,
							NextLevelConvert: map[string]bp.RequestConvert{
								"days": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Days",
									ForceGet:    isUpdate,
								},
								"date": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Date",
									ForceGet:    isUpdate,
								},
							},
						},
						"transitions": {
							ConvertType: bp.ConvertJsonObjectArray,
							TargetField: "Transitions",
							ForceGet:    isUpdate,
							NextLevelConvert: map[string]bp.RequestConvert{
								"days": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Days",
									ForceGet:    isUpdate,
								},
								"date": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Date",
									ForceGet:    isUpdate,
								},
								"storage_class": {
									ConvertType: bp.ConvertDefault,
									TargetField: "StorageClass",
									ForceGet:    isUpdate,
								},
							},
						},
						"non_current_version_expiration": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "NonCurrentVersionExpiration",
							ForceGet:    isUpdate,
							NextLevelConvert: map[string]bp.RequestConvert{
								"non_current_days": {
									ConvertType: bp.ConvertDefault,
									TargetField: "NonCurrentDays",
									ForceGet:    isUpdate,
								},
							},
						},
						"non_current_version_transitions": {
							ConvertType: bp.ConvertJsonObjectArray,
							TargetField: "NonCurrentVersionTransitions",
							ForceGet:    isUpdate,
							NextLevelConvert: map[string]bp.RequestConvert{
								"non_current_days": {
									ConvertType: bp.ConvertDefault,
									TargetField: "NonCurrentDays",
									ForceGet:    isUpdate,
								},
								"storage_class": {
									ConvertType: bp.ConvertDefault,
									TargetField: "StorageClass",
									ForceGet:    isUpdate,
								},
							},
						},
						"tags": {
							ConvertType: bp.ConvertJsonObjectArray,
							TargetField: "Tags",
							ForceGet:    isUpdate,
							NextLevelConvert: map[string]bp.RequestConvert{
								"key": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Key",
									ForceGet:    isUpdate,
								},
								"value": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Value",
									ForceGet:    isUpdate,
								},
							},
						},
						"filter": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "Filter",
							ForceGet:    isUpdate,
							NextLevelConvert: map[string]bp.RequestConvert{
								"object_size_greater_than": {
									ConvertType: bp.ConvertDefault,
									TargetField: "ObjectSizeGreaterThan",
									ForceGet:    isUpdate,
								},
								"object_size_less_than": {
									ConvertType: bp.ConvertDefault,
									TargetField: "ObjectSizeLessThan",
									ForceGet:    isUpdate,
								},
							},
						},
						"abort_incomplete_multipart_upload": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "AbortIncompleteMultipartUpload",
							ForceGet:    isUpdate,
							NextLevelConvert: map[string]bp.RequestConvert{
								"days_after_initiation": {
									ConvertType: bp.ConvertDefault,
									TargetField: "DaysAfterInitiation",
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
				resp, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					ContentType: bp.ApplicationJSON,
					HttpMethod:  bp.PUT,
					Domain:      (*call.SdkParam)[bp.BypassDomain].(string),
					UrlParam: map[string]string{
						"lifecycle": "",
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

func (s *VestackTosBucketLifecycleService) ReadResourceId(id string) string {
	return id
}
