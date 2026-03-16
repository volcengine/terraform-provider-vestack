package iam_group_user

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamGroupUserService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewIamGroupUserService(c *bp.SdkClient) *VestackIamGroupUserService {
	return &VestackIamGroupUserService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackIamGroupUserService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamGroupUserService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageOffsetQuery(m, "Limit", "Offset", 100, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListGroupsForUser"

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
		results, err = bp.ObtainSdkValue("Result.UserGroups", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.UserGroups is not Slice")
		}
		return data, err
	})
}

func (s *VestackIamGroupUserService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	return data, nil
}

func (s *VestackIamGroupUserService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *VestackIamGroupUserService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (VestackIamGroupUserService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackIamGroupUserService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackIamGroupUserService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackIamGroupUserService) DatasourceResources(d *schema.ResourceData, r *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"user_name": {
				TargetField: "UserName",
			},
			"query": {
				TargetField: "Query",
			},
		},
		NameField:    "UserGroupName",
		IdField:      "UserGroupName",
		CollectField: "user_groups",
		ResponseConverts: map[string]bp.ResponseConvert{
			"UserGroupID": {
				TargetField: "user_group_id",
			},
			"UserGroupName": {
				TargetField: "user_group_name",
			},
			"DisplayName": {
				TargetField: "display_name",
			},
			"Description": {
				TargetField: "description",
			},
			"JoinDate": {
				TargetField: "join_date",
			},
		},
	}
}

func (s *VestackIamGroupUserService) ReadResourceId(id string) string {
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
