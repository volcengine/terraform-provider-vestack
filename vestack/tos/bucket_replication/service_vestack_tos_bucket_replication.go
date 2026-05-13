package tos_bucket_replication

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackTosBucketReplicationService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewTosBucketReplicationService(c *bp.SdkClient) *VestackTosBucketReplicationService {
	return &VestackTosBucketReplicationService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackTosBucketReplicationService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackTosBucketReplicationService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return data, err
}

func (s *VestackTosBucketReplicationService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	tos := s.Client.BypassSvcClient
	var (
		ok bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	action := "GetBucketReplication"
	logger.Debug(logger.ReqFormat, action, id)
	resp, err := tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     id,
		UrlParam: map[string]string{
			"replication": "",
			"progress":    "",
		},
	}, nil)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp, err)
	if data, ok = (*resp)[bp.BypassResponse].(map[string]interface{}); !ok {
		return data, errors.New("GetBucketReplication Resp is not map")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("tos_bucket_replication %s not exist ", id)
	}

	data["BucketName"] = id
	if v, exist := data["Rules"]; exist {
		if rules, ok := v.([]interface{}); ok {
			for _, rule := range rules {
				if rule1, ok := rule.(map[string]interface{}); ok {
					if v1, exist1 := rule1["Destination"]; exist1 {
						if destination, ok1 := v1.(map[string]interface{}); ok1 {
							rule1["Destination"] = []interface{}{destination}
						}
					}
					if v1, exist1 := rule1["AccessControlTranslation"]; exist1 {
						if accessControlTranslation, ok1 := v1.(map[string]interface{}); ok1 {
							rule1["AccessControlTranslation"] = []interface{}{accessControlTranslation}
						}
					}
				}
			}
		}
	}

	return data, err
}

func (s *VestackTosBucketReplicationService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *VestackTosBucketReplicationService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"Role": {
				TargetField: "role",
			},
			"Rules": {
				TargetField: "rules",
			},
			"ID": {
				TargetField: "id",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackTosBucketReplicationService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateReplication(resourceData, resource, false)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketReplicationService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateReplication(resourceData, resource, true)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketReplicationService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteBucketReplication",
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
						"replication": "",
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
							return resource.NonRetryableError(fmt.Errorf("error on reading tos bucket replication on delete %q, %w", s.ReadResourceId(d.Id()), callErr))
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

func (s *VestackTosBucketReplicationService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackTosBucketReplicationService) createOrUpdateReplication(resourceData *schema.ResourceData, resource *schema.Resource, isUpdate bool) bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "PutBucketReplication",
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
				"role": {
					ConvertType: bp.ConvertDefault,
					TargetField: "Role",
					ForceGet:    isUpdate,
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
						"status": {
							ConvertType: bp.ConvertDefault,
							TargetField: "Status",
							ForceGet:    isUpdate,
						},
						"prefix_set": {
							ConvertType: bp.ConvertJsonArray,
							TargetField: "PrefixSet",
							ForceGet:    isUpdate,
						},
						"destination": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "Destination",
							ForceGet:    isUpdate,
							NextLevelConvert: map[string]bp.RequestConvert{
								"bucket": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Bucket",
									ForceGet:    isUpdate,
								},
								"location": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Location",
									ForceGet:    isUpdate,
								},
								"storage_class": {
									ConvertType: bp.ConvertDefault,
									TargetField: "StorageClass",
									ForceGet:    isUpdate,
								},
								"storage_class_inherit_directive": {
									ConvertType: bp.ConvertDefault,
									TargetField: "StorageClassInheritDirective",
									ForceGet:    isUpdate,
								},
							},
						},
						"historical_object_replication": {
							ConvertType: bp.ConvertDefault,
							TargetField: "HistoricalObjectReplication",
							ForceGet:    isUpdate,
						},
						"access_control_translation": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "AccessControlTranslation",
							ForceGet:    isUpdate,
							NextLevelConvert: map[string]bp.RequestConvert{
								"owner": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Owner",
									ForceGet:    isUpdate,
								},
							},
						},
						"transfer_type": {
							ConvertType: bp.ConvertDefault,
							TargetField: "TransferType",
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
						"replication": "",
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

func (s *VestackTosBucketReplicationService) ReadResourceId(id string) string {
	return id
}
