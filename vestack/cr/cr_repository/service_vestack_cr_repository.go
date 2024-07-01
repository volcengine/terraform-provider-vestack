package cr_repository

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackCrRepositoryService struct {
	Client *bp.SdkClient
}

func NewCrRepositoryService(c *bp.SdkClient) *VestackCrRepositoryService {
	return &VestackCrRepositoryService{
		Client: c,
	}
}

func (s *VestackCrRepositoryService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackCrRepositoryService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)

	pageCall := func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListRepositories"

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

func (s *VestackCrRepositoryService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		return data, fmt.Errorf("the id format must be 'registry:namespace:name'")
	}
	registry := parts[0]
	namespace := parts[1]
	name := parts[2]

	req := map[string]interface{}{
		"Registry": registry,
		"Filter": map[string]interface{}{
			"Namespaces": []string{namespace},
			"Names":      []string{name},
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
		return data, fmt.Errorf("CrRepository %s is not exist", name)
	}
	return data, err
}

func (s *VestackCrRepositoryService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, name string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *VestackCrRepositoryService) WithResourceResponseHandlers(instance map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return instance, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackCrRepositoryService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateRepository",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertAll,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				registry := d.Get("registry").(string)
				namespace := d.Get("namespace").(string)
				name := d.Get("name").(string)
				id := registry + ":" + namespace + ":" + name
				d.SetId(id)
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackCrRepositoryService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateRepository",
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Registry"] = resourceData.Get("registry").(string)
				(*call.SdkParam)["Namespace"] = resourceData.Get("namespace").(string)
				(*call.SdkParam)["Name"] = resourceData.Get("name").(string)
				if resourceData.HasChange("description") {
					(*call.SdkParam)["Description"] = resourceData.Get("description").(string)
				}
				if resourceData.HasChange("access_level") {
					(*call.SdkParam)["AccessLevel"] = resourceData.Get("access_level").(string)
				}
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

func (s *VestackCrRepositoryService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteRepository",
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				id := resourceData.Id()
				parts := strings.Split(id, ":")
				if len(parts) != 3 {
					return false, fmt.Errorf("the id format must be 'registry:namespace:name'")
				}
				(*call.SdkParam)["Registry"] = parts[0]
				(*call.SdkParam)["Namespace"] = parts[1]
				(*call.SdkParam)["Name"] = parts[2]
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

func (s *VestackCrRepositoryService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType:  bp.ContentTypeJson,
		IdField:      "Name",
		CollectField: "repositories",
		RequestConverts: map[string]bp.RequestConvert{
			"namespaces": {
				TargetField: "Filter.Namespaces",
				ConvertType: bp.ConvertJsonArray,
			},
			"names": {
				TargetField: "Filter.Names",
				ConvertType: bp.ConvertJsonArray,
			},
			"access_levels": {
				TargetField: "Filter.AccessLevels",
				ConvertType: bp.ConvertJsonArray,
			},
		},
	}
}

func (s *VestackCrRepositoryService) ReadResourceId(id string) string {
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
