package tos_bucket_mirror_back

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackTosBucketMirrorBackService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewTosBucketMirrorBackService(c *bp.SdkClient) *VestackTosBucketMirrorBackService {
	return &VestackTosBucketMirrorBackService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackTosBucketMirrorBackService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackTosBucketMirrorBackService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return data, err
}

func (s *VestackTosBucketMirrorBackService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	tos := s.Client.BypassSvcClient
	var (
		ok bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	action := "GetBucketMirrorBack"
	logger.Debug(logger.ReqFormat, action, id)
	resp, err := tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     id,
		UrlParam: map[string]string{
			"mirror": "",
		},
	}, nil)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp, err)
	if data, ok = (*resp)[bp.BypassResponse].(map[string]interface{}); !ok {
		return data, errors.New("GetBucketMirrorBack Resp is not map")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("tos_bucket_mirror_back %s not exist ", id)
	}

	data["BucketName"] = id
	if v, exist := data["Rules"]; exist {
		if rules, ok := v.([]interface{}); ok {
			for _, rule := range rules {
				if rule1, ok := rule.(map[string]interface{}); ok {
					if v1, exist1 := rule1["Condition"]; exist1 {
						if condition, ok1 := v1.(map[string]interface{}); ok1 {
							rule1["Condition"] = []interface{}{condition}
						}
					}
					if v1, exist1 := rule1["Redirect"]; exist1 {
						if redirect, ok1 := v1.(map[string]interface{}); ok1 {
							if v2, exist2 := redirect["PublicSource"]; exist2 {
								if publicSource, ok2 := v2.(map[string]interface{}); ok2 {
									if v3, exist3 := publicSource["SourceEndpoint"]; exist3 {
										if sourceEndpoint, ok3 := v3.(map[string]interface{}); ok3 {
											publicSource["SourceEndpoint"] = []interface{}{sourceEndpoint}
										}
									}
									redirect["PublicSource"] = []interface{}{publicSource}
								}
							}
							if v2, exist2 := redirect["MirrorHeader"]; exist2 {
								if mirrorHeader, ok2 := v2.(map[string]interface{}); ok2 {
									redirect["MirrorHeader"] = []interface{}{mirrorHeader}
								}
							}
							if v2, exist2 := redirect["FetchHeaderToMetaDataRules"]; exist2 {
								if fetchHeaderToMetaDataRules, ok2 := v2.(map[string]interface{}); ok2 {
									redirect["FetchHeaderToMetaDataRules"] = []interface{}{fetchHeaderToMetaDataRules}
								}
							}
							if v2, exist2 := redirect["Transform"]; exist2 {
								if transform, ok2 := v2.(map[string]interface{}); ok2 {
									if v3, exist3 := transform["ReplaceKeyPrefix"]; exist3 {
										if replaceKeyPrefix, ok3 := v3.(map[string]interface{}); ok3 {
											transform["ReplaceKeyPrefix"] = []interface{}{replaceKeyPrefix}
										}
									}
									redirect["Transform"] = []interface{}{transform}
								}
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

func (s *VestackTosBucketMirrorBackService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *VestackTosBucketMirrorBackService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
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

func (s *VestackTosBucketMirrorBackService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateMirrorBack(resourceData, resource, false)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketMirrorBackService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateMirrorBack(resourceData, resource, true)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketMirrorBackService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteBucketMirrorBack",
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
						"mirror": "",
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
							return resource.NonRetryableError(fmt.Errorf("error on reading tos bucket mirror_back on delete %q, %w", s.ReadResourceId(d.Id()), callErr))
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

func (s *VestackTosBucketMirrorBackService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackTosBucketMirrorBackService) createOrUpdateMirrorBack(resourceData *schema.ResourceData, resource *schema.Resource, isUpdate bool) bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "PutBucketMirrorBack",
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
						"redirect": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "Redirect",
							ForceGet:    isUpdate,
							NextLevelConvert: map[string]bp.RequestConvert{
								"redirect_type": {
									ConvertType: bp.ConvertDefault,
									TargetField: "RedirectType",
									ForceGet:    isUpdate,
								},
								"fetch_source_on_redirect": {
									ConvertType: bp.ConvertDefault,
									TargetField: "FetchSourceOnRedirect",
									ForceGet:    isUpdate,
								},
								"fetch_source_on_redirect_with_query": {
									ConvertType: bp.ConvertDefault,
									TargetField: "FetchSourceOnRedirectWithQuery",
									ForceGet:    isUpdate,
								},
								"pass_query": {
									ConvertType: bp.ConvertDefault,
									TargetField: "PassQuery",
									ForceGet:    isUpdate,
								},
								"follow_redirect": {
									ConvertType: bp.ConvertDefault,
									TargetField: "FollowRedirect",
									ForceGet:    isUpdate,
								},
								"mirror_header": {
									ConvertType: bp.ConvertJsonObject,
									TargetField: "MirrorHeader",
									ForceGet:    isUpdate,
									NextLevelConvert: map[string]bp.RequestConvert{
										"pass_all": {
											ConvertType: bp.ConvertDefault,
											TargetField: "PassAll",
											ForceGet:    isUpdate,
										},
										"pass": {
											ConvertType: bp.ConvertJsonArray,
											TargetField: "Pass",
											ForceGet:    isUpdate,
										},
										"remove": {
											ConvertType: bp.ConvertJsonArray,
											TargetField: "Remove",
											ForceGet:    isUpdate,
										},
										"set": {
											ConvertType: bp.ConvertJsonObjectArray,
											TargetField: "Set",
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
									},
								},
								"public_source": {
									ConvertType: bp.ConvertJsonObject,
									TargetField: "PublicSource",
									ForceGet:    isUpdate,
									NextLevelConvert: map[string]bp.RequestConvert{
										"source_endpoint": {
											ConvertType: bp.ConvertJsonObject,
											TargetField: "SourceEndpoint",
											ForceGet:    isUpdate,
											NextLevelConvert: map[string]bp.RequestConvert{
												"primary": {
													ConvertType: bp.ConvertJsonArray,
													TargetField: "Primary",
													ForceGet:    isUpdate,
												},
												"follower": {
													ConvertType: bp.ConvertJsonArray,
													TargetField: "Follower",
													ForceGet:    isUpdate,
												},
											},
										},
										"fixed_endpoint": {
											ConvertType: bp.ConvertDefault,
											TargetField: "FixedEndpoint",
											ForceGet:    isUpdate,
										},
									},
								},
								"private_source": {
									ConvertType: bp.ConvertJsonObject,
									TargetField: "PrivateSource",
									ForceGet:    isUpdate,
									NextLevelConvert: map[string]bp.RequestConvert{
										"source_endpoint": {
											ConvertType: bp.ConvertJsonObject,
											TargetField: "SourceEndpoint",
											ForceGet:    isUpdate,
											NextLevelConvert: map[string]bp.RequestConvert{
												"primary": {
													ConvertType: bp.ConvertJsonObjectArray,
													TargetField: "Primary",
													ForceGet:    isUpdate,
													NextLevelConvert: map[string]bp.RequestConvert{
														"endpoint": {
															ConvertType: bp.ConvertDefault,
															TargetField: "Endpoint",
															ForceGet:    isUpdate,
														},
														"bucket_name": {
															ConvertType: bp.ConvertDefault,
															TargetField: "BucketName",
															ForceGet:    isUpdate,
														},
														"credential_provider": {
															ConvertType: bp.ConvertJsonObject,
															TargetField: "CredentialProvider",
															ForceGet:    isUpdate,
															NextLevelConvert: map[string]bp.RequestConvert{
																"role": {
																	ConvertType: bp.ConvertDefault,
																	TargetField: "Role",
																	ForceGet:    isUpdate,
																},
															},
														},
													},
												},
												"follower": {
													ConvertType: bp.ConvertJsonObjectArray,
													TargetField: "Follower",
													ForceGet:    isUpdate,
													NextLevelConvert: map[string]bp.RequestConvert{
														"endpoint": {
															ConvertType: bp.ConvertDefault,
															TargetField: "Endpoint",
															ForceGet:    isUpdate,
														},
														"bucket_name": {
															ConvertType: bp.ConvertDefault,
															TargetField: "BucketName",
															ForceGet:    isUpdate,
														},
														"credential_provider": {
															ConvertType: bp.ConvertJsonObject,
															TargetField: "CredentialProvider",
															ForceGet:    isUpdate,
															NextLevelConvert: map[string]bp.RequestConvert{
																"role": {
																	ConvertType: bp.ConvertDefault,
																	TargetField: "Role",
																	ForceGet:    isUpdate,
																},
															},
														},
													},
												},
											},
										},
									},
								},
								"transform": {
									ConvertType: bp.ConvertJsonObject,
									TargetField: "Transform",
									ForceGet:    isUpdate,
									NextLevelConvert: map[string]bp.RequestConvert{
										"with_key_prefix": {
											ConvertType: bp.ConvertDefault,
											TargetField: "WithKeyPrefix",
											ForceGet:    isUpdate,
										},
										"with_key_suffix": {
											ConvertType: bp.ConvertDefault,
											TargetField: "WithKeySuffix",
											ForceGet:    isUpdate,
										},
										"replace_key_prefix": {
											ConvertType: bp.ConvertJsonObject,
											TargetField: "ReplaceKeyPrefix",
											ForceGet:    isUpdate,
											NextLevelConvert: map[string]bp.RequestConvert{
												"key_prefix": {
													ConvertType: bp.ConvertDefault,
													TargetField: "KeyPrefix",
													ForceGet:    isUpdate,
												},
												"replace_with": {
													ConvertType: bp.ConvertDefault,
													TargetField: "ReplaceWith",
													ForceGet:    isUpdate,
												},
											},
										},
										"remove_key_prefix": {
											ConvertType: bp.ConvertDefault,
											TargetField: "RemoveKeyPrefix",
											ForceGet:    isUpdate,
										},
									},
								},
							},
						},
						"condition": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "Condition",
							ForceGet:    isUpdate,
							NextLevelConvert: map[string]bp.RequestConvert{
								"http_code": {
									ConvertType: bp.ConvertDefault,
									TargetField: "HttpCode",
									ForceGet:    isUpdate,
								},
								"key_prefix": {
									ConvertType: bp.ConvertDefault,
									TargetField: "KeyPrefix",
									ForceGet:    isUpdate,
								},
								"key_suffix": {
									ConvertType: bp.ConvertDefault,
									TargetField: "KeySuffix",
									ForceGet:    isUpdate,
								},
								"allow_host": {
									ConvertType: bp.ConvertJsonArray,
									TargetField: "AllowHost",
									ForceGet:    isUpdate,
								},
								"http_method": {
									ConvertType: bp.ConvertJsonArray,
									TargetField: "HttpMethod",
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
						"mirror": "",
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

func (s *VestackTosBucketMirrorBackService) ReadResourceId(id string) string {
	return id
}
