package route_table_associate

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/subnet"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/vpc"
)

type VestackRouteTableAssociateService struct {
	Client *bp.SdkClient
}

func NewRouteTableAssociateService(c *bp.SdkClient) *VestackRouteTableAssociateService {
	return &VestackRouteTableAssociateService{
		Client: c,
	}
}

func (s *VestackRouteTableAssociateService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackRouteTableAssociateService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
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
}

func (s *VestackRouteTableAssociateService) ReadResource(resourceData *schema.ResourceData, associateId string) (data map[string]interface{}, err error) {
	var (
		results        []interface{}
		ok             bool
		associate      bool
		subnetIds      interface{}
		tmpSubnetIds   []interface{}
		routeTableId   string
		targetSubnetId string
		ids            []string
	)

	if associateId == "" {
		associateId = s.ReadResourceId(resourceData.Id())
	}

	ids = strings.Split(associateId, ":")
	if len(ids) != 2 {
		return map[string]interface{}{}, fmt.Errorf("invalid route table associate id: %v", associateId)
	}
	routeTableId = ids[0]
	targetSubnetId = ids[1]

	req := map[string]interface{}{
		"RouteTableId": routeTableId,
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
		return data, fmt.Errorf("route table %s not exist ", routeTableId)
	}
	subnetIds, err = bp.ObtainSdkValue("SubnetIds", data)
	if err != nil {
		return data, err
	}
	if subnetIds == nil {
		return data, errors.New("not associate")
	}
	tmpSubnetIds, ok = subnetIds.([]interface{})
	if !ok {
		return data, errors.New("subnet ids is not string slice")
	}
	for _, subnetId := range tmpSubnetIds {
		if subnetId.(string) == targetSubnetId {
			associate = true
			break
		}
	}
	if !associate {
		return data, errors.New("not associate")
	}
	return data, err
}

func (s *VestackRouteTableAssociateService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				demo   map[string]interface{}
				status = "Associate"
			)
			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				if !strings.Contains(err.Error(), "not associate") {
					return nil, "", err
				}
				status = "Available"
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status, nil
		},
	}
}

func (VestackRouteTableAssociateService) WithResourceResponseHandlers(routeTables map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return routeTables, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackRouteTableAssociateService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var vpcId string
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AssociateRouteTable",
			ConvertMode: bp.RequestConvertAll,
			LockId: func(d *schema.ResourceData) string {
				return vpcId
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				subnetId := resourceData.Get("subnet_id").(string)
				resp, err := subnet.NewSubnetService(s.Client).ReadResource(resourceData, subnetId)
				if err != nil {
					return false, err
				}
				vpcId = resp["VpcId"].(string)
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.VpcClient.AssociateRouteTableCommon(call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprint((*call.SdkParam)["RouteTableId"], ":", (*call.SdkParam)["SubnetId"]))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Associate"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			// 外部定义vpcId无法传入ExtraRefresh中
			ExtraRefreshCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (map[bp.ResourceService]*bp.StateRefresh, error) {
				return map[bp.ResourceService]*bp.StateRefresh{
					vpc.NewVpcService(s.Client): {
						Target:     []string{"Available"},
						Timeout:    resourceData.Timeout(schema.TimeoutCreate),
						ResourceId: vpcId,
					},
				}, nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackRouteTableAssociateService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}

func (s *VestackRouteTableAssociateService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	var vpcId string
	ids := strings.Split(s.ReadResourceId(resourceData.Id()), ":")
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DisassociateRouteTable",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"RouteTableId": ids[0],
				"SubnetId":     ids[1],
			},
			LockId: func(d *schema.ResourceData) string {
				return vpcId
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				subnetId := resourceData.Get("subnet_id").(string)
				resp, err := subnet.NewSubnetService(s.Client).ReadResource(resourceData, subnetId)
				if err != nil {
					return false, err
				}
				vpcId = resp["VpcId"].(string)
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.VpcClient.DisassociateRouteTableCommon(call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutDelete),
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						return resource.NonRetryableError(fmt.Errorf("error on reading route table associate on delete %q, %w", d.Id(), callErr))
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

func (s *VestackRouteTableAssociateService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *VestackRouteTableAssociateService) ReadResourceId(id string) string {
	return id
}
