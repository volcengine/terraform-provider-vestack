package iam_caller_identity

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamCallerIdentityService struct {
	Client *bp.SdkClient
}

func NewIamCallerIdentityService(c *bp.SdkClient) *VestackIamCallerIdentityService {
	return &VestackIamCallerIdentityService{
		Client: c,
	}
}

func (s *VestackIamCallerIdentityService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamCallerIdentityService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	action := "GetCallerIdentity"
	logger.Debug(logger.ReqFormat, action, m)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
	if err != nil {
		return data, err
	}

	result, err := bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return data, err
	}
	if result == nil {
		return []interface{}{}, nil
	}

	return []interface{}{result}, nil
}

func (s *VestackIamCallerIdentityService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ResponseConverts: map[string]bp.ResponseConvert{
			"AccountId": {
				TargetField: "account_id",
				Convert: func(i interface{}) interface{} {
					if v, ok := i.(float64); ok {
						return strconv.FormatFloat(v, 'f', 0, 64)
					}
					return i
				},
			},
			"Trn": {
				TargetField: "trn",
			},
			"IdentityType": {
				TargetField: "identity_type",
			},
			"IdentityId": {
				TargetField: "identity_id",
				Convert: func(i interface{}) interface{} {
					if v, ok := i.(float64); ok {
						return strconv.FormatFloat(v, 'f', 0, 64)
					}
					return i
				},
			},
		},
		CollectField: "caller_identities",
	}
}

func (s *VestackIamCallerIdentityService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	return nil, nil
}
func (s *VestackIamCallerIdentityService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}
func (s *VestackIamCallerIdentityService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	return nil
}
func (s *VestackIamCallerIdentityService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}
func (s *VestackIamCallerIdentityService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}
func (s *VestackIamCallerIdentityService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return nil
}
func (s *VestackIamCallerIdentityService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "sts",
		Action:      actionName,
		Version:     "2018-01-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		RegionType:  bp.Global,
	}
}
