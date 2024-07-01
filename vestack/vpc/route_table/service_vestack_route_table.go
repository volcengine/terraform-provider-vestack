package route_table

import (
	"errors"
	"fmt"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/vpc"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackRouteTableService struct {
	Client *bp.SdkClient
}

func NewRouteTableService(c *bp.SdkClient) *VestackRouteTableService {
	return &VestackRouteTableService{
		Client: c,
	}
}

func (s *VestackRouteTableService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackRouteTableService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		tables  []interface{}
		res     = make([]interface{}, 0)
		ids     interface{}
		idsMap  = make(map[string]bool)
		ok      bool
	)
	tables, err = bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		vpcClient := s.Client.VpcClient
		action := "DescribeRouteTableList"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = vpcClient.DescribeRouteTableListCommon(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = vpcClient.DescribeRouteTableListCommon(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = bp.ObtainSdkValue("Result.RouterTableList", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.RouterTableList is not Slice")
		}
		return data, err
	})
	if err != nil {
		return tables, err
	}
	ids, ok = m["RouteTableIds"]
	if !ok || ids == nil {
		return tables, nil
	}
	for _, id := range ids.(*schema.Set).List() {
		idsMap[strings.Trim(id.(string), " ")] = true
	}
	if len(idsMap) == 0 {
		return tables, nil
	}
	for _, entry := range tables {
		if _, ok = idsMap[entry.(map[string]interface{})["RouteTableId"].(string)]; ok {
			res = append(res, entry)
		}
	}
	return res, nil
}

func (s *VestackRouteTableService) ReadResource(resourceData *schema.ResourceData, tableId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if tableId == "" {
		tableId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"RouteTableId": tableId,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("value is not map")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("route table %s not exist ", tableId)
	}
	return data, err
}

func (s *VestackRouteTableService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				demo map[string]interface{}
			)
			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			_, err = bp.ObtainSdkValue("RouteTableId", demo)
			if err != nil {
				return nil, "", err
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, "Available", err
		},
	}
}

func (VestackRouteTableService) WithResourceResponseHandlers(tables map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return tables, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackRouteTableService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateRouteTable",
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.VpcClient.CreateRouteTableCommon(call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.RouteTableId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				vpc.NewVpcService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("vpc_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackRouteTableService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyRouteTableAttributes",
			ConvertMode: bp.RequestConvertAll,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["RouteTableId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.VpcClient.ModifyRouteTableAttributesCommon(call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				vpc.NewVpcService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("vpc_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackRouteTableService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteRouteTable",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"RouteTableId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.VpcClient.DeleteRouteTableCommon(call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading route table on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 3*time.Minute)
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackRouteTableService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "RouteTableIds",
			},
		},
		NameField:    "RouteTableName",
		IdField:      "RouteTableId",
		CollectField: "route_tables",
		ResponseConverts: map[string]bp.ResponseConvert{
			"RouteTableId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *VestackRouteTableService) ReadResourceId(id string) string {
	return id
}

func (s *VestackRouteTableService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "vpc",
		ResourceType:         "routetable",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}
