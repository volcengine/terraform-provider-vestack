package kubeconfig

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
	"time"
)

type VestackVkeKubeconfigService struct {
	Client *bp.SdkClient
}

func NewVkeKubeconfigService(c *bp.SdkClient) *VestackVkeKubeconfigService {
	return &VestackVkeKubeconfigService{
		Client: c,
	}
}

func (s *VestackVkeKubeconfigService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackVkeKubeconfigService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)

	data, err = bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 10, 1, func(m map[string]interface{}) ([]interface{}, error) {
		action := "ListKubeconfigs"
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

		results, err = bp.ObtainSdkValue("Result.Items", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Items is not Slice")
		}
		return data, err
	})
	if err != nil {
		return data, err
	}

	return data, err
}

func (s *VestackVkeKubeconfigService) ReadResource(resourceData *schema.ResourceData, kubeconfigId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if kubeconfigId == "" {
		kubeconfigId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"Filter": map[string]interface{}{
			"Ids": []string{kubeconfigId},
		},
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
		return data, fmt.Errorf("Vke Kubeconfig %s not exist ", kubeconfigId)
	}

	return data, err
}

func (s *VestackVkeKubeconfigService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh:    nil, // no need to refresh resources
	}
}

func (VestackVkeKubeconfigService) WithResourceResponseHandlers(kubeconfig map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return kubeconfig, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackVkeKubeconfigService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateKubeconfig",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//创建Kubeconfig
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := bp.ObtainSdkValue("Result.Id", *resp)
				d.SetId(id.(string))
				//异步操作 这里需要等5s
				time.Sleep(time.Duration(5) * time.Second)
				return nil
			},
		},
	}
	return []bp.Callback{callback}

}

func (s *VestackVkeKubeconfigService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackVkeKubeconfigService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteKubeconfigs",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Ids"] = []string{resourceData.Id()}
				(*call.SdkParam)["ClusterId"] = resourceData.Get("cluster_id")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//删除Kubeconfig
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//异步操作 这里需要等5s
				time.Sleep(time.Duration(5) * time.Second)
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 10*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading kubeconfig on delete %q, %w", d.Id(), callErr))
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

func (s *VestackVkeKubeconfigService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "Filter.Ids",
				ConvertType: bp.ConvertJsonArray,
			},
			"cluster_ids": {
				TargetField: "Filter.ClusterIds",
				ConvertType: bp.ConvertJsonArray,
			},
			"types": {
				TargetField: "Filter.Types",
				ConvertType: bp.ConvertJsonArray,
			},
		},
		ContentType:  bp.ContentTypeJson,
		IdField:      "Id",
		CollectField: "kubeconfigs",
		ResponseConverts: map[string]bp.ResponseConvert{
			"Id": {
				TargetField: "kubeconfig_id",
				KeepDefault: true,
			},
		},
	}
}

func (s *VestackVkeKubeconfigService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vke",
		Version:     "2022-05-12",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
