package iam_oauth_provider

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamOAuthProviderService struct {
	Client *bp.SdkClient
}

func NewIamOAuthProviderService(c *bp.SdkClient) *VestackIamOAuthProviderService {
	return &VestackIamOAuthProviderService{
		Client: c,
	}
}

func (s *VestackIamOAuthProviderService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamOAuthProviderService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{
		{
			Call: bp.SdkCall{
				Action:      "CreateOAuthProvider",
				ConvertMode: bp.RequestConvertIgnore,
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					param, err := bp.ResourceDateToRequest(d, resource, false, s.createRequestConvert(), bp.RequestConvertInConvert, bp.ContentTypeDefault)
					if err != nil {
						return nil, err
					}
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), &param)
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					if v, ok := d.Get("oauth_provider_name").(string); ok && v != "" {
						d.SetId(v)
					}
					return nil
				},
			},
		},
	}
}

func (s *VestackIamOAuthProviderService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{
		{
			Call: bp.SdkCall{
				Action:      "UpdateOAuthProvider",
				ConvertMode: bp.RequestConvertIgnore,
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					param, err := bp.ResourceDateToRequest(d, resource, true, s.createRequestConvert(), bp.RequestConvertInConvert, bp.ContentTypeDefault)
					if err != nil {
						return nil, err
					}
					param["OAuthProviderName"] = d.Id()
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), &param)
				},
			},
		},
	}
}

func (s *VestackIamOAuthProviderService) createRequestConvert() map[string]bp.RequestConvert {
	return map[string]bp.RequestConvert{
		"oauth_provider_name": {TargetField: "OAuthProviderName"},
		"sso_type":            {TargetField: "SSOType"},
		"status":              {TargetField: "Status"},
		"description":         {TargetField: "Description"},
		"client_id":           {TargetField: "ClientId"},
		"client_secret":       {TargetField: "ClientSecret"},
		"user_info_url":       {TargetField: "UserInfoURL"},
		"token_url":           {TargetField: "TokenURL"},
		"authorize_url":       {TargetField: "AuthorizeURL"},
		"authorize_template":  {TargetField: "AuthorizeTemplate"},
		"scope":               {TargetField: "Scope"},
		"identity_map_type":   {TargetField: "IdentityMapType"},
		"idp_identity_key":    {TargetField: "IdpIdentityKey"},
	}
}

func (s *VestackIamOAuthProviderService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{
		{
			Call: bp.SdkCall{
				Action:      "DeleteOAuthProvider",
				ConvertMode: bp.RequestConvertIgnore,
				SdkParam: &map[string]interface{}{
					"OAuthProviderName": data.Id(),
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
			},
		},
	}
}

func (s *VestackIamOAuthProviderService) ReadResource(d *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	action := "GetOAuthProvider"
	if id == "" {
		id = d.Id()
	}
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &map[string]interface{}{
		"OAuthProviderName": id,
	})
	if err != nil {
		return nil, err
	}
	result, err := bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	if resMap, ok := result.(map[string]interface{}); ok {
		return resMap, nil
	}
	return nil, fmt.Errorf("result is not map[string]interface{}")
}

func (s *VestackIamOAuthProviderService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	action := "GetOAuthProvider"
	logger.Debug(logger.ReqFormat, action, m)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &m)
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

func (s *VestackIamOAuthProviderService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ResponseConverts: map[string]bp.ResponseConvert{
			"OAuthProviderName": {TargetField: "oauth_provider_name"},
			"ProviderId":        {TargetField: "provider_id"},
			"SSOType":           {TargetField: "sso_type"},
			"Status":            {TargetField: "status"},
			"Description":       {TargetField: "description"},
			"ClientId":          {TargetField: "client_id"},
			"ClientSecret":      {TargetField: "client_secret"},
			"UserInfoURL":       {TargetField: "user_info_url"},
			"TokenURL":          {TargetField: "token_url"},
			"AuthorizeURL":      {TargetField: "authorize_url"},
			"AuthorizeTemplate": {TargetField: "authorize_template"},
			"Scope":             {TargetField: "scope"},
			"IdentityMapType":   {TargetField: "identity_map_type"},
			"IdpIdentityKey":    {TargetField: "idp_identity_key"},
			"Trn":               {TargetField: "trn"},
			"CreateDate":        {TargetField: "create_date"},
			"UpdateDate":        {TargetField: "update_date"},
		},
		CollectField: "providers",
		RequestConverts: map[string]bp.RequestConvert{
			"oauth_provider_name": {TargetField: "OAuthProviderName"},
		},
	}
}

func (s *VestackIamOAuthProviderService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}
func (s *VestackIamOAuthProviderService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	return nil
}
func (s *VestackIamOAuthProviderService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	method := bp.GET
	if actionName == "CreateOAuthProvider" || actionName == "UpdateOAuthProvider" {
		method = bp.POST
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
