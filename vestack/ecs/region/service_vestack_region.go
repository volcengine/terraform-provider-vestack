package region

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackRegionService struct {
	Client *bp.SdkClient
}

func (v *VestackRegionService) GetClient() *bp.SdkClient {
	return v.Client
}

func (v *VestackRegionService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp      *map[string]interface{}
		nextToken interface{}
		results   interface{}
		next      string
		ok        bool
	)
	return bp.WithNextTokenQuery(condition, "MaxResults", "NextToken", 10, nil, func(m map[string]interface{}) ([]interface{}, string, error) {
		action := "DescribeRegions"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = v.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
		} else {
			resp, err = v.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
		}
		if err != nil {
			return nil, next, err
		}
		logger.Debug(logger.RespFormat, action, condition, *resp)

		results, err = bp.ObtainSdkValue("Result.Regions", *resp)
		if err != nil {
			return nil, next, err
		}
		nextToken, err = bp.ObtainSdkValue("Result.NextToken", *resp)
		if err != nil {
			return nil, next, err
		}
		next, ok = nextToken.(string)
		if !ok {
			return nil, next, fmt.Errorf("next token must be a string")
		}
		if results == nil {
			results = make([]interface{}, 0)
		}

		if data, ok = results.([]interface{}); !ok {
			return nil, next, errors.New("Result.Regions is not Slice")
		}

		return data, next, err
	})
}

func (v *VestackRegionService) ReadResource(data *schema.ResourceData, s string) (map[string]interface{}, error) {
	return nil, nil
}

func (v *VestackRegionService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			return nil, "", err
		},
	}
}

func (v *VestackRegionService) WithResourceResponseHandlers(region map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return region, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (v *VestackRegionService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (v *VestackRegionService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (v *VestackRegionService) RemoveResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (v *VestackRegionService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "RegionIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "RegionId",
		IdField:      "RegionId",
		CollectField: "regions",
		ResponseConverts: map[string]bp.ResponseConvert{
			"RegionId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (v *VestackRegionService) ReadResourceId(s string) string {
	return s
}

func NewRegionService(c *bp.SdkClient) *VestackRegionService {
	return &VestackRegionService{
		Client: c,
	}
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "ecs",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}
