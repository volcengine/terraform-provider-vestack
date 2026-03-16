package tos_bucket_website

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackTosBucketWebsiteService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewTosBucketWebsiteService(c *bp.SdkClient) *VestackTosBucketWebsiteService {
	return &VestackTosBucketWebsiteService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackTosBucketWebsiteService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackTosBucketWebsiteService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return data, err
}

func (s *VestackTosBucketWebsiteService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	tos := s.Client.BypassSvcClient
	var (
		ok bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	action := "GetBucketWebsite"
	logger.Debug(logger.ReqFormat, action, id)
	resp, err := tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     id,
		UrlParam: map[string]string{
			"website": "",
		},
	}, nil)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp, err)
	if data, ok = (*resp)[bp.BypassResponse].(map[string]interface{}); !ok {
		return data, errors.New("GetBucketWebsite Resp is not map")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("tos_bucket_website %s not exist ", id)
	}

	data["BucketName"] = id

	if v, exist := data["RoutingRules"]; exist {
		if routingRules, ok := v.([]interface{}); ok {
			for _, rule := range routingRules {
				if rule1, ok := rule.(map[string]interface{}); ok {
					if v1, exist1 := rule1["Condition"]; exist1 {
						if condition, ok1 := v1.(map[string]interface{}); ok1 {
							rule1["Condition"] = []interface{}{condition}
						}
					}
					if v1, exist1 := rule1["Redirect"]; exist1 {
						if redirect, ok1 := v1.(map[string]interface{}); ok1 {
							if v2, ok := redirect["ReplaceKeyPrefixWith"]; ok {
								redirect["ReplaceKeyPrefixWith"] = v2
							}
							if v2, ok := redirect["ReplaceKeyWith"]; ok {
								redirect["ReplaceKeyWith"] = v2
							}
							rule1["Redirect"] = []interface{}{redirect}
						}
					}
				}
			}
		}
	}

	return data, err
}

func (s *VestackTosBucketWebsiteService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *VestackTosBucketWebsiteService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"IndexDocument": {
				TargetField: "index_document",
			},
			"ErrorDocument": {
				TargetField: "error_document",
			},
			"RedirectAllRequestsTo": {
				TargetField: "redirect_all_requests_to",
			},
			"RoutingRules": {
				TargetField: "routing_rules",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackTosBucketWebsiteService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateWebsite(resourceData, resource, false)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketWebsiteService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateWebsite(resourceData, resource, true)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketWebsiteService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteBucketWebsite",
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
						"website": "",
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
							return resource.NonRetryableError(fmt.Errorf("error on reading tos bucket website on delete %q, %w", s.ReadResourceId(d.Id()), callErr))
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

func (s *VestackTosBucketWebsiteService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackTosBucketWebsiteService) createOrUpdateWebsite(resourceData *schema.ResourceData, resource *schema.Resource, isUpdate bool) bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "PutBucketWebsite",
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
				"index_document": {
					ConvertType: bp.ConvertJsonObject,
					TargetField: "IndexDocument",
					ForceGet:    isUpdate,
					NextLevelConvert: map[string]bp.RequestConvert{
						"suffix": {
							ConvertType: bp.ConvertDefault,
							TargetField: "Suffix",
							ForceGet:    isUpdate,
						},
						"support_sub_dir": {
							ConvertType: bp.ConvertDefault,
							TargetField: "SupportSubDir",
							ForceGet:    isUpdate,
						},
					},
				},
				"error_document": {
					ConvertType: bp.ConvertJsonObject,
					TargetField: "ErrorDocument",
					ForceGet:    isUpdate,
					NextLevelConvert: map[string]bp.RequestConvert{
						"key": {
							ConvertType: bp.ConvertDefault,
							TargetField: "Key",
							ForceGet:    isUpdate,
						},
					},
				},
				"redirect_all_requests_to": {
					ConvertType: bp.ConvertJsonObject,
					TargetField: "RedirectAllRequestsTo",
					ForceGet:    isUpdate,
					NextLevelConvert: map[string]bp.RequestConvert{
						"host_name": {
							ConvertType: bp.ConvertDefault,
							TargetField: "HostName",
							ForceGet:    isUpdate,
						},
						"protocol": {
							ConvertType: bp.ConvertDefault,
							TargetField: "Protocol",
							ForceGet:    isUpdate,
						},
					},
				},
				"routing_rules": {
					ConvertType: bp.ConvertJsonObjectArray,
					TargetField: "RoutingRules",
					NextLevelConvert: map[string]bp.RequestConvert{
						"condition": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "Condition",
							ForceGet:    isUpdate,
							NextLevelConvert: map[string]bp.RequestConvert{
								"key_prefix_equals": {
									ConvertType: bp.ConvertDefault,
									TargetField: "KeyPrefixEquals",
									ForceGet:    isUpdate,
								},
								"http_error_code_returned_equals": {
									ConvertType: bp.ConvertDefault,
									TargetField: "HttpErrorCodeReturnedEquals",
									ForceGet:    isUpdate,
								},
							},
						},
						"redirect": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "Redirect",
							NextLevelConvert: map[string]bp.RequestConvert{
								"protocol": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Protocol",
									ForceGet:    isUpdate,
								},
								"host_name": {
									ConvertType: bp.ConvertDefault,
									TargetField: "HostName",
									ForceGet:    isUpdate,
								},
								"replace_key_with": {
									ConvertType: bp.ConvertDefault,
									TargetField: "ReplaceKeyWith",
								},
								"replace_key_prefix_with": {
									ConvertType: bp.ConvertDefault,
									TargetField: "ReplaceKeyPrefixWith",
								},
								"http_redirect_code": {
									ConvertType: bp.ConvertDefault,
									TargetField: "HttpRedirectCode",
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
						"website": "",
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

func (s *VestackTosBucketWebsiteService) ReadResourceId(id string) string {
	return id
}
