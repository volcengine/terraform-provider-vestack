package iam_identity_provider

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamIdentityProviderService struct {
	Client *bp.SdkClient
}

func NewIamIdentityProviderService(c *bp.SdkClient) *VestackIamIdentityProviderService {
	return &VestackIamIdentityProviderService{
		Client: c,
	}
}

func (s *VestackIamIdentityProviderService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamIdentityProviderService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return bp.WithPageNumberQuery(m, "Limit", "Offset", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		universalClient := s.Client.UniversalClient
		action := "ListIdentityProviders"
		logger.Debug(logger.ReqFormat, action, condition)
		resp, err := universalClient.DoCall(getUniversalInfo(action), &condition)
		if err != nil {
			return nil, err
		}
		logger.Debug(logger.RespFormat, action, resp)
		results, err := bp.ObtainSdkValue("Result.IdentityProviders", *resp)
		if err != nil {
			return nil, err
		}
		if results == nil {
			return []interface{}{}, nil
		}
		if data, ok := results.([]interface{}); ok {
			return data, nil
		}
		return nil, fmt.Errorf("Result.IdentityProviders is not Slice")
	})
}

func (s *VestackIamIdentityProviderService) DatasourceResources(d *schema.ResourceData, r *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{},
		CollectField:    "providers",
		ResponseConverts: map[string]bp.ResponseConvert{
			"SSOType": {
				TargetField: "sso_type",
			},
		},
	}
}

func (s *VestackIamIdentityProviderService) CreateResource(d *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackIamIdentityProviderService) ReadResource(d *schema.ResourceData, id string) (map[string]interface{}, error) {
	return nil, nil
}

func (s *VestackIamIdentityProviderService) ModifyResource(d *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackIamIdentityProviderService) RemoveResource(d *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackIamIdentityProviderService) RefreshResourceState(d *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *VestackIamIdentityProviderService) ReadResourceId(id string) string {
	return id
}

func (s *VestackIamIdentityProviderService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	return nil
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
