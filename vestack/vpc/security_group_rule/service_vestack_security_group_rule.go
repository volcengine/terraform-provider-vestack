package security_group_rule

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackSecurityGroupRuleService struct {
	Client *bp.SdkClient
}

type Direction string

const (
	DirectionIngress = Direction("ingress")
	DirectionEgress  = Direction("egress")
)

func NewSecurityGroupRuleService(c *bp.SdkClient) *VestackSecurityGroupRuleService {
	return &VestackSecurityGroupRuleService{
		Client: c,
	}
}

func (s *VestackSecurityGroupRuleService) GetClient() *bp.SdkClient {
	return s.Client
}

func (VestackSecurityGroupRuleService) WithResourceResponseHandlers(rule map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return rule, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackSecurityGroupRuleService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var action string
	direction := resourceData.Get("direction").(string)
	if direction == string(DirectionEgress) {
		action = "AuthorizeSecurityGroupEgress"
	} else {
		action = "AuthorizeSecurityGroupIngress"
	}

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      action,
			ConvertMode: bp.RequestConvertAll,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if resourceData.Get("cidr_ip") == "" && resourceData.Get("source_group_id") == "" {
					return false, fmt.Errorf("At least one of cidr_ip and source_group_id exists. ")
				}
				protocol := resourceData.Get("protocol").(string)
				start := resourceData.Get("port_start").(int)
				end := resourceData.Get("port_end").(int)
				if err := validateProtocol(protocol, start, end); err != nil {
					return false, err
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				if direction == string(DirectionEgress) {
					return s.Client.VpcClient.AuthorizeSecurityGroupEgressCommon(call.SdkParam)
				} else {
					return s.Client.VpcClient.AuthorizeSecurityGroupIngressCommon(call.SdkParam)
				}
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				var (
					securityGroupId = resourceData.Get("security_group_id")
					cidrIp          = resourceData.Get("cidr_ip")
					protocol        = resourceData.Get("protocol")
					portStart       = resourceData.Get("port_start")
					portEnd         = resourceData.Get("port_end")
					sourceGroupId   = resourceData.Get("source_group_id")
					dir             = resourceData.Get("direction")
					policy          = resourceData.Get("policy")
					priority        = resourceData.Get("priority")
				)
				d.SetId(fmt.Sprintf("%v:%v:%v:%v:%v:%v:%v:%v:%v",
					securityGroupId, protocol, portStart,
					portEnd, cidrIp, sourceGroupId,
					dir, policy, priority))
				return nil
			},
		},
	}

	return []bp.Callback{callback}
}

func (s *VestackSecurityGroupRuleService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)

	vpcClient := s.Client.VpcClient
	action := "DescribeSecurityGroupAttributes"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = vpcClient.DescribeSecurityGroupAttributesCommon(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = vpcClient.DescribeSecurityGroupAttributesCommon(&condition)
		if err != nil {
			return data, err
		}
	}
	logger.Debug(logger.RespFormat, action, condition, *resp)

	results, err = bp.ObtainSdkValue("Result.Permissions", *resp)

	if err != nil {
		return data, err
	}
	if results == nil {
		results = []interface{}{}
	}
	if data, ok = results.([]interface{}); !ok {
		return data, errors.New("Result.Permissions is not Slice")
	}

	securityGroupStatus, err := bp.ObtainSdkValue("Result.Status", *resp)
	if err != nil {
		return nil, err
	}

	// Resource 里定义了 status，已经发布无法修改，这里将 status 回填
	for index := range data {
		data[index].(map[string]interface{})["Status"] = securityGroupStatus

		ele := data[index].(map[string]interface{})
		data[index].(map[string]interface{})["PortStart"] = int(ele["PortStart"].(float64))
		data[index].(map[string]interface{})["PortEnd"] = int(ele["PortEnd"].(float64))
		data[index].(map[string]interface{})["Priority"] = int(ele["Priority"].(float64))
	}

	return data, err
}

func (s *VestackSecurityGroupRuleService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	ids := strings.Split(id, ":")

	req := map[string]interface{}{
		"SecurityGroupId": ids[0],
		"Direction":       resourceData.Get("direction"),
		"Protocol":        resourceData.Get("protocol"),
	}
	if len(resourceData.Get("cidr_ip").(string)) > 0 {
		req["CidrIp"] = resourceData.Get("cidr_ip")
	}
	if len(resourceData.Get("source_group_id").(string)) > 0 {
		req["SourceGroupId"] = resourceData.Get("source_group_id")
	}

	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}

	for _, v := range results {
		data = v.(map[string]interface{})

		if data["PortStart"] != resourceData.Get("port_start") {
			continue
		}
		if data["PortEnd"] != resourceData.Get("port_end") {
			continue
		}
		if data["Policy"] != resourceData.Get("policy") {
			continue
		}
		if data["Priority"] != resourceData.Get("priority") {
			continue
		}
		return data, nil
	}

	return data, fmt.Errorf("SecurityGroupRule %s not exist ", id)
}

func (s *VestackSecurityGroupRuleService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *VestackSecurityGroupRuleService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var action string
	direction := resourceData.Get("direction").(string)
	if direction == string(DirectionEgress) {
		action = "ModifySecurityGroupRuleDescriptionsEgress"
	} else {
		action = "ModifySecurityGroupRuleDescriptionsIngress"
	}

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      action,
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"description": {
					TargetField: "Description",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				items := strings.Split(d.Id(), ":")

				start, _ := strconv.Atoi(items[2])
				end, _ := strconv.Atoi(items[3])

				(*call.SdkParam)["SecurityGroupId"] = items[0]
				(*call.SdkParam)["Protocol"] = items[1]
				(*call.SdkParam)["PortStart"] = start
				(*call.SdkParam)["PortEnd"] = end
				if len(items[4]) > 0 {
					(*call.SdkParam)["CidrIp"] = items[4]
				}
				if len(items[5]) > 0 {
					(*call.SdkParam)["SourceGroupId"] = items[5]
				}
				(*call.SdkParam)["Policy"] = items[7]
				(*call.SdkParam)["Priority"] = resourceData.Get("priority")

				// validate protocol
				if err := validateProtocol(items[1], start, end); err != nil {
					return false, err
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				if direction == string(DirectionEgress) {
					return s.Client.VpcClient.ModifySecurityGroupRuleDescriptionsEgressCommon(call.SdkParam)
				} else {
					return s.Client.VpcClient.ModifySecurityGroupRuleDescriptionsIngressCommon(call.SdkParam)
				}

			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackSecurityGroupRuleService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	var action string
	direction := resourceData.Get("direction").(string)
	if direction == string(DirectionEgress) {
		action = "RevokeSecurityGroupEgress"
	} else {
		action = "RevokeSecurityGroupIngress"
	}

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      action,
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"SecurityGroupId": resourceData.Get("security_group_id"),
				"Protocol":        resourceData.Get("protocol"),
				"PortStart":       resourceData.Get("port_start"),
				"PortEnd":         resourceData.Get("port_end"),
				"CidrIp":          resourceData.Get("cidr_ip"),
				"SourceGroupId":   resourceData.Get("source_group_id"),
				"Policy":          resourceData.Get("policy"),
				"Priority":        resourceData.Get("priority"),
			},

			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				if direction == string(DirectionEgress) {
					return s.Client.VpcClient.RevokeSecurityGroupEgressCommon(call.SdkParam)
				} else {
					return s.Client.VpcClient.RevokeSecurityGroupIngressCommon(call.SdkParam)
				}

			},
		},
	}
	return []bp.Callback{callback}
}

func (s *VestackSecurityGroupRuleService) DatasourceResources(d *schema.ResourceData, r *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		CollectField: "security_group_rules",
		ExtraData: func(sourceData []interface{}) ([]interface{}, error) {
			var next []interface{}
			for _, i := range sourceData {
				i.(map[string]interface{})["SecurityGroupId"] = d.Get("security_group_id")
				next = append(next, i)
			}
			return next, nil
		},
	}
}

func (s *VestackSecurityGroupRuleService) ReadResourceId(id string) string {
	return id
}

func validateProtocol(protocol string, start, end int) error {
	switch protocol {
	case "tcp":
		if start < 1 || end < 1 {
			return fmt.Errorf("Protocol is tcp,Port start or end must between 1-65535. ")
		}
	case "udp":
		if start < 1 || end < 1 {
			return fmt.Errorf("Protocol is udp,Port start or end must between 1-65535. ")
		}
	case "icmp":
		if start != -1 || end != -1 {
			return fmt.Errorf("Protocol is icmp,Port start or end must -1. ")
		}
	case "all":
		if start != -1 || end != -1 {
			return fmt.Errorf("Protocol is all,Port start or end must -1. ")
		}
	}
	return nil
}
