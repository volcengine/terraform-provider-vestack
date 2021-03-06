package acl

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ve "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackAclService struct {
	Client     *ve.SdkClient
	Dispatcher *ve.Dispatcher
}

func NewAclService(c *ve.SdkClient) *VestackAclService {
	return &VestackAclService{
		Client:     c,
		Dispatcher: &ve.Dispatcher{},
	}
}

func (s *VestackAclService) GetClient() *ve.SdkClient {
	return s.Client
}

func (s *VestackAclService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return ve.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		clb := s.Client.ClbClient
		action := "DescribeAcls"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = clb.DescribeAclsCommon(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = clb.DescribeAclsCommon(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = ve.ObtainSdkValue("Result.Acls", *resp)
		if err != nil {
			return data, err
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Acls is not Slice")
		}
		return data, err
	})
}

func (s *VestackAclService) ReadResource(resourceData *schema.ResourceData, aclId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if aclId == "" {
		aclId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"AclIds.1": aclId,
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
		return data, fmt.Errorf("acl %s not exist ", aclId)
	}

	//????????????
	var (
		resp *map[string]interface{}
	)
	action := "DescribeAclAttributes"
	condition := make(map[string]interface{})
	condition["AclId"] = aclId
	clb := s.Client.ClbClient
	logger.Debug(logger.ReqFormat, action, condition)
	resp, err = clb.DescribeAclAttributesCommon(&condition)
	entries, _ := ve.ObtainSdkValue("Result.AclEntries", *resp)
	logger.Debug(logger.ReqFormat, action, condition, entries)
	logger.Debug(logger.ReqFormat, action, condition, data)
	if entries != nil {
		data["AclEntries"] = entries
	}
	return data, err
}

func (s *VestackAclService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh:    nil, // no need to refresh resources
	}

}

func (VestackAclService) WithResourceResponseHandlers(acl map[string]interface{}) []ve.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]ve.ResponseConvert, error) {
		return acl, nil, nil
	}
	return []ve.ResourceResponseHandler{handler}

}

func (s *VestackAclService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []ve.Callback {
	var callbacks []ve.Callback

	callback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "CreateAcl",
			ConvertMode: ve.RequestConvertAll,
			Convert: map[string]ve.RequestConvert{
				"acl_entries": {
					Ignore: true,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.ClbClient.CreateAclCommon(call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *ve.SdkClient, resp *map[string]interface{}, call ve.SdkCall) error {
				id, _ := ve.ObtainSdkValue("Result.AclId", *resp)
				d.SetId(id.(string))
				return nil
			},
		},
	}

	callbacks = append(callbacks, callback)
	//????????????
	entryCallback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "AddAclEntries",
			ConvertMode: ve.RequestConvertInConvert,
			Convert: map[string]ve.RequestConvert{
				"acl_entries": {
					ConvertType: ve.ConvertListN,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (bool, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				if len(*call.SdkParam) > 0 {
					(*call.SdkParam)["AclId"] = d.Id()
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.ClbClient.AddAclEntriesCommon(call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, entryCallback)
	return callbacks

}

func (s *VestackAclService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []ve.Callback {
	var callbacks []ve.Callback

	callback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "ModifyAclAttributes",
			ConvertMode: ve.RequestConvertAll,
			Convert: map[string]ve.RequestConvert{
				"acl_entries": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (bool, error) {
				(*call.SdkParam)["AclId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.ClbClient.ModifyAclAttributesCommon(call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, callback)

	//????????????
	add, remove, _, _ := ve.GetSetDifference("acl_entries", resourceData, ve.ClbAclEntryHash, false)

	entryRemoveCallback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "RemoveAclEntries",
			ConvertMode: ve.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (bool, error) {
				if len(remove.List()) > 0 {
					(*call.SdkParam)["AclId"] = d.Id()
					for index, entry := range remove.List() {
						(*call.SdkParam)["Entries."+strconv.Itoa(index+1)] = entry.(map[string]interface{})["entry"].(string)
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.ClbClient.RemoveAclEntriesCommon(call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *ve.SdkClient, resp *map[string]interface{}, call ve.SdkCall) error {
				//???????????????????????? ?????????????????????
				time.Sleep(time.Duration(5) * time.Second)
				return nil
			},
		},
	}
	callbacks = append(callbacks, entryRemoveCallback)

	entryAddCallback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "AddAclEntries",
			ConvertMode: ve.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (bool, error) {
				if len(add.List()) > 0 {
					(*call.SdkParam)["AclId"] = d.Id()
					for index, entry := range add.List() {
						(*call.SdkParam)["AclEntries."+strconv.Itoa(index+1)+"."+"Entry"] = entry.(map[string]interface{})["entry"].(string)
						(*call.SdkParam)["AclEntries."+strconv.Itoa(index+1)+"."+"Description"] = entry.(map[string]interface{})["description"].(string)
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.ClbClient.AddAclEntriesCommon(call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, entryAddCallback)

	return callbacks
}

func (s *VestackAclService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []ve.Callback {
	callback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "DeleteAcl",
			ConvertMode: ve.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"AclId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.ClbClient.DeleteAclCommon(call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall, baseErr error) error {
				//?????????????????????
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if ve.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading acl on delete %q, %w", d.Id(), callErr))
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
	return []ve.Callback{callback}
}

func (s *VestackAclService) DatasourceResources(*schema.ResourceData, *schema.Resource) ve.DataSourceInfo {
	return ve.DataSourceInfo{
		RequestConverts: map[string]ve.RequestConvert{
			"ids": {
				TargetField: "AclIds",
				ConvertType: ve.ConvertWithN,
			},
		},
		NameField:    "AclName",
		IdField:      "AclId",
		CollectField: "acls",
		ResponseConverts: map[string]ve.ResponseConvert{
			"AclId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *VestackAclService) ReadResourceId(id string) string {
	return id
}
