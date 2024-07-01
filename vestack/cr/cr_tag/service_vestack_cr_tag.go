package cr_tag

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

type VestackCrTagService struct {
	Client *bp.SdkClient
}

func NewCrTagService(c *bp.SdkClient) *VestackCrTagService {
	return &VestackCrTagService{
		Client: c,
	}
}

func (s *VestackCrTagService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackCrTagService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)

	pageCall := func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListTags"

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

		logger.Debug(logger.RespFormat, action, condition, *resp)
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

func (s *VestackCrTagService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	logger.DebugInfo("read resource id :%s", id)
	parts := strings.Split(id, ":")
	if len(parts) != 4 {
		return data, fmt.Errorf("the format of import id must be 'registry:namespace:repository:tag...'")
	}

	registry := parts[0]
	namespace := parts[1]
	repository := parts[2]
	name := parts[3]

	req := map[string]interface{}{
		"Registry":   registry,
		"Namespace":  namespace,
		"Repository": repository,
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
		return data, fmt.Errorf("cr tag %s not exist", name)
	}
	return data, err
}

func (s *VestackCrTagService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, name string) *resource.StateChangeConf {
	return nil
}

func (s *VestackCrTagService) WithResourceResponseHandlers(instance map[string]interface{}) []bp.ResourceResponseHandler {
	if _, ok := instance["ImageAttributes"]; !ok {
		instance["ImageAttributes"] = []interface{}{}
	}
	if _, ok := instance["ChartAttribute"]; !ok {
		instance["ChartAttribute"] = map[string]interface{}{}
	}
	return []bp.ResourceResponseHandler{}
}

func (s *VestackCrTagService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackCrTagService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackCrTagService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteTags",
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				id := resourceData.Id()
				parts := strings.Split(id, ":")
				if len(parts) != 4 {
					return false, fmt.Errorf("the id format must be 'registry:namespace:repository:tag'")
				}
				(*call.SdkParam)["Registry"] = parts[0]
				(*call.SdkParam)["Namespace"] = parts[1]
				(*call.SdkParam)["Repository"] = parts[2]
				(*call.SdkParam)["Names"] = []string{parts[3]}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackCrTagService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType:  bp.ContentTypeJson,
		IdField:      "Name",
		CollectField: "tags",
		RequestConverts: map[string]bp.RequestConvert{
			"names": {
				TargetField: "Filter.Names",
				ConvertType: bp.ConvertJsonArray,
			},
			"types": {
				TargetField: "Filter.Types",
				ConvertType: bp.ConvertJsonArray,
			},
		},
	}
}

func (s *VestackCrTagService) ReadResourceId(id string) string {
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
