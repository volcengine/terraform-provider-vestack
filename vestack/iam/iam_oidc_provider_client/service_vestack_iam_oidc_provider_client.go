package iam_oidc_provider_client

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

type VestackIamOidcProviderClientService struct {
	Client *bp.SdkClient
}

func NewIamOidcProviderClientService(c *bp.SdkClient) *VestackIamOidcProviderClientService {
	return &VestackIamOidcProviderClientService{
		Client: c,
	}
}

func (s *VestackIamOidcProviderClientService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamOidcProviderClientService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{
		{
			Call: bp.SdkCall{
				Action:      "AddClientIDToOIDCProvider",
				ConvertMode: bp.RequestConvertIgnore,
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					param, err := bp.ResourceDateToRequest(d, resource, false, s.createRequestConvert(), bp.RequestConvertInConvert, bp.ContentTypeDefault)
					if err != nil {
						return nil, err
					}
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), &param)
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					pName, _ := d.Get("oidc_provider_name").(string)
					cId, _ := d.Get("client_id").(string)
					if pName != "" && cId != "" {
						d.SetId(fmt.Sprintf("%s:%s", pName, cId))
					}
					return nil
				},
			},
		},
	}
}

func (s *VestackIamOidcProviderClientService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{
		{
			Call: bp.SdkCall{
				Action:      "RemoveClientIDFromOIDCProvider",
				ConvertMode: bp.RequestConvertIgnore,
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					parts := strings.Split(d.Id(), ":")
					if len(parts) != 2 {
						return nil, fmt.Errorf("invalid id format")
					}
					param := map[string]interface{}{
						"OIDCProviderName": parts[0],
						"ClientID":         parts[1],
					}
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), &param)
				},
			},
		},
	}
}

func (s *VestackIamOidcProviderClientService) createRequestConvert() map[string]bp.RequestConvert {
	return map[string]bp.RequestConvert{
		"oidc_provider_name": {TargetField: "OIDCProviderName"},
		"client_id":          {TargetField: "ClientID"},
	}
}

func (s *VestackIamOidcProviderClientService) ReadResource(d *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	if id == "" {
		id = d.Id()
	}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format")
	}
	providerName := parts[0]
	clientId := parts[1]

	action := "GetOIDCProvider"
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &map[string]interface{}{
		"OIDCProviderName": providerName,
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return nil, nil
		}
		return nil, err
	}

	clientIds, err := bp.ObtainSdkValue("Result.ClientIDs", *resp)
	if err != nil {
		return nil, err
	}
	if clientIds == nil {
		return nil, nil
	}

	list, ok := clientIds.([]interface{})
	if !ok {
		return nil, fmt.Errorf("ClientIDs is not a list")
	}

	for _, v := range list {
		if vStr, ok := v.(string); ok && vStr == clientId {
			return map[string]interface{}{
				"oidc_provider_name": providerName,
				"client_id":          clientId,
				"id":                 id,
			}, nil
		}
	}

	return nil, nil
}

func (s *VestackIamOidcProviderClientService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return nil, nil
}

func (s *VestackIamOidcProviderClientService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackIamOidcProviderClientService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *VestackIamOidcProviderClientService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackIamOidcProviderClientService) ReadResourceId(id string) string {
	return id
}

func (s *VestackIamOidcProviderClientService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	return nil
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	method := bp.POST
	if actionName == "GetOIDCProvider" {
		method = bp.GET
	}
	return bp.UniversalInfo{
		ServiceName: "iam",
		Action:      actionName,
		Version:     "2018-01-01",
		HttpMethod:  method,
		ContentType: bp.Default,
		RegionType:  bp.Global,
	}
}
