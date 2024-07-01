package iam_access_key

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/encryption"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamAccessKeyService struct {
	Client *bp.SdkClient
}

func NewIamAccessKeyService(c *bp.SdkClient) *VestackIamAccessKeyService {
	return &VestackIamAccessKeyService{
		Client: c,
	}
}

func (s *VestackIamAccessKeyService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamAccessKeyService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
		idSet   = make(map[string]bool)
	)
	if _, ok = m["AccessKeyIds.1"]; ok {
		i := 1
		for {
			filed := fmt.Sprintf("AccessKeyIds.%d", i)
			tmpId, ok := m[filed]
			if !ok {
				break
			}
			idSet[tmpId.(string)] = true
			i++
			delete(m, filed)
		}
	}
	cens, err := bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		universalClient := s.Client.UniversalClient
		action := "ListAccessKeys"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = universalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = universalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, err
			}
		}
		logger.Debug(logger.RespFormat, action, resp)
		results, err = bp.ObtainSdkValue("Result.AccessKeyMetadata", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.AccessKeyMetadata is not Slice")
		}
		return data, err
	})
	if err != nil || len(idSet) == 0 {
		return cens, err
	}

	res := make([]interface{}, 0)
	for _, cen := range cens {
		if !idSet[cen.(map[string]interface{})["AccessKeyId"].(string)] {
			continue
		}
		res = append(res, cen)
	}
	return res, nil
}

func (s *VestackIamAccessKeyService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"AccessKeyIds.1": id,
	}
	if resourceData.Get("user_name") != nil && len(resourceData.Get("user_name").(string)) > 0 {
		req["UserName"] = resourceData.Get("user_name").(string)
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("access key %s not exist ", id)
	}
	return data, err
}

func (s *VestackIamAccessKeyService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				demo       map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "Error")
			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", demo)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("access key status error, status:%s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}

}

func (VestackIamAccessKeyService) WithResourceResponseHandlers(v map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return v, map[string]bp.ResponseConvert{}, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackIamAccessKeyService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callbacks := make([]bp.Callback, 0)

	// 创建ak
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateAccessKey",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"user_name": {
					ConvertType: bp.ConvertDefault,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, *resp)
				//注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := bp.ObtainSdkValue("Result.AccessKey.AccessKeyId", *resp)
				d.SetId(id.(string))
				sk, _ := bp.ObtainSdkValue("Result.AccessKey.SecretAccessKey", *resp)
				if v, ok := d.GetOk("pgp_key"); ok && len(v.(string)) > 0 {
					pgpKey := v.(string)
					encryptionKey, err := encryption.RetrieveGPGKey(pgpKey)
					if err != nil {
						return fmt.Errorf("get gpg key error: %s", err.Error())
					}
					fingerprint, encrypted, err := encryption.EncryptValue(encryptionKey, sk.(string), "Vestack IAM Access Key Secret")
					if err != nil {
						return fmt.Errorf("encrypt secret err: %s", err.Error())
					}
					_ = d.Set("key_fingerprint", fingerprint)
					_ = d.Set("encrypted_secret", encrypted)
				} else {
					_ = d.Set("secret", sk.(string))
				}
				if output, ok := d.GetOk("secret_file"); ok && output != nil {
					akSk, _ := bp.ObtainSdkValue("Result.AccessKey", *resp)
					if err := writeToFile(output.(string), akSk); err != nil {
						return fmt.Errorf("write secret to file err: %s", err.Error())
					}
				}

				return nil
			},
		},
	}
	callbacks = append(callbacks, callback)

	// 更新ak状态
	if resourceData.Get("status") != nil {
		callbacks = append(callbacks, s.updateAccessKeyStatus(resourceData.Get("status").(string), resourceData))
	}

	return callbacks
}

func (s *VestackIamAccessKeyService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callbacks := make([]bp.Callback, 0)
	if resourceData.HasChange("status") {
		callbacks = append(callbacks, s.updateAccessKeyStatus(resourceData.Get("status").(string), resourceData))
	}
	return callbacks
}

func (s *VestackIamAccessKeyService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callbacks := make([]bp.Callback, 0)

	// 删除前需要将ak禁用
	callbacks = append(callbacks, s.updateAccessKeyStatus("inactive", resourceData))

	// 删除sk
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:         "DeleteAccessKey",
			ConvertMode:    bp.RequestConvertIgnore,
			RequestIdField: "AccessKeyId",
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading access key on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
		},
	}
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackIamAccessKeyService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackIamAccessKeyService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "iam",
		Action:      actionName,
		Version:     "2018-01-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
	}
}

func (s *VestackIamAccessKeyService) updateAccessKeyStatus(status string, resourceData *schema.ResourceData) bp.Callback {
	return bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateAccessKey",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"Status": status,
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["AccessKeyId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading access key on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					logger.Debug(logger.ErrFormat, call.Action, callErr)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{status},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
}
