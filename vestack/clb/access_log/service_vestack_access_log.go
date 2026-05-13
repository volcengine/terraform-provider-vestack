package access_log

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

type VestackAccessLogService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewAccessLogService(c *bp.SdkClient) *VestackAccessLogService {
	return &VestackAccessLogService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackAccessLogService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackAccessLogService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	data, err = bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeLoadBalancers"

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
		results, err = bp.ObtainSdkValue("Result.LoadBalancers", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.LoadBalancers is not Slice")
		}
		return data, err
	})
	if err != nil {
		return data, err
	}
	return data, err
}

func (s *VestackAccessLogService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"LoadBalancerIds.1": id,
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
		return data, fmt.Errorf("clb %s not exist ", id)
	}
	return data, err
}

func (s *VestackAccessLogService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				d      map[string]interface{}
				status interface{}
			)
			d, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", d)
			if err != nil {
				return nil, "", err
			}
			return d, status.(string), err
		},
	}
}

func (s *VestackAccessLogService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "EnableAccessLog",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"load_balancer_id": {
					TargetField: "LoadBalancerId",
				},
				"delivery_type": {
					TargetField: "DeliveryType",
				},
				"bucket_name": {
					TargetField: "BucketName",
				},
				"tls_project_id": {
					TargetField: "TlsProjectId",
				},
				"tls_topic_id": {
					TargetField: "TlsTopicId",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				// Validate required fields based on delivery_type
				deliveryType, ok := d.GetOk("delivery_type")
				if !ok || deliveryType == "tos" {
					// For TOS delivery, bucket_name is required
					if _, ok := d.GetOk("bucket_name"); !ok {
						return false, fmt.Errorf("bucket_name is required when delivery_type is 'tos'")
					}
				} else if deliveryType == "tls" {
					// For TLS delivery, both tls_project_id and tls_topic_id are required
					if _, ok := d.GetOk("tls_project_id"); !ok {
						return false, fmt.Errorf("tls_project_id is required when delivery_type is 'tls'")
					}
					if _, ok := d.GetOk("tls_topic_id"); !ok {
						return false, fmt.Errorf("tls_topic_id is required when delivery_type is 'tls'")
					}
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				// Set resource ID as LoadBalancerId for access log
				loadBalancerId, _ := bp.ObtainSdkValue("LoadBalancerId", *call.SdkParam)
				d.SetId(loadBalancerId.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("load_balancer_id").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (VestackAccessLogService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackAccessLogService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	// Access log doesn't support update operation
	return []bp.Callback{}
}

func (s *VestackAccessLogService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	clbId := resourceData.Get("load_balancer_id").(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DisableAccessLog",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["LoadBalancerId"] = clbId
				if deliveryType, ok := d.GetOk("delivery_type"); ok {
					(*call.SdkParam)["DeliveryType"] = deliveryType.(string)
				} else {
					(*call.SdkParam)["DeliveryType"] = "tos"
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			LockId: func(d *schema.ResourceData) string {
				return clbId
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackAccessLogService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackAccessLogService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "clb",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
