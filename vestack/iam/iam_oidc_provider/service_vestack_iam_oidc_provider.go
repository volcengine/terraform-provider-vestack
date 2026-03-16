package iam_oidc_provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamOidcProviderService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewIamOidcProviderService(c *bp.SdkClient) *VestackIamOidcProviderService {
	return &VestackIamOidcProviderService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackIamOidcProviderService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamOidcProviderService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageOffsetQuery(m, "Limit", "Offset", 100, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListOIDCProviders"

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
		results, err = bp.ObtainSdkValue("Result.OIDCProviders", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.OIDCProviders is not Slice")
		}
		for _, ele := range data {
			provider, ok := ele.(map[string]interface{})
			if !ok {
				continue
			}
			query := map[string]interface{}{
				"OIDCProviderName": provider["ProviderName"],
			}
			action = "GetOIDCProvider"
			logger.Debug(logger.ReqFormat, action, query)
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &query)
			if err != nil {
				return data, err
			}
			logger.Debug(logger.RespFormat, action, query, *resp)

			issuerUrl, err := bp.ObtainSdkValue("Result.IssuerURL", *resp)
			if err != nil {
				return data, err
			}
			provider["IssuerURL"] = issuerUrl

			issuanceLimitTime, err := bp.ObtainSdkValue("Result.IssuanceLimitTime", *resp)
			if err != nil {
				return data, err
			}
			provider["IssuanceLimitTime"] = issuanceLimitTime

			clientIds, err := bp.ObtainSdkValue("Result.ClientIDs", *resp)
			if err != nil {
				return data, err
			}
			provider["ClientIDs"] = clientIds

			thumbprints, err := bp.ObtainSdkValue("Result.Thumbprints", *resp)
			if err != nil {
				return data, err
			}
			provider["Thumbprints"] = thumbprints
		}
		return data, err
	})
}

func (s *VestackIamOidcProviderService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		temp    map[string]interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if temp, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
		if name, ok := temp["ProviderName"].(string); ok && name == id {
			data = temp
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("iam_oidc_provider %s not exist ", id)
	}
	return data, err
}

func (s *VestackIamOidcProviderService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (VestackIamOidcProviderService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"OIDCProviderName": {
				TargetField: "oidc_provider_name",
			},
			"IssuerURL": {
				TargetField: "issuer_url",
			},
			"ClientIDs": {
				TargetField: "client_ids",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackIamOidcProviderService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateOIDCProvider",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"oidc_provider_name": {
					TargetField: "OIDCProviderName",
				},
				"issuer_url": {
					TargetField: "IssuerURL",
				},
				"client_ids": {
					TargetField: "ClientIDs",
					ConvertType: bp.ConvertWithN,
				},
				"thumbprints": {
					TargetField: "Thumbprints",
					ConvertType: bp.ConvertWithN,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, err := bp.ObtainSdkValue("Result.OIDCProviderName", *resp)
				if err != nil {
					return err
				}
				if s, ok := id.(string); ok && s != "" {
					d.SetId(s)
				} else {
					return errors.New("Result.OIDCProviderName is not string")
				}
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackIamOidcProviderService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateOIDCProvider",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"oidc_provider_name": {
					TargetField: "OIDCProviderName",
					ForceGet:    true,
				},
				"description": {
					TargetField: "Description",
				},
				"issuance_limit_time": {
					TargetField: "IssuanceLimitTime",
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
		},
	}
	callbacks = append(callbacks, callback)

	// client_ids
	if resourceData.HasChange("client_ids") {
		addClientIds, removeClientIds, _, _ := bp.GetSetDifference("client_ids", resourceData, schema.HashString, false)
		for _, element := range addClientIds.List() {
			callbacks = append(callbacks, s.updateOidcProviderCallback(resourceData, "AddClientIDToOIDCProvider", "ClientID", element))
		}
		for _, element := range removeClientIds.List() {
			callbacks = append(callbacks, s.updateOidcProviderCallback(resourceData, "RemoveClientIDFromOIDCProvider", "ClientID", element))
		}
	}

	// thumbprints
	if resourceData.HasChange("thumbprints") {
		addThumbprints, removeThumbprints, _, _ := bp.GetSetDifference("thumbprints", resourceData, schema.HashString, false)
		for _, element := range addThumbprints.List() {
			callbacks = append(callbacks, s.updateOidcProviderCallback(resourceData, "AddThumbprintToOIDCProvider", "Thumbprint", element))
		}
		for _, element := range removeThumbprints.List() {
			callbacks = append(callbacks, s.updateOidcProviderCallback(resourceData, "RemoveThumbprintFromOIDCProvider", "Thumbprint", element))
		}
	}

	return callbacks
}

func (s *VestackIamOidcProviderService) updateOidcProviderCallback(resourceData *schema.ResourceData, action, field string, element interface{}) bp.Callback {
	return bp.Callback{
		Call: bp.SdkCall{
			Action:      action,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["OIDCProviderName"] = d.Id()
				s, ok := element.(string)
				if !ok {
					return false, errors.New("element is not string")
				}
				(*call.SdkParam)[field] = s
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, resp)
				return resp, err
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Id()
			},
		},
	}
}

func (s *VestackIamOidcProviderService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteOIDCProvider",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"OIDCProviderName": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(5*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading iam oidc provider on delete %q, %w", d.Id(), callErr))
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
	return []bp.Callback{callback}
}

func (s *VestackIamOidcProviderService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		NameField:    "OIDCProviderName",
		CollectField: "oidc_providers",
		ResponseConverts: map[string]bp.ResponseConvert{
			"IssuerURL": {
				TargetField: "issuer_url",
			},
			"ClientIDs": {
				TargetField: "client_ids",
			},
		},
	}
}

func (s *VestackIamOidcProviderService) ReadResourceId(id string) string {
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
