package cr_registry_state

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackCrRegistryStateService struct {
	Client *bp.SdkClient
}

func NewCrRegistryStateService(c *bp.SdkClient) *VestackCrRegistryStateService {
	return &VestackCrRegistryStateService{
		Client: c,
	}
}

func (s *VestackCrRegistryStateService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackCrRegistryStateService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)

	pageCall := func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListRegistries"

		logger.Debug(logger.ReqFormat, action, condition)
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

		logger.Debug(logger.RespFormat, action, resp)
		results, err = bp.ObtainSdkValue("Result.Items", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}

		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Results.Items is not slice")
		}

		return data, err
	}

	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, pageCall)
}

func (s *VestackCrRegistryStateService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)

	name := resourceData.Get("name").(string)

	req := map[string]interface{}{
		"Filter": map[string]interface{}{
			"Names": []string{name},
		},
	}

	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}

	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("value is not a map")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("CrRegistry %s is not exist", id)
	}
	return data, err
}

func (s *VestackCrRegistryStateService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			failStates = append(failStates, "Failed")
			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			logger.Debug("Refresh CrRegistry status resp:%v", "ReadResource", demo)

			status, err = bp.ObtainSdkValue("Status.Phase", demo)
			if err != nil {
				return nil, "", err
			}

			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("CrRegistry status error,status %s", status.(string))
				}
			}
			return demo, status.(string), err
		},
	}
}

func (VestackCrRegistryStateService) WithResourceResponseHandlers(instance map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return instance, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackCrRegistryStateService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	action := "StartRegistry"
	instanceAction := resourceData.Get("action").(string)
	if instanceAction != "Start" {
		return []bp.Callback{}
	}

	name := resourceData.Get("name")

	logger.DebugInfo("CreateResource state:name %s", name.(string))

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      action,
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"action": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				data, err := s.ReadResource(resourceData, name.(string))
				if err != nil {
					return false, err
				}

				// check instance status,only "{Stopped, [Released]}" can do start.
				phase, err := bp.ObtainSdkValue("Status.Phase", data)
				if err != nil {
					return false, err
				}
				if phase.(string) != "Stopped" {
					return false, fmt.Errorf("CrRegistry status is %s,can not start", phase.(string))
				}

				conditions, err := bp.ObtainSdkValue("Status.Conditions", data)
				if err != nil {
					return false, err
				}

				conds, ok := conditions.([]interface{})
				if !ok {
					return false, fmt.Errorf("CrRegistry status conditions not a slice,conditions:%v", conditions)
				}

				condHit := false
				for _, v := range conds {
					if v.(string) == "Released" {
						condHit = true
						break
					}
				}
				if !condHit {
					return false, fmt.Errorf("CrRegistry status is %s,not Released,can not start", name)
				}

				(*call.SdkParam)["Name"] = name
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id := d.Get("name").(string)
				d.SetId(id)
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Running"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}

	return []bp.Callback{callback}
}

func (s *VestackCrRegistryStateService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackCrRegistryStateService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackCrRegistryStateService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackCrRegistryStateService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "cr",
		Version:     "2022-05-12",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
