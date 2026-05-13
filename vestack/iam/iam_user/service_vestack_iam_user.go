package iam_user

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamUserService struct {
	Client *bp.SdkClient
}

func NewIamUserService(c *bp.SdkClient) *VestackIamUserService {
	return &VestackIamUserService{
		Client: c,
	}
}

func (s *VestackIamUserService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamUserService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	cens, err := bp.WithPageOffsetQuery(m, "Limit", "Offset", 100, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		universalClient := s.Client.UniversalClient
		action := "ListUsers"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = universalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = universalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, err
			}
		}
		logger.Debug(logger.RespFormat, action, resp)
		results, err = bp.ObtainSdkValue("Result.UserMetadata", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.UserMetadata is not Slice")
		}
		data, err = removeSystemTags(data)
		return data, err
	})
	return cens, err
}

func (s *VestackIamUserService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"Query": id,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		m, ok := v.(map[string]interface{})
		if !ok {
			return data, errors.New("result item is not map")
		}
		if name, ok := m["UserName"].(string); ok && name == id {
			data = m
			break
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("user %s not exist ", id)
	}

	return data, err
}

func (s *VestackIamUserService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (VestackIamUserService) WithResourceResponseHandlers(v map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return v, map[string]bp.ResponseConvert{
			"AccountId": {
				TargetField: "account_id",
				Convert: func(i interface{}) interface{} {
					if v, ok := i.(float64); ok {
						return strconv.FormatFloat(v, 'f', 0, 64)
					}
					return ""
				},
			},
			"Id": {
				TargetField: "user_id",
				Convert: func(i interface{}) interface{} {
					if v, ok := i.(float64); ok {
						return strconv.FormatFloat(v, 'f', 0, 64)
					}
					return ""
				},
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackIamUserService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateUser",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertListN,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, resp)
				v, ok := d.Get("user_name").(string)
				if !ok || v == "" {
					return errors.New("user_name is not string or empty")
				}
				d.SetId(v)
				return nil
			},
		},
	}

	return []bp.Callback{callback}
}

func (s *VestackIamUserService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateUser",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"user_name": {
					TargetField: "NewUserName",
					ConvertType: bp.ConvertDefault,
				},
				"display_name": {
					TargetField: "NewDisplayName",
					ConvertType: bp.ConvertDefault,
					Convert:     defaultConvert,
				},
				"mobile_phone": {
					TargetField: "NewMobilePhone",
				},
				"email": {
					TargetField: "NewEmail",
					Convert:     defaultConvert,
				},
				"description": {
					TargetField: "NewDescription",
					ConvertType: bp.ConvertDefault,
					Convert:     defaultConvert,
				},
			},
			RequestIdField: "UserName",
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				if d.HasChange("user_name") {
					v, ok := d.Get("user_name").(string)
					if !ok || v == "" {
						return errors.New("user_name is not string or empty")
					}
					d.SetId(v)
				}
				return nil
			},
		},
	}
	callbacks = append(callbacks, callback)
	setResourceTagsCallbacks := s.setResourceTags(resourceData, "User", callbacks)
	return setResourceTagsCallbacks
}

func (s *VestackIamUserService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:         "DeleteUser",
			ConvertMode:    bp.RequestConvertIgnore,
			RequestIdField: "UserName",
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}

	return []bp.Callback{callback}
}

func (s *VestackIamUserService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"query": {
				TargetField: "Query",
			},
		},
		NameField:    "UserName",
		IdField:      "UserName",
		CollectField: "users",
		ResponseConverts: map[string]bp.ResponseConvert{
			"Id": {
				TargetField: "user_id",
				Convert: func(i interface{}) interface{} {
					if v, ok := i.(float64); ok {
						return strconv.FormatFloat(v, 'f', 0, 64)
					}
					return ""
				},
			},
			"AccountId": {
				TargetField: "account_id",
				Convert: func(i interface{}) interface{} {
					if v, ok := i.(float64); ok {
						return strconv.FormatFloat(v, 'f', 0, 64)
					}
					return ""
				},
			},
			"Tags": {
				TargetField: "tags",
				Convert: func(i interface{}) interface{} {
					if i == nil {
						return nil
					}
					var tags []map[string]interface{}
					if list, ok := i.([]interface{}); ok {
						for _, v := range list {
							if m, ok := v.(map[string]interface{}); ok {
								tag := make(map[string]interface{})
								if key, ok := m["Key"].(string); ok {
									tag["key"] = key
								}
								if value, ok := m["Value"].(string); ok {
									tag["value"] = value
								}
								tags = append(tags, tag)
							}
						}
					}
					return tags
				},
			},
		},
	}
}

func (s *VestackIamUserService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "iam",
		Action:      actionName,
		Version:     "2018-01-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		RegionType:  bp.Global,
	}
}

func (s *VestackIamUserService) setResourceTags(resourceData *schema.ResourceData, resourceType string, callbacks []bp.Callback) []bp.Callback {
	addedTags, removedTags, _, _ := bp.GetSetDifference("tags", resourceData, bp.TagsHash, false)

	removeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UntagResources",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if removedTags != nil && len(removedTags.List()) > 0 {
					(*call.SdkParam)["ResourceNames.1"] = resourceData.Id()
					(*call.SdkParam)["ResourceType"] = resourceType
					for index, tag := range removedTags.List() {
						tm, ok := tag.(map[string]interface{})
						if !ok {
							return false, errors.New("tag item is not map")
						}
						key, ok := tm["key"].(string)
						if !ok {
							return false, errors.New("tag key is not string")
						}
						(*call.SdkParam)["TagKeys."+strconv.Itoa(index+1)] = key
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, removeCallback)

	addCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "TagResources",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if addedTags != nil && len(addedTags.List()) > 0 {
					(*call.SdkParam)["ResourceNames.1"] = resourceData.Id()
					(*call.SdkParam)["ResourceType"] = resourceType
					for index, tag := range addedTags.List() {
						tm, ok := tag.(map[string]interface{})
						if !ok {
							return false, errors.New("tag item is not map")
						}
						key, ok := tm["key"].(string)
						if !ok {
							return false, errors.New("tag key is not string")
						}
						value, ok := tm["value"].(string)
						if !ok {
							return false, errors.New("tag value is not string")
						}
						(*call.SdkParam)["Tags."+strconv.Itoa(index+1)+".Key"] = key
						(*call.SdkParam)["Tags."+strconv.Itoa(index+1)+".Value"] = value
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, addCallback)

	return callbacks
}

func removeSystemTags(data []interface{}) ([]interface{}, error) {
	var (
		ok      bool
		result  map[string]interface{}
		results []interface{}
		tags    []interface{}
	)
	for _, d := range data {
		if result, ok = d.(map[string]interface{}); !ok {
			return results, errors.New("The elements in data are not map ")
		}
		tags, ok = result["Tags"].([]interface{})
		if ok {
			tags = bp.FilterSystemTags(tags)
			result["Tags"] = tags
		}
		results = append(results, result)
	}
	return results, nil
}
