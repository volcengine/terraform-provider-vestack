package iam_access_key_last_used

import (
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamAccessKeyLastUsedService struct {
	Client *bp.SdkClient
}

func NewIamAccessKeyLastUsedService(c *bp.SdkClient) *VestackIamAccessKeyLastUsedService {
	return &VestackIamAccessKeyLastUsedService{
		Client: c,
	}
}

func (s *VestackIamAccessKeyLastUsedService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamAccessKeyLastUsedService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp *map[string]interface{}
	)
	action := "GetAccessKeyLastUsed"
	logger.Debug(logger.ReqFormat, action, m)
	if m == nil {
		return data, errors.New("missing params")
	}

	resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &m)
	if err != nil {
		return data, err
	}

	result, err := bp.ObtainSdkValue("Result.AccessKeyLastUsed", *resp)
	if err != nil {
		return data, err
	}
	if result == nil {
		return []interface{}{}, nil
	}

	return []interface{}{result}, nil
}

func (s *VestackIamAccessKeyLastUsedService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"access_key_id": {TargetField: "AccessKeyId"},
			"user_name":     {TargetField: "UserName"},
		},
		ResponseConverts: map[string]bp.ResponseConvert{
			"Region": {
				TargetField: "region",
			},
			"Service": {
				TargetField: "service",
			},
			"RequestTime": {
				TargetField: "request_time",
			},
		},
		CollectField: "access_key_last_useds",
	}
}

func (s *VestackIamAccessKeyLastUsedService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	return nil, nil
}

func (s *VestackIamAccessKeyLastUsedService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *VestackIamAccessKeyLastUsedService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	return nil
}

func (s *VestackIamAccessKeyLastUsedService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}

func (s *VestackIamAccessKeyLastUsedService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}

func (s *VestackIamAccessKeyLastUsedService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return nil
}

func (s *VestackIamAccessKeyLastUsedService) ReadResourceId(id string) string {
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
