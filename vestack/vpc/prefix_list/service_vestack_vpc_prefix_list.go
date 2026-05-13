package prefix_list

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackVpcPrefixListService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewVpcPrefixListService(c *bp.SdkClient) *VestackVpcPrefixListService {
	return &VestackVpcPrefixListService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackVpcPrefixListService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackVpcPrefixListService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribePrefixLists"
		bytes, _ := json.Marshal(condition)
		logger.Debug(logger.ReqFormat, action, string(bytes))
		if condition == nil {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, err
			}
		}
		respBytes, _ := json.Marshal(resp)
		logger.Debug(logger.RespFormat, action, condition, string(respBytes))
		results, err = bp.ObtainSdkValue("Result.PrefixLists", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Items is not Slice")
		}
		for index, ele := range data {
			prefixList := ele.(map[string]interface{})
			prefixId := prefixList["PrefixListId"].(string)
			query := map[string]interface{}{
				"PrefixListId": prefixId,
			}
			action = "DescribePrefixListAssociations"
			logger.Debug(logger.ReqFormat, action, string(bytes))
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &query)
			if err != nil {
				return data, err
			}
			logger.Debug(logger.RespFormat, action, condition, string(respBytes))
			prefixListAssociations, err := bp.ObtainSdkValue("Result.PrefixListAssociations", *resp)
			if err != nil {
				return data, err
			}
			data[index].(map[string]interface{})["PrefixListAssociations"] = prefixListAssociations
			action = "DescribePrefixListEntries"
			logger.Debug(logger.ReqFormat, action, string(bytes))
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &query)
			if err != nil {
				return data, err
			}
			logger.Debug(logger.RespFormat, action, condition, string(respBytes))
			prefixListEntries, err := bp.ObtainSdkValue("Result.PrefixListEntries", *resp)
			if err != nil {
				return data, err
			}
			data[index].(map[string]interface{})["PrefixListEntries"] = prefixListEntries
		}
		return data, err
	})
}

func (s *VestackVpcPrefixListService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"PrefixListIds.1": id,
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
		return data, fmt.Errorf("vpc_prefix_list %s not exist ", id)
	}
	return data, err
}

func (s *VestackVpcPrefixListService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				d          map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "Error")
			if err = resource.Retry(20*time.Minute, func() *resource.RetryError {
				d, err = s.ReadResource(resourceData, id)
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
			status, err = bp.ObtainSdkValue("Status", d)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("vpc_prefix_list status error, status: %s", status.(string))
				}
			}
			return d, status.(string), err
		},
	}
}

func (s *VestackVpcPrefixListService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreatePrefixList",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"prefix_list_entries": {
					TargetField: "PrefixListEntries",
					ConvertType: bp.ConvertListN,
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertListN,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.PrefixListId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (VestackVpcPrefixListService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackVpcPrefixListService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyPrefixList",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"prefix_list_name": {
					ConvertType: bp.ConvertDefault,
				},
				"description": {
					ConvertType: bp.ConvertDefault,
				},
				"max_entries": {
					ConvertType: bp.ConvertDefault,
				},
				"prefix_list_entries": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["PrefixListId"] = d.Id()
				addEntries, removeEntries, _, _ := bp.GetSetDifference("prefix_list_entries", d, entriesHash, false)
				if addEntries != nil && addEntries.Len() > 0 {
					for index, entry := range addEntries.List() {
						if cidr, ok := entry.(map[string]interface{})["cidr"]; ok {
							(*call.SdkParam)["AddPrefixListEntries."+strconv.Itoa(index+1)+".Cidr"] = cidr
						}
						if description, ok := entry.(map[string]interface{})["description"]; ok {
							(*call.SdkParam)["AddPrefixListEntries."+strconv.Itoa(index+1)+".Description"] = description
						}
					}
				}
				if removeEntries != nil && removeEntries.Len() > 0 {
					for index, entry := range removeEntries.List() {
						if cidr, ok := entry.(map[string]interface{})["cidr"]; ok {
							(*call.SdkParam)["RemovePrefixListEntries."+strconv.Itoa(index+1)+".Cidr"] = cidr
						}
					}
				}
				delete(*call.SdkParam, "Tags")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, callback)
	return callbacks
}

func (s *VestackVpcPrefixListService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeletePrefixList",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"PrefixListId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackVpcPrefixListService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "PrefixListIds",
				ConvertType: bp.ConvertWithN,
			},
			"tag_filters": {
				TargetField: "TagFilters",
				ConvertType: bp.ConvertListN,
				NextLevelConvert: map[string]bp.RequestConvert{
					"values": {
						ConvertType: bp.ConvertWithN,
					},
				},
			},
		},
		NameField:    "PrefixListName",
		IdField:      "PrefixListId",
		CollectField: "prefix_lists",
		ResponseConverts: map[string]bp.ResponseConvert{
			"PrefixListId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *VestackVpcPrefixListService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vpc",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
