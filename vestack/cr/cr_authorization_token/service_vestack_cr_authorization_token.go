package cr_authorization_token

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackCrAuthorizationTokenService struct {
	Client *bp.SdkClient
}

func NewCrAuthorizationTokenService(c *bp.SdkClient) *VestackCrAuthorizationTokenService {
	return &VestackCrAuthorizationTokenService{
		Client: c,
	}
}

func (s *VestackCrAuthorizationTokenService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackCrAuthorizationTokenService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	action := "GetAuthorizationToken"
	logger.Debug(logger.ReqFormat, action, condition)
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

	logger.Debug(logger.RespFormat, action, resp)
	results, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return data, err
	}
	if results == nil {
		return data, fmt.Errorf("GetAuthorizationToken return an empty result")
	}

	token, err := bp.ObtainSdkValue("Result.Token", *resp)
	if err != nil {
		return data, err
	}
	username, err := bp.ObtainSdkValue("Result.Username", *resp)
	if err != nil {
		return data, err
	}
	expireTime, err := bp.ObtainSdkValue("Result.ExpireTime", *resp)
	if err != nil {
		return data, err
	}

	user := map[string]interface{}{
		"Token":      token,
		"Username":   username,
		"ExpireTime": expireTime,
	}

	return []interface{}{user}, err
}

func (s *VestackCrAuthorizationTokenService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	return data, err
}

func (s *VestackCrAuthorizationTokenService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (VestackCrAuthorizationTokenService) WithResourceResponseHandlers(instance map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return instance, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackCrAuthorizationTokenService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackCrAuthorizationTokenService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return s.CreateResource(resourceData, resource)
}

func (s *VestackCrAuthorizationTokenService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackCrAuthorizationTokenService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType:  bp.ContentTypeJson,
		CollectField: "tokens",
	}
}

func (s *VestackCrAuthorizationTokenService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "cr",
		Version:     "2022-05-12",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
