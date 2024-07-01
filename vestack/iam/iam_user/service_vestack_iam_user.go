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
		nameSet = make(map[string]bool)
	)
	if _, ok = m["UserNames.1"]; ok {
		i := 1
		for {
			filed := fmt.Sprintf("UserNames.%d", i)
			tmpId, ok := m[filed]
			if !ok {
				break
			}
			nameSet[tmpId.(string)] = true
			i++
			delete(m, filed)
		}
	}
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
		return data, err
	})
	if err != nil || len(nameSet) == 0 {
		return cens, err
	}

	res := make([]interface{}, 0)
	for _, cen := range cens {
		if !nameSet[cen.(map[string]interface{})["UserName"].(string)] {
			continue
		}
		res = append(res, cen)
	}
	return res, nil
}

func (s *VestackIamUserService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"UserNames.1": id,
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
					return strconv.FormatFloat(i.(float64), 'f', 0, 64)
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
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, resp)
				d.SetId(d.Get("user_name").(string))
				return nil
			},
		},
	}

	return []bp.Callback{callback}
}

func (s *VestackIamUserService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
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
					d.SetId(d.Get("user_name").(string))
				}
				return nil
			},
		},
	}
	return []bp.Callback{callback}
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
			"user_names": {
				TargetField: "UserNames",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "UserName",
		IdField:      "UserName",
		CollectField: "users",
		ResponseConverts: map[string]bp.ResponseConvert{
			"AccountId": {
				TargetField: "account_id",
				Convert: func(i interface{}) interface{} {
					return strconv.FormatFloat(i.(float64), 'f', 0, 64)
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
	}
}
