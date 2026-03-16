package iam_oidc_provider_thumbprint

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

type VestackIamOidcProviderThumbprintService struct {
	Client *bp.SdkClient
}

func NewIamOidcProviderThumbprintService(c *bp.SdkClient) *VestackIamOidcProviderThumbprintService {
	return &VestackIamOidcProviderThumbprintService{
		Client: c,
	}
}

func (s *VestackIamOidcProviderThumbprintService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamOidcProviderThumbprintService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{
		{
			Call: bp.SdkCall{
				Action:      "AddThumbprintToOIDCProvider",
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
					thumbprint, _ := d.Get("thumbprint").(string)
					if pName != "" && thumbprint != "" {
						d.SetId(fmt.Sprintf("%s:%s", pName, thumbprint))
					}
					return nil
				},
			},
		},
	}
}

func (s *VestackIamOidcProviderThumbprintService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{
		{
			Call: bp.SdkCall{
				Action:      "RemoveThumbprintFromOIDCProvider",
				ConvertMode: bp.RequestConvertIgnore,
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					parts := strings.Split(d.Id(), ":")
					if len(parts) != 2 {
						return nil, fmt.Errorf("invalid id format")
					}
					param := map[string]interface{}{
						"OIDCProviderName": parts[0],
						"Thumbprint":       parts[1],
					}
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), &param)
				},
			},
		},
	}
}

func (s *VestackIamOidcProviderThumbprintService) createRequestConvert() map[string]bp.RequestConvert {
	return map[string]bp.RequestConvert{
		"oidc_provider_name": {TargetField: "OIDCProviderName"},
		"thumbprint":         {TargetField: "Thumbprint"},
	}
}

func (s *VestackIamOidcProviderThumbprintService) ReadResource(d *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	if id == "" {
		id = d.Id()
	}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format")
	}
	providerName := parts[0]
	thumbprint := parts[1]

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

	thumbprints, err := bp.ObtainSdkValue("Result.Thumbprints", *resp)
	if err != nil {
		return nil, err
	}
	if thumbprints == nil {
		return nil, nil
	}

	list, ok := thumbprints.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Thumbprints is not a list")
	}

	for _, v := range list {
		if vStr, ok := v.(string); ok && vStr == thumbprint {
			return map[string]interface{}{
				"oidc_provider_name": providerName,
				"thumbprint":         thumbprint,
				"id":                 id,
			}, nil
		}
	}

	return nil, nil
}

func (s *VestackIamOidcProviderThumbprintService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return nil, nil
}

func (s *VestackIamOidcProviderThumbprintService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackIamOidcProviderThumbprintService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *VestackIamOidcProviderThumbprintService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackIamOidcProviderThumbprintService) ReadResourceId(id string) string {
	return id
}

func (s *VestackIamOidcProviderThumbprintService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
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
