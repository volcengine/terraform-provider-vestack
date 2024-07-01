package bucket_policy

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackTosBucketPolicyService struct {
	Client *bp.SdkClient
}

func NewTosBucketPolicyService(c *bp.SdkClient) *VestackTosBucketPolicyService {
	return &VestackTosBucketPolicyService{
		Client: c,
	}
}

func (s *VestackTosBucketPolicyService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackTosBucketPolicyService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	tos := s.Client.BypassSvcClient
	var (
		action  string
		resp    *map[string]interface{}
		results interface{}
	)
	action = "GetBucketPolicy"
	logger.Debug(logger.ReqFormat, action, nil)
	resp, err = tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     condition[bp.BypassDomain].(string),
		UrlParam: map[string]string{
			"policy": "",
		},
	}, nil)
	if err != nil {
		return data, err
	}
	results, err = bp.ObtainSdkValue(bp.BypassResponse, *resp)
	if err != nil {
		return data, err
	}

	if len(results.(map[string]interface{})) == 0 {
		return data, fmt.Errorf("bucket Policy %s not exist ", condition[bp.BypassDomain].(string))
	}

	data = append(data, map[string]interface{}{
		"Policy": results.(map[string]interface{}),
	})
	return data, err
}

func (s *VestackTosBucketPolicyService) ReadResource(resourceData *schema.ResourceData, instanceId string) (data map[string]interface{}, err error) {
	bucketName := resourceData.Get("bucket_name").(string)
	if instanceId == "" {
		instanceId = s.ReadResourceId(resourceData.Id())
	} else {
		instanceId = s.ReadResourceId(instanceId)
	}

	var (
		ok      bool
		results []interface{}
	)

	logger.Debug(logger.ReqFormat, "GetBucketPolicy", bucketName+":"+instanceId)
	condition := map[string]interface{}{
		bp.BypassDomain: bucketName,
	}
	results, err = s.ReadResources(condition)

	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, fmt.Errorf("Value is not map ")
		}
	}

	if len(data) == 0 {
		return data, fmt.Errorf("bucket Policy %s not exist ", instanceId)
	}

	return data, nil
}

func (s *VestackTosBucketPolicyService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *VestackTosBucketPolicyService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return m, map[string]bp.ResponseConvert{
			"Policy": {
				Convert: func(i interface{}) interface{} {
					b, _ := json.Marshal(i)
					return string(b)
				},
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackTosBucketPolicyService) putBucketPolicy(data *schema.ResourceData, resource *schema.Resource, isUpdate bool) bp.Callback {
	return bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "PutBucketPolicy",
			ConvertMode:     bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"bucket_name": {
					ConvertType: bp.ConvertDefault,
					TargetField: "BucketName",
					ForceGet:    isUpdate,
					SpecialParam: &bp.SpecialParam{
						Type: bp.DomainParam,
					},
				},
				"policy": {
					ForceGet:    isUpdate,
					ConvertType: bp.ConvertDefault,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				j := (*call.SdkParam)[bp.BypassParam].(map[string]interface{})["Policy"]
				data := map[string]interface{}{}
				err := json.Unmarshal([]byte(j.(string)), &data)
				if err != nil {
					return false, err
				}
				delete((*call.SdkParam)[bp.BypassParam].(map[string]interface{}), "Policy")
				for k, v := range data {
					(*call.SdkParam)[bp.BypassParam].(map[string]interface{})[k] = v
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				param := (*call.SdkParam)[bp.BypassParam].(map[string]interface{})
				return s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					HttpMethod:  bp.PUT,
					ContentType: bp.ApplicationJSON,
					UrlParam: map[string]string{
						"policy": "",
					},
					Domain: (*call.SdkParam)[bp.BypassDomain].(string),
				}, &param)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId((*call.SdkParam)[bp.BypassDomain].(string) + ":POLICY")
				return nil
			},
		},
	}
}

func (s *VestackTosBucketPolicyService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{s.putBucketPolicy(data, resource, false)}
}

func (s *VestackTosBucketPolicyService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{s.putBucketPolicy(data, resource, true)}
}

func (s *VestackTosBucketPolicyService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "DeleteBucketPolicy",
			ConvertMode:     bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"bucket_name": {
					ConvertType: bp.ConvertDefault,
					TargetField: "BucketName",
					ForceGet:    true,
					SpecialParam: &bp.SpecialParam{
						Type: bp.DomainParam,
					},
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					HttpMethod: bp.DELETE,
					Domain:     (*call.SdkParam)[bp.BypassDomain].(string),
					UrlParam: map[string]string{
						"policy": "",
					},
				}, nil)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading tos bucket policy on delete %q, %w", s.ReadResourceId(d.Id()), callErr))
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

func (s *VestackTosBucketPolicyService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackTosBucketPolicyService) ReadResourceId(id string) string {
	return id
}
