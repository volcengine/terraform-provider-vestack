package tos_bucket_customdomain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackTosBucketCustomDomainService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewTosBucketCustomDomainService(c *bp.SdkClient) *VestackTosBucketCustomDomainService {
	return &VestackTosBucketCustomDomainService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackTosBucketCustomDomainService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackTosBucketCustomDomainService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return data, err
}

func (s *VestackTosBucketCustomDomainService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	tos := s.Client.BypassSvcClient
	var (
		ok bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	bucketName, customDomain, err := parseId(id)
	if err != nil {
		return data, err
	}

	action := "GetBucketCustomDomain"
	logger.Debug(logger.ReqFormat, action, id)
	resp, err := tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     bucketName,
		UrlParam: map[string]string{
			"customdomain": "",
		},
	}, nil)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp, err)
	if data, ok = (*resp)[bp.BypassResponse].(map[string]interface{}); !ok {
		return data, errors.New("GetBucketCustomDomain Resp is not map")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("tos_bucket_customdomain %s not exist ", id)
	}

	// Find the specific custom domain
	if v, exist := data["CustomDomainRules"]; exist {
		if domains, ok := v.([]interface{}); ok {
			found := false
			for _, domain := range domains {
				if domainMap, ok := domain.(map[string]interface{}); ok {
					if domainName, exist := domainMap["Domain"]; exist && domainName.(string) == customDomain {
						data = domainMap
						data["BucketName"] = bucketName
						data["CustomDomain"] = customDomain
						found = true
						break
					}
				}
			}
			if !found {
				return data, fmt.Errorf("tos_bucket_customdomain %s not exist ", id)
			}
		}
	}

	return data, err
}

func (s *VestackTosBucketCustomDomainService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *VestackTosBucketCustomDomainService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	return []bp.ResourceResponseHandler{}
}

func (s *VestackTosBucketCustomDomainService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateCustomDomain(resourceData, resource, false)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketCustomDomainService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := s.createOrUpdateCustomDomain(resourceData, resource, true)
	callbacks = append(callbacks, callback)

	return callbacks
}

func (s *VestackTosBucketCustomDomainService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteBucketCustomDomain",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["BucketName"] = d.Id()
				bucketName, customDomain, err := parseId(d.Id())
				if err != nil {
					return false, err
				}
				(*call.SdkParam)["CustomDomain"] = customDomain
				(*call.SdkParam)["BucketName"] = bucketName
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					ContentType: bp.ApplicationJSON,
					HttpMethod:  bp.DELETE,
					Domain:      (*call.SdkParam)["BucketName"].(string),
					UrlParam: map[string]string{
						"customdomain": (*call.SdkParam)["CustomDomain"].(string),
					},
				}, nil)
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
							return resource.NonRetryableError(fmt.Errorf("error on reading tos bucket customdomain on delete %q, %w", s.ReadResourceId(d.Id()), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("bucket_name").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackTosBucketCustomDomainService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackTosBucketCustomDomainService) createOrUpdateCustomDomain(resourceData *schema.ResourceData, resource *schema.Resource, isUpdate bool) bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "PutBucketCustomDomain",
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
				"custom_domain_rule": {
					ConvertType: bp.ConvertJsonObject,
					TargetField: "CustomDomainRule",
					ForceGet:    isUpdate,
					NextLevelConvert: map[string]bp.RequestConvert{
						"domain": {
							ConvertType: bp.ConvertDefault,
							TargetField: "Domain",
							ForceGet:    isUpdate,
						},
						"protocol": {
							ConvertType: bp.ConvertDefault,
							TargetField: "Protocol",
							ForceGet:    isUpdate,
						},
						"cert_id": {
							ConvertType: bp.ConvertDefault,
							TargetField: "CertId",
							ForceGet:    isUpdate,
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
				(*call.SdkParam)[bp.BypassParam] = sourceParam

				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				param := (*call.SdkParam)[bp.BypassParam].(map[string]interface{})
				resp, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					ContentType: bp.ApplicationJSON,
					HttpMethod:  bp.PUT,
					Domain:      (*call.SdkParam)[bp.BypassDomain].(string),
					UrlParam: map[string]string{
						"customdomain": "",
					},
				}, &param)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId((*call.SdkParam)[bp.BypassDomain].(string) + ":" + (*call.SdkParam)[bp.BypassParam].(map[string]interface{})["CustomDomainRule"].(map[string]interface{})["Domain"].(string))
				return nil
			},
		},
	}

	return callback
}

func (s *VestackTosBucketCustomDomainService) ReadResourceId(id string) string {
	return id
}

func parseId(id string) (string, string, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid ID format, expected format: bucket_name:custom_domain")
	}
	return parts[0], parts[1], nil
}
