package iam_user_policy_attachment

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

type VestackIamUserPolicyAttachmentService struct {
	Client *bp.SdkClient
}

func NewIamUserPolicyAttachmentService(c *bp.SdkClient) *VestackIamUserPolicyAttachmentService {
	return &VestackIamUserPolicyAttachmentService{
		Client: c,
	}
}

func (s *VestackIamUserPolicyAttachmentService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamUserPolicyAttachmentService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	action := "ListAttachedUserPolicies"
	logger.Debug(logger.ReqFormat, action, m)
	if m == nil {
		resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &m)
		if err != nil {
			return data, err
		}
	}

	logger.Debug(logger.RespFormat, action, m, *resp)

	results, err = bp.ObtainSdkValue("Result.AttachedPolicyMetadata", *resp)
	if err != nil {
		return data, err
	}
	if results == nil {
		results = []interface{}{}
	}
	if data, ok = results.([]interface{}); !ok {
		return data, errors.New("Result.AttachedPolicyMetadata is not Slice")
	}
	return data, err
}

func (s *VestackIamUserPolicyAttachmentService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 3 {
		return data, fmt.Errorf("import id is invalid")
	}
	req := map[string]interface{}{
		"UserName": ids[0],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("value is not map")
		} else if ids[1] == data["PolicyName"].(string) && ids[2] == data["PolicyType"].(string) {
			return data, err
		}
	}
	return data, fmt.Errorf("User policy attachment %s not exist ", id)
}

func (s *VestackIamUserPolicyAttachmentService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *VestackIamUserPolicyAttachmentService) WithResourceResponseHandlers(attachment map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return attachment, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackIamUserPolicyAttachmentService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AttachUserPolicy",
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprintf("%s:%s:%s", d.Get("user_name").(string),
					d.Get("policy_name").(string), d.Get("policy_type").(string)))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackIamUserPolicyAttachmentService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackIamUserPolicyAttachmentService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DetachUserPolicy",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ids := strings.Split(d.Id(), ":")
				(*call.SdkParam)["UserName"] = ids[0]
				(*call.SdkParam)["PolicyName"] = ids[1]
				(*call.SdkParam)["PolicyType"] = ids[2]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
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
							return resource.NonRetryableError(fmt.Errorf("error on reading iam user attachment on delete %q, %w", d.Id(), callErr))
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

func (s *VestackIamUserPolicyAttachmentService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackIamUserPolicyAttachmentService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "iam",
		Version:     "2018-01-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}
