package tos_bucket_encryption

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackTosBucketEncryptionService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewTosBucketEncryptionService(c *bp.SdkClient) *VestackTosBucketEncryptionService {
	return &VestackTosBucketEncryptionService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackTosBucketEncryptionService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackTosBucketEncryptionService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return data, err
}

func (s *VestackTosBucketEncryptionService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	tos := s.Client.BypassSvcClient
	var (
		ok bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	action := "GetBucketEncryption"
	logger.Debug(logger.ReqFormat, action, id)
	resp, err := tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     id,
		UrlParam: map[string]string{
			"encryption": "",
		},
	}, nil)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp, err)
	if data, ok = (*resp)[bp.BypassResponse].(map[string]interface{}); !ok {
		return data, errors.New("GetBucketEncryption Resp is not map ")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("tos_bucket_encryption %s not exist ", id)
	}

	data["BucketName"] = id
	if v, exist := data["Rule"]; exist {
		if rule, ok := v.(map[string]interface{}); ok {
			if v1, exist1 := rule["ApplyServerSideEncryptionByDefault"]; exist1 {
				if encryption, ok1 := v1.(map[string]interface{}); ok1 {
					rule["ApplyServerSideEncryptionByDefault"] = []interface{}{encryption}
				}
			}
			data["Rule"] = []interface{}{rule}
		}
	}

	return data, err
}

func (s *VestackTosBucketEncryptionService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (VestackTosBucketEncryptionService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"SSEAlgorithm": {
				TargetField: "sse_algorithm",
			},
			"KMSDataEncryption": {
				TargetField: "kms_data_encryption",
			},
			"KMSMasterKeyID": {
				TargetField: "kms_master_key_id",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackTosBucketEncryptionService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateEncryption(resourceData, resource, false)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketEncryptionService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateEncryption(resourceData, resource, true)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketEncryptionService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteBucketEncryption",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["BucketName"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					ContentType: bp.ApplicationJSON,
					HttpMethod:  bp.DELETE,
					Domain:      (*call.SdkParam)["BucketName"].(string),
					UrlParam: map[string]string{
						"encryption": "",
					},
				}, nil)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				return resource.Retry(5*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading tos bucket encryption on delete %q, %w", s.ReadResourceId(d.Id()), callErr))
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

func (s *VestackTosBucketEncryptionService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackTosBucketEncryptionService) createOrUpdateEncryption(resourceData *schema.ResourceData, resource *schema.Resource, isUpdate bool) bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "PutBucketEncryption",
			ConvertMode:     bp.RequestConvertInConvert,
			ContentType:     bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"bucket_name": {
					ConvertType: bp.ConvertDefault,
					TargetField: "BucketName",
					SpecialParam: &bp.SpecialParam{
						Type: bp.DomainParam,
					},
					ForceGet: isUpdate,
				},
				"rule": {
					ConvertType: bp.ConvertJsonObject,
					TargetField: "Rule",
					NextLevelConvert: map[string]bp.RequestConvert{
						"apply_server_side_encryption_by_default": {
							ConvertType: bp.ConvertJsonObject,
							TargetField: "ApplyServerSideEncryptionByDefault",
							NextLevelConvert: map[string]bp.RequestConvert{
								"sse_algorithm": {
									ConvertType: bp.ConvertDefault,
									TargetField: "SSEAlgorithm",
								},
								"kms_data_encryption": {
									ConvertType: bp.ConvertDefault,
									TargetField: "KMSDataEncryption",
								},
								"kms_master_key_id": {
									ConvertType: bp.ConvertDefault,
									TargetField: "KMSMasterKeyID",
								},
							},
						},
					},
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				var sourceParam map[string]interface{}
				sourceParam, err := bp.SortAndStartTransJson((*call.SdkParam)[bp.BypassParam].(map[string]interface{}))
				if err != nil {
					return false, err
				}

				if rule, ok := sourceParam["Rule"].(map[string]interface{}); ok {
					if applyServerSideEncryptionByDefault, ok := rule["ApplyServerSideEncryptionByDefault"].(map[string]interface{}); ok {
						if sseAlgorithm, ok := applyServerSideEncryptionByDefault["SSEAlgorithm"]; ok {
							if sseAlgorithm != "kms" {
								delete(applyServerSideEncryptionByDefault, "KMSDataEncryption")
								delete(applyServerSideEncryptionByDefault, "KMSMasterKeyID")
							}
						}
					}
				}

				(*call.SdkParam)[bp.BypassParam] = sourceParam

				bytes, err := json.Marshal((*call.SdkParam)[bp.BypassParam].(map[string]interface{}))
				if err != nil {
					return false, err
				}
				hash := md5.New()
				io.WriteString(hash, string(bytes))
				contentMd5 := base64.StdEncoding.EncodeToString(hash.Sum(nil))

				(*call.SdkParam)[bp.BypassHeader].(map[string]string)["Content-MD5"] = contentMd5

				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)

				param := (*call.SdkParam)[bp.BypassParam].(map[string]interface{})
				resp, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					HttpMethod:  bp.PUT,
					ContentType: bp.ApplicationJSON,
					Domain:      (*call.SdkParam)[bp.BypassDomain].(string),
					Header:      (*call.SdkParam)[bp.BypassHeader].(map[string]string),
					UrlParam: map[string]string{
						"encryption": "",
					},
				}, &param)

				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId((*call.SdkParam)[bp.BypassDomain].(string))
				return nil
			},
		},
	}

	return callback
}

func (s *VestackTosBucketEncryptionService) ReadResourceId(id string) string {
	return id
}
