package tos_bucket_inventory

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

type VestackTosBucketInventoryService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewTosBucketInventoryService(c *bp.SdkClient) *VestackTosBucketInventoryService {
	return &VestackTosBucketInventoryService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackTosBucketInventoryService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackTosBucketInventoryService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	tos := s.Client.BypassSvcClient
	var (
		action  string
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	action = "ListBucketInventory"
	logger.Debug(logger.ReqFormat, action, nil)
	resp, err = tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     condition[bp.BypassDomain].(string),
		UrlParam: map[string]string{
			"inventory": "",
		},
	}, nil)
	if err != nil {
		return data, err
	}
	results, err = bp.ObtainSdkValue(bp.BypassResponse+".InventoryConfigurations", *resp)
	if err != nil {
		return data, err
	}

	if results == nil {
		results = []interface{}{}
	}
	if data, ok = results.([]interface{}); !ok {
		return data, errors.New("InventoryConfigurations is not Slice")
	}
	return data, err
}

func (s *VestackTosBucketInventoryService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	tos := s.Client.BypassSvcClient
	var (
		ok bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return data, fmt.Errorf("invalid tos inventory id: %s", id)
	}

	action := "GetBucketInventory"
	logger.Debug(logger.ReqFormat, action, id)
	resp, err := tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     ids[0],
		UrlParam: map[string]string{
			"inventory": "",
			"id":        ids[1],
		},
	}, nil)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp, err)
	if data, ok = (*resp)[bp.BypassResponse].(map[string]interface{}); !ok {
		return data, errors.New("GetBucketInventory Resp is not map ")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("tos_bucket_inventory %s not exist ", id)
	}

	data["BucketName"] = ids[0]
	if destination, ok := data["Destination"].(map[string]interface{}); ok {
		if tosBucketDestination, ok := destination["TOSBucketDestination"].(map[string]interface{}); ok {
			destination["TOSBucketDestination"] = []interface{}{tosBucketDestination}
		}
	}

	return data, err
}

func (s *VestackTosBucketInventoryService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (VestackTosBucketInventoryService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"TOSBucketDestination": {
				TargetField: "tos_bucket_destination",
			},
			"Id": {
				TargetField: "inventory_id",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackTosBucketInventoryService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	//create inventory
	callback := s.createOrUpdateInventory(resourceData, resource, false)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketInventoryService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	//create inventory
	callback := s.createOrUpdateInventory(resourceData, resource, true)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketInventoryService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteBucketInventory",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ids := strings.Split(d.Id(), ":")
				if len(ids) != 2 {
					return false, fmt.Errorf("invalid tos inventory id: %s", d.Id())
				}
				(*call.SdkParam)["BucketName"] = ids[0]
				(*call.SdkParam)["Id"] = ids[1]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					ContentType: bp.ApplicationJSON,
					HttpMethod:  bp.DELETE,
					Domain:      (*call.SdkParam)["BucketName"].(string),
					UrlParam: map[string]string{
						"inventory": "",
						"id":        (*call.SdkParam)["Id"].(string),
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
							return resource.NonRetryableError(fmt.Errorf("error on reading tos inventory on delete %q, %w", s.ReadResourceId(d.Id()), callErr))
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

func (s *VestackTosBucketInventoryService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	bucketName := data.Get("bucket_name")
	inventoryId, ok := data.GetOk("inventory_id")
	return bp.DataSourceInfo{
		ServiceCategory: bp.ServiceBypass,
		RequestConverts: map[string]bp.RequestConvert{
			"bucket_name": {
				ConvertType: bp.ConvertDefault,
				SpecialParam: &bp.SpecialParam{
					Type: bp.DomainParam,
				},
			},
			"inventory_id": {
				Ignore: true,
			},
		},
		NameField:    "Id",
		IdField:      "InventoryId",
		CollectField: "inventory_configurations",
		ResponseConverts: map[string]bp.ResponseConvert{
			"TOSBucketDestination": {
				TargetField: "tos_bucket_destination",
			},
		},
		ExtraData: func(sourceData []interface{}) (extraData []interface{}, err error) {
			for _, v := range sourceData {
				if ok {
					if inventoryId.(string) == v.(map[string]interface{})["Id"].(string) {
						v.(map[string]interface{})["InventoryId"] = bucketName.(string) + ":" + v.(map[string]interface{})["Id"].(string)
						extraData = append(extraData, v)
						break
					} else {
						continue
					}
				} else {
					v.(map[string]interface{})["InventoryId"] = bucketName.(string) + ":" + v.(map[string]interface{})["Id"].(string)
					extraData = append(extraData, v)
				}

			}
			return extraData, err
		},
	}
}

func (s *VestackTosBucketInventoryService) createOrUpdateInventory(resourceData *schema.ResourceData, resource *schema.Resource, isUpdate bool) bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "PutBucketInventory",
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
				"inventory_id": {
					ConvertType: bp.ConvertDefault,
					TargetField: "Id",
					ForceGet:    isUpdate,
				},
				"is_enabled": {
					ConvertType: bp.ConvertDefault,
					TargetField: "IsEnabled",
					ForceGet:    true,
				},
				"included_object_versions": {
					ConvertType: bp.ConvertDefault,
					TargetField: "IncludedObjectVersions",
					ForceGet:    true,
				},
				"schedule": {
					ConvertType: bp.ConvertJsonObject,
					TargetField: "Schedule",
					ForceGet:    true,
				},
				"filter": {
					ConvertType: bp.ConvertJsonObject,
					TargetField: "Filter",
					ForceGet:    true,
				},
				"optional_fields": {
					ConvertType: bp.ConvertJsonObject,
					TargetField: "OptionalFields",
					ForceGet:    true,
					NextLevelConvert: map[string]bp.RequestConvert{
						"field": {
							ConvertType: bp.ConvertJsonArray,
							TargetField: "Field",
						},
					},
				},
				"destination": {
					ConvertType: bp.ConvertJsonObject,
					TargetField: "Destination",
					ForceGet:    true,
					NextLevelConvert: map[string]bp.RequestConvert{
						"tos_bucket_destination": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "TosBucketDestination",
						},
					},
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				id := d.Get("inventory_id")
				(*call.SdkParam)["InventoryId"] = id.(string)

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
				//创建 Inventory
				param := (*call.SdkParam)[bp.BypassParam].(map[string]interface{})
				resp, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					HttpMethod:  bp.PUT,
					ContentType: bp.ApplicationJSON,
					Domain:      (*call.SdkParam)[bp.BypassDomain].(string),
					UrlParam: map[string]string{
						"inventory": "",
						"id":        (*call.SdkParam)["InventoryId"].(string),
					},
				}, &param)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId((*call.SdkParam)[bp.BypassDomain].(string) + ":" + (*call.SdkParam)["InventoryId"].(string))
				return nil
			},
		},
	}

	return callback
}

func (s *VestackTosBucketInventoryService) ReadResourceId(id string) string {
	return id
}
