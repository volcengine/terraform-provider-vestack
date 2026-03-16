package tos_bucket_notification

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

type VestackTosBucketNotificationService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewTosBucketNotificationService(c *bp.SdkClient) *VestackTosBucketNotificationService {
	return &VestackTosBucketNotificationService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackTosBucketNotificationService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackTosBucketNotificationService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return data, err
}

func (s *VestackTosBucketNotificationService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	tos := s.Client.BypassSvcClient
	var (
		ok bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return data, fmt.Errorf("invalid tos notification id: %s", id)
	}

	action := "GetBucketNotificationV2"
	logger.Debug(logger.ReqFormat, action, id)
	resp, err := tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     ids[0],
		UrlParam: map[string]string{
			"notification_v2": "",
		},
	}, nil)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp, err)
	if data, ok = (*resp)[bp.BypassResponse].(map[string]interface{}); !ok {
		return data, errors.New("GetBucketNotificationV2 Resp is not map ")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("tos_bucket_notification %s not exist ", id)
	}

	rules, ok := data["Rules"].([]interface{})
	if !ok {
		return data, fmt.Errorf("tos_bucket_notification %s rules is not slice ", id)
	}
	if len(rules) == 0 {
		return data, fmt.Errorf("tos_bucket_notification %s rules not exist ", id)
	}

	// 根据 rule id 查找对应 rule
	rule := make(map[string]interface{})
	for _, v := range rules {
		ruleMap, ok := v.(map[string]interface{})
		if !ok {
			return data, fmt.Errorf("tos_bucket_notification %s rule is not map ", id)
		}
		if ids[1] == ruleMap["RuleId"] {
			rule = ruleMap
			break
		}
	}
	if len(rule) == 0 {
		return data, fmt.Errorf("tos_bucket_notification %s rule not exist ", id)
	}

	if destination, exist := rule["Destination"]; exist {
		if destinationMap, ok := destination.(map[string]interface{}); ok {
			rule["Destination"] = []interface{}{destinationMap}
		}
	}
	if filter, exist := rule["Filter"]; exist {
		if filterMap, ok := filter.(map[string]interface{}); ok {
			if tosKey, exist := filterMap["TOSKey"]; exist {
				if tosKeyMap, ok := tosKey.(map[string]interface{}); ok {
					filterMap["TOSKey"] = []interface{}{tosKeyMap}
				}
			}
			rule["Filter"] = []interface{}{filterMap}
		}
	}
	data["Rules"] = []interface{}{rule}
	data["BucketName"] = ids[0]

	return data, err
}

func (s *VestackTosBucketNotificationService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (VestackTosBucketNotificationService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"TOSKey": {
				TargetField: "tos_key",
			},
			"VeFaaS": {
				TargetField: "ve_faas",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackTosBucketNotificationService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	//create inventory
	callback := s.createOrUpdateNotification(resourceData, resource, false)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketNotificationService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	//create inventory
	callback := s.createOrUpdateNotification(resourceData, resource, true)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketNotificationService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	ids := strings.Split(resourceData.Id(), ":")
	if len(ids) != 2 {
		return []bp.Callback{{
			Err: fmt.Errorf("invalid tos bucket notification id: %s", resourceData.Id()),
		}}
	}
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "PutBucketNotificationV2",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["BucketName"] = ids[0]
				(*call.SdkParam)[bp.BypassParam] = make(map[string]interface{})
				//version := d.Get("version").(string)
				//(*call.SdkParam)[bp.BypassParam] = map[string]interface{}{
				//	"Version": version,
				//}
				return true, nil
			},
			AfterLocked: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) error {
				// 获取存量 rules 信息
				action := "GetBucketNotificationV2"
				logger.Debug(logger.ReqFormat, action, d.Get("bucket_name"))
				data, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					HttpMethod: bp.GET,
					Domain:     ids[0],
					UrlParam: map[string]string{
						"notification_v2": "",
					},
				}, nil)
				logger.Debug(logger.RespFormat, action, data, err)
				if err != nil {
					return err
				}

				v, _ := bp.ObtainSdkValue("Rules", (*data)[bp.BypassResponse])
				rules, ok := v.([]interface{})
				if !ok {
					return fmt.Errorf("tos_bucket_notification %s rules is not slice ", d.Id())
				}
				for index, v := range rules {
					ruleMap, ok := v.(map[string]interface{})
					if !ok {
						return fmt.Errorf("tos_bucket_notification %s rule is not map ", d.Id())
					}
					if ids[1] == ruleMap["RuleId"] {
						rules = append(rules[:index], rules[index+1:]...)
						break
					}
				}
				(*call.SdkParam)[bp.BypassParam].(map[string]interface{})["Rules"] = rules
				return nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				param := (*call.SdkParam)[bp.BypassParam].(map[string]interface{})
				resp, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					HttpMethod:  bp.PUT,
					ContentType: bp.ApplicationJSON,
					Domain:      (*call.SdkParam)["BucketName"].(string),
					UrlParam: map[string]string{
						"notification_v2": "",
					},
				}, &param)
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
							return resource.NonRetryableError(fmt.Errorf("error on reading tos bucket realtime log on delete %q, %w", s.ReadResourceId(d.Id()), callErr))
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

func (s *VestackTosBucketNotificationService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackTosBucketNotificationService) createOrUpdateNotification(resourceData *schema.ResourceData, resource *schema.Resource, isUpdate bool) bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "PutBucketNotificationV2",
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
				//"version": {
				//	ConvertType: bp.ConvertDefault,
				//	TargetField: "Version",
				//	ForceGet:    isUpdate,
				//},
				"rules": {
					ConvertType: bp.ConvertJsonObjectArray,
					TargetField: "Rules",
					//ForceGet: true,
					NextLevelConvert: map[string]bp.RequestConvert{
						"rule_id": {
							ConvertType: bp.ConvertDefault,
							TargetField: "RuleId",
							ForceGet:    isUpdate,
						},
						"events": {
							ConvertType: bp.ConvertJsonArray,
							TargetField: "Events",
							ForceGet:    isUpdate,
						},
						"destination": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "Destination",
							ForceGet:    isUpdate,
							NextLevelConvert: map[string]bp.RequestConvert{
								"ve_faas": {
									ConvertType: bp.ConvertJsonObjectArray,
									TargetField: "VeFaaS",
									ForceGet:    isUpdate,
									NextLevelConvert: map[string]bp.RequestConvert{
										"function_id": {
											ConvertType: bp.ConvertDefault,
											TargetField: "FunctionId",
										},
									},
								},
							},
						},
						"filter": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "Filter",
							ForceGet:    isUpdate,
							NextLevelConvert: map[string]bp.RequestConvert{
								"tos_key": {
									ConvertType: bp.ConvertJsonObject,
									TargetField: "TOSKey",
									ForceGet:    isUpdate,
									NextLevelConvert: map[string]bp.RequestConvert{
										"filter_rules": {
											ConvertType: bp.ConvertJsonObjectArray,
											TargetField: "FilterRules",
											ForceGet:    isUpdate,
											NextLevelConvert: map[string]bp.RequestConvert{
												"name": {
													ConvertType: bp.ConvertDefault,
													TargetField: "Name",
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
							},
						},
					},
				},
			},
			//BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
			//	id := d.Get("rules.0.rule_id")
			//	(*call.SdkParam)["RuleId"] = id.(string)
			//
			//	var sourceParam map[string]interface{}
			//	sourceParam, err := bp.SortAndStartTransJson((*call.SdkParam)[bp.BypassParam].(map[string]interface{}))
			//	if err != nil {
			//		return false, err
			//	}
			//	(*call.SdkParam)[bp.BypassParam] = sourceParam
			//
			//	return true, nil
			//},
			AfterLocked: s.beforePutBucketNotification(isUpdate),
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				//创建 Notification
				param := (*call.SdkParam)[bp.BypassParam].(map[string]interface{})
				resp, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					HttpMethod:  bp.PUT,
					ContentType: bp.ApplicationJSON,
					Domain:      (*call.SdkParam)[bp.BypassDomain].(string),
					UrlParam: map[string]string{
						"notification_v2": "",
					},
				}, &param)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId((*call.SdkParam)[bp.BypassDomain].(string) + ":" + d.Get("rules.0.rule_id").(string))
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("bucket_name").(string)
			},
		},
	}

	return callback
}

func (s *VestackTosBucketNotificationService) beforePutBucketNotification(isUpdate bool) bp.CallFunc {

	return func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) error {
		action := "GetBucketNotificationV2"
		logger.Debug(logger.ReqFormat, action, d.Get("bucket_name"))
		data, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
			HttpMethod: bp.GET,
			Domain:     (*call.SdkParam)[bp.BypassDomain].(string),
			UrlParam: map[string]string{
				"notification_v2": "",
			},
		}, nil)
		logger.Debug(logger.RespFormat, action, data, err)
		return s.beforeTosPutNotification(d, call, data, err, isUpdate)
	}
}

func (s *VestackTosBucketNotificationService) beforeTosPutNotification(d *schema.ResourceData, call bp.SdkCall, data *map[string]interface{}, err error, isUpdate bool) error {
	if err != nil {
		return err
	}
	id := d.Get("rules.0.rule_id")

	var sourceAclParam map[string]interface{}
	sourceAclParam, err = bp.SortAndStartTransJson((*call.SdkParam)[bp.BypassParam].(map[string]interface{}))
	if err != nil {
		return err
	}
	v, _ := bp.ObtainSdkValue("Rules", (*data)[bp.BypassResponse])
	rules, ok := v.([]interface{})
	if !ok {
		return fmt.Errorf("tos_bucket_notification %s rules is not slice ", id)
	}
	if len(rules) == 0 && isUpdate {
		return fmt.Errorf("tos_bucket_notification %s rules is empty", id)
	}
	// 根据 rule id 查找对应 rule
	rule := make(map[string]interface{})
	for index, v := range rules {
		ruleMap, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("tos_bucket_notification %s rule is not map ", id)
		}
		if id == ruleMap["RuleId"] {
			rule = ruleMap
			rules = append(rules[:index], rules[index+1:]...)
			break
		}
	}
	if len(rule) == 0 && isUpdate {
		return fmt.Errorf("tos_bucket_notification %s rule not exist ", id)
	}
	if len(rule) > 0 && !isUpdate {
		return fmt.Errorf("tos_bucket_notification %s rule is existed ", id)
	}
	// merge rule
	rulesParam, _ := bp.ObtainSdkValue("Rules", sourceAclParam)
	if rulesParam != nil {
		_, ok := rulesParam.([]interface{})
		if !ok {
			return fmt.Errorf("tos_bucket_notification %s rules is not slice ", id)
		}
		rulesParam = append(rulesParam.([]interface{}), rules...)
	}
	sourceAclParam["Rules"] = rulesParam

	(*call.SdkParam)[bp.BypassParam] = sourceAclParam
	return nil
}

func (s *VestackTosBucketNotificationService) ReadResourceId(id string) string {
	return id
}
