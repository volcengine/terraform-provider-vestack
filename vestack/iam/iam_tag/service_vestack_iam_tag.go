package iam_tag

import (
	"errors"
	"fmt"
	"strconv"

	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIamTagService struct {
	Client *bp.SdkClient
}

func NewIamTagService(c *bp.SdkClient) *VestackIamTagService {
	return &VestackIamTagService{
		Client: c,
	}
}

func (s *VestackIamTagService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIamTagService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	maxResults := 100
	if val, ok := m["MaxResults"].(int); ok && val > 0 {
		maxResults = val
	}
	return bp.WithNextTokenQuery(m, "MaxResults", "NextToken", maxResults, nil, func(condition map[string]interface{}) ([]interface{}, string, error) {
		universalClient := s.Client.UniversalClient
		action := "ListTagsForResources"

		// Convert MaxResults to string as required by IAM API
		if val, ok := condition["MaxResults"].(int); ok {
			condition["MaxResults"] = strconv.Itoa(val)
		}

		if _, ok := condition["NextToken"]; !ok {
			condition["NextToken"] = ""
		}

		logger.Debug(logger.ReqFormat, action, condition)

		resp, err := universalClient.DoCall(getUniversalInfo(action), &condition)
		if err != nil {
			return nil, "", err
		}

		logger.Debug(logger.RespFormat, action, resp)
		results, err := bp.ObtainSdkValue("Result.ResourceTags", *resp)
		if err != nil {
			return nil, "", err
		}
		if results == nil {
			results = []interface{}{}
		}
		rawTags, ok := results.([]interface{})
		if !ok {
			return nil, "", errors.New("Result.ResourceTags is not Slice")
		}

		var nextTokenStr string
		nextToken, _ := bp.ObtainSdkValue("Result.NextToken", *resp)
		if nextToken != nil {
			if s, ok := nextToken.(string); ok {
				nextTokenStr = s
			}
		}

		return rawTags, nextTokenStr, nil
	})
}

func (s *VestackIamTagService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	// For a single tag resource, we might need to filter from ListTagsForResources
	return nil, errors.New("ReadResource not implemented for iam_tag")
}

func (s *VestackIamTagService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (VestackIamTagService) WithResourceResponseHandlers(v map[string]interface{}) []bp.ResourceResponseHandler {
	return nil
}

func (s *VestackIamTagService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "TagResources",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ResourceType"] = d.Get("resource_type")
				vNames := d.Get("resource_names")
				resourceNames, ok := vNames.([]interface{})
				if !ok {
					return false, errors.New("resource_names is not []interface{}")
				}
				for i, name := range resourceNames {
					nameStr, ok := name.(string)
					if !ok {
						return false, errors.New("resource_name is not string")
					}
					(*call.SdkParam)[fmt.Sprintf("ResourceNames.%d", i+1)] = nameStr
				}
				vTags := d.Get("tags")
				tagsSet, ok := vTags.(*schema.Set)
				if !ok {
					return false, errors.New("tags is not *schema.Set")
				}
				tags := tagsSet.List()
				for i, tag := range tags {
					tm, ok := tag.(map[string]interface{})
					if !ok {
						return false, errors.New("tag item is not map")
					}
					(*call.SdkParam)[fmt.Sprintf("Tags.%d.Key", i+1)] = tm["key"]
					(*call.SdkParam)[fmt.Sprintf("Tags.%d.Value", i+1)] = tm["value"]
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprintf("%s:%d", d.Get("resource_type"), time.Now().UnixNano()))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackIamTagService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	// Tags are usually recreated if changed, or handled by adding/removing.
	// For this standalone resource, we'll just recreate it if key/value changes.
	return nil
}

func (s *VestackIamTagService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UntagResources",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ResourceType"] = d.Get("resource_type")
				vNames := d.Get("resource_names")
				resourceNames, ok := vNames.([]interface{})
				if !ok {
					return false, errors.New("resource_names is not []interface{}")
				}
				for i, name := range resourceNames {
					nameStr, ok := name.(string)
					if !ok {
						return false, errors.New("resource_name is not string")
					}
					(*call.SdkParam)[fmt.Sprintf("ResourceNames.%d", i+1)] = nameStr
				}
				vTags := d.Get("tags")
				tagsSet, ok := vTags.(*schema.Set)
				if !ok {
					return false, errors.New("tags is not *schema.Set")
				}
				tags := tagsSet.List()
				for i, tag := range tags {
					tm, ok := tag.(map[string]interface{})
					if !ok {
						return false, errors.New("tag item is not map")
					}
					(*call.SdkParam)[fmt.Sprintf("TagKeys.%d", i+1)] = tm["key"]
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

func (s *VestackIamTagService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"resource_type": {
				TargetField: "ResourceType",
			},
			"resource_names": {
				TargetField: "ResourceNames",
				ConvertType: bp.ConvertListN,
			},
		},
		CollectField: "resource_tags",
		ResponseConverts: map[string]bp.ResponseConvert{
			"ResourceType": {
				TargetField: "resource_type",
			},
			"ResourceName": {
				TargetField: "resource_name",
			},
			"TagKey": {
				TargetField: "tag_key",
			},
			"TagValue": {
				TargetField: "tag_value",
			},
			"NextToken": {
				TargetField: "next_token",
				Chain:       "resource_tags",
			},
		},
	}
}

func (s *VestackIamTagService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "iam",
		Action:      actionName,
		Version:     "2018-01-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		RegionType:  bp.Global,
	}
}
