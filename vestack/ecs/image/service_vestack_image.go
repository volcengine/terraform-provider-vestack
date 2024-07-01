package image

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackImageService struct {
	Client *bp.SdkClient
}

func NewImageService(c *bp.SdkClient) *VestackImageService {
	return &VestackImageService{
		Client: c,
	}
}

func (s *VestackImageService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackImageService) ReadResources(condition map[string]interface{}) ([]interface{}, error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithNextTokenQuery(condition, "MaxResults", "NextToken", 20, nil, func(m map[string]interface{}) (data []interface{}, next string, err error) {
		ecs := s.Client.EcsClient
		action := "DescribeInstances"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = ecs.DescribeImagesCommon(nil)
			if err != nil {
				return data, next, err
			}
		} else {
			resp, err = ecs.DescribeImagesCommon(&condition)
			if err != nil {
				return data, next, err
			}
		}
		logger.Debug(logger.RespFormat, action, condition, *resp)

		results, err = bp.ObtainSdkValue("Result.Images", *resp)
		if err != nil {
			return data, next, err
		}
		nextToken, err := bp.ObtainSdkValue("Result.NextToken", *resp)
		if err != nil {
			return data, next, err
		}
		next = nextToken.(string)
		if results == nil {
			results = []interface{}{}
		}

		if _, ok = results.([]interface{}); !ok {
			return data, next, errors.New("Result.Images is not Slice")
		}

		return results.([]interface{}), next, err
	})
}

func (s *VestackImageService) ReadResource(resourceData *schema.ResourceData, imageId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if imageId == "" {
		imageId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"ImageIds.1": imageId,
	}

	results, err = s.ReadResources(req)
	if err != nil {
		return nil, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return nil, errors.New("Value is not map ")
		}
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("Image %s not exist ", imageId)
	}
	return data, err
}

func (s *VestackImageService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *VestackImageService) WithResourceResponseHandlers(image map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return image, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackImageService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackImageService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackImageService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *VestackImageService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "ImageIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "ImageName",
		IdField:      "ImageId",
		CollectField: "images",
	}
}

func (s *VestackImageService) ReadResourceId(id string) string {
	return id
}
