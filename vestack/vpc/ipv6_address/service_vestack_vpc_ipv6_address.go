package ipv6_address

import (
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackIpv6AddressService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewIpv6AddressService(c *bp.SdkClient) *VestackIpv6AddressService {
	return &VestackIpv6AddressService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *VestackIpv6AddressService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackIpv6AddressService) ReadResources(condition map[string]interface{}) (ipv6Addresses []interface{}, err error) {
	var (
		resp               *map[string]interface{}
		data               []interface{}
		results            interface{}
		next               string
		ok                 bool
		ecsInstance        map[string]interface{}
		networkInterfaces  []interface{}
		networkInterfaceId string
	)
	data, err = bp.WithNextTokenQuery(condition, "MaxResults", "NextToken", 20, nil, func(m map[string]interface{}) ([]interface{}, string, error) {
		action := "DescribeInstances"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = s.Client.UniversalClient.DoCall(getEcsUniversalInfo(action), nil)
			if err != nil {
				return data, next, err
			}
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getEcsUniversalInfo(action), &condition)
			if err != nil {
				return data, next, err
			}
		}
		logger.Debug(logger.RespFormat, action, condition, *resp)

		results, err = bp.ObtainSdkValue("Result.Instances", *resp)
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

		if data, ok = results.([]interface{}); !ok {
			return data, next, errors.New("Result.Instances is not Slice")
		}
		return data, next, err
	})

	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return ipv6Addresses, nil
	}
	if ecsInstance, ok = data[0].(map[string]interface{}); !ok {
		return ipv6Addresses, errors.New("Value is not map ")
	} else {
		// query primary network interface info of the ecs instance
		if networkInterfaces, ok = ecsInstance["NetworkInterfaces"].([]interface{}); !ok {
			return ipv6Addresses, errors.New("Instances.NetworkInterfaces is not Slice")
		}
		for _, networkInterface := range networkInterfaces {
			if networkInterfaceMap, ok := networkInterface.(map[string]interface{}); ok &&
				networkInterfaceMap["Type"] == "primary" {
				networkInterfaceId = networkInterfaceMap["NetworkInterfaceId"].(string)
			}
		}

		action := "DescribeNetworkInterfaces"
		req := map[string]interface{}{
			"NetworkInterfaceIds.1": networkInterfaceId,
		}
		logger.Debug(logger.ReqFormat, action, req)
		res, err := s.Client.UniversalClient.DoCall(getVpcUniversalInfo(action), &req)
		if err != nil {
			logger.Info("DescribeNetworkInterfaces error:", err)
			return ipv6Addresses, err
		}
		logger.Debug(logger.RespFormat, action, condition, *res)

		networkInterfaceInfos, err := bp.ObtainSdkValue("Result.NetworkInterfaceSets", *res)
		if err != nil {
			logger.Info("ObtainSdkValue Result.NetworkInterfaceSets error:", err)
			return ipv6Addresses, err
		}
		if ipv6Sets, ok := networkInterfaceInfos.([]interface{})[0].(map[string]interface{})["IPv6Sets"].([]interface{}); ok {
			for _, ipv6Address := range ipv6Sets {
				ipv6AddressMap := make(map[string]interface{})
				ipv6AddressMap["Ipv6Address"] = ipv6Address
				ipv6Addresses = append(ipv6Addresses, ipv6AddressMap)
			}
		}
	}

	return ipv6Addresses, err
}

func (s *VestackIpv6AddressService) ReadResource(resourceData *schema.ResourceData, allocationId string) (data map[string]interface{}, err error) {
	return data, err
}

func (s *VestackIpv6AddressService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (VestackIpv6AddressService) WithResourceResponseHandlers(ipv6Address map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return ipv6Address, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *VestackIpv6AddressService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}

func (s *VestackIpv6AddressService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}

func (s *VestackIpv6AddressService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return nil
}

func (s *VestackIpv6AddressService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"associated_instance_id": {
				TargetField: "InstanceIds.1",
			},
		},
		//IdField:      "AllocationId",
		CollectField: "ipv6_addresses",
	}
}

func (s *VestackIpv6AddressService) ReadResourceId(id string) string {
	return id
}

func getEcsUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "ecs",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}

func getVpcUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vpc",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
