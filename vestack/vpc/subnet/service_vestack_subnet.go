package subnet

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/vpc"
)

type VestackSubnetService struct {
	Client *bp.SdkClient
}

func NewSubnetService(c *bp.SdkClient) *VestackSubnetService {
	return &VestackSubnetService{
		Client: c,
	}
}

func (s *VestackSubnetService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackSubnetService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		vpcClient := s.Client.VpcClient
		action := "DescribeSubnets"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = vpcClient.DescribeSubnetsCommon(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = vpcClient.DescribeSubnetsCommon(&condition)
			if err != nil {
				return data, err
			}
		}

		logger.Debug(logger.RespFormat, action, condition, *resp)

		results, err = bp.ObtainSdkValue("Result.Subnets", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Subnets is not Slice")
		}
		return data, err
	})
}

func (s *VestackSubnetService) ReadResource(resourceData *schema.ResourceData, subnetId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if subnetId == "" {
		subnetId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"SubnetIds.1": subnetId,
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
		return data, fmt.Errorf("Subnet %s not exist ", subnetId)
	}
	return data, err
}

func (s *VestackSubnetService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			failStates = append(failStates, "Error")
			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", demo)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("subnet status error, status:%s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}

}

func (VestackSubnetService) WithResourceResponseHandlers(subnet map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		if ipv6CidrBlock, ok1 := subnet["Ipv6CidrBlock"]; ok1 && ipv6CidrBlock.(string) != "" {
			subnet["EnableIpv6"] = true

			ipv6Address, _, err := net.ParseCIDR(ipv6CidrBlock.(string))
			if err != nil {
				return subnet, nil, err
			}
			bits := strings.Split(ipv6Address.String(), ":")
			if len(bits) < 4 {
				subnet["Ipv6CidrBlock"] = 0
			} else {
				temp := bits[3]
				temp = strings.Repeat("0", 4-len(temp)) + temp
				ipv6CidrValue, err := strconv.ParseInt(temp[2:], 16, 9)
				if err != nil {
					return subnet, nil, err
				}
				subnet["Ipv6CidrBlock"] = int(ipv6CidrValue)
			}

		} else {
			subnet["EnableIpv6"] = false
			delete(subnet, "Ipv6CidrBlock")
		}
		return subnet, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackSubnetService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action: "CreateSubnet",
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"ipv6_cidr_block": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ipv6CidrBlock, exists := d.GetOkExists("ipv6_cidr_block")
				if exists {
					(*call.SdkParam)["Ipv6CidrBlock"] = ipv6CidrBlock
				}

				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.VpcClient.CreateSubnetCommon(call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.SubnetId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
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

func (s *VestackSubnetService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifySubnetAttributes",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"ipv6_cidr_block": {
					Ignore: true,
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["SubnetId"] = d.Id()

				if d.HasChange("enable_ipv6") && d.Get("enable_ipv6").(bool) {
					ipv6CidrBlock, exists := d.GetOkExists("ipv6_cidr_block")
					if exists {
						(*call.SdkParam)["Ipv6CidrBlock"] = ipv6CidrBlock
					}
				}

				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.VpcClient.ModifySubnetAttributesCommon(call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackSubnetService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	id := resourceData.Id()
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteSubnet",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"SubnetId": id,
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("vpc_id").(string)
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.VpcClient.DeleteSubnetCommon(call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 3*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading subnet on delete %q, %w", d.Id(), callErr))
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

func (s *VestackSubnetService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "SubnetIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "SubnetName",
		IdField:      "SubnetId",
		CollectField: "subnets",
		ResponseConverts: map[string]bp.ResponseConvert{
			"SubnetId": {
				TargetField: "id",
			},
			"RouteTable.RouteTableId": {
				TargetField: "route_table_id",
			},
			"RouteTable.RouteTableType": {
				TargetField: "route_table_type",
			},
		},
	}
}

func (s *VestackSubnetService) ReadResourceId(id string) string {
	return id
}
