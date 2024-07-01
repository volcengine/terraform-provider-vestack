package ecs_key_pair

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackEcsKeyPairService struct {
	Client *bp.SdkClient
}

func NewEcsKeyPairService(c *bp.SdkClient) *VestackEcsKeyPairService {
	return &VestackEcsKeyPairService{
		Client: c,
	}
}

func (s *VestackEcsKeyPairService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackEcsKeyPairService) ReadResources(condition map[string]interface{}) ([]interface{}, error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithNextTokenQuery(condition, "MaxResults", "NextToken", 20, nil, func(m map[string]interface{}) (data []interface{}, next string, err error) {
		client := s.Client.UniversalClient
		action := "DescribeKeyPairs"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = client.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, next, err
			}
		} else {
			resp, err = client.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, next, err
			}
		}
		logger.Debug(logger.RespFormat, action, condition, *resp)

		results, err = bp.ObtainSdkValue("Result.KeyPairs", *resp)
		if err != nil {
			return data, next, err
		}
		nextToken, err := bp.ObtainSdkValue("Result.NextToken", *resp)
		if err != nil {
			return data, next, err
		}
		next = nextToken.(string)
		if results == nil {
			results = []interface{}{}
		}

		if data, ok = results.([]interface{}); !ok {
			return data, next, errors.New("Result.KeyPairs is not Slice")
		}
		return data, next, err
	})
}

func (s *VestackEcsKeyPairService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"KeyPairIds.1": id,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, fmt.Errorf("Value is not map ")
		}
	}

	if len(data) == 0 {
		return data, fmt.Errorf("Ecs key pair %s not exist ", id)
	}
	return data, nil
}

func (s *VestackEcsKeyPairService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *VestackEcsKeyPairService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackEcsKeyPairService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	action := "CreateKeyPair"
	if _, ok := resourceData.GetOk("public_key"); ok {
		action = "ImportKeyPair"
	}
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      action,
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.KeyPairId", *resp)
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, id)
				d.SetId(id.(string))

				if action == "CreateKeyPair" {
					body, _ := bp.ObtainSdkValue("Result.PrivateKey", *resp)
					// save data into file
					if file, ok := d.GetOk("key_file"); ok {
						_ = ioutil.WriteFile(file.(string), []byte(body.(string)), 0600)
						_ = os.Chmod(file.(string), 0400)
					}
				}
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackEcsKeyPairService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyKeyPairAttribute",
			ConvertMode: bp.RequestConvertAll,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["KeyPairId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackEcsKeyPairService) RemoveResource(d *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteKeyPairs",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"KeyPairNames.1": d.Get("key_pair_name"),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading key pair on delete %q, %w", d.Id(), callErr))
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

func (s *VestackEcsKeyPairService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"key_pair_ids": {
				TargetField: "KeyPairIds",
				ConvertType: bp.ConvertWithN,
			},
			"key_pair_names": {
				TargetField: "KeyPairNames",
				ConvertType: bp.ConvertWithN,
			},
		},
		ResponseConverts: map[string]bp.ResponseConvert{
			"KeyPairId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
		NameField:    "KeyPairName",
		IdField:      "KeyPairId",
		CollectField: "key_pairs",
	}
}

func (VestackEcsKeyPairService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "ecs",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}
