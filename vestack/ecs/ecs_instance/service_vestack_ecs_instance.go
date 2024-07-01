package ecs_instance

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_deployment_set_associate"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/subnet"
	"golang.org/x/sync/semaphore"
	"golang.org/x/time/rate"
)

var rateInfo *bp.RateInfo

func init() {
	rateInfo = &bp.RateInfo{
		Create: &bp.Rate{
			Limiter:   rate.NewLimiter(4, 10),
			Semaphore: semaphore.NewWeighted(14),
		},
		Update: &bp.Rate{
			Limiter:   rate.NewLimiter(4, 10),
			Semaphore: semaphore.NewWeighted(14),
		},
		Read: &bp.Rate{
			Limiter:   rate.NewLimiter(4, 10),
			Semaphore: semaphore.NewWeighted(14),
		},
		Delete: &bp.Rate{
			Limiter:   rate.NewLimiter(4, 10),
			Semaphore: semaphore.NewWeighted(14),
		},
		Data: &bp.Rate{
			Limiter:   rate.NewLimiter(4, 10),
			Semaphore: semaphore.NewWeighted(14),
		},
	}
}

type VestackEcsService struct {
	Client        *bp.SdkClient
	SubnetService *subnet.VestackSubnetService
}

func NewEcsService(c *bp.SdkClient) *VestackEcsService {
	return &VestackEcsService{
		Client:        c,
		SubnetService: subnet.NewSubnetService(c),
	}
}

func (s *VestackEcsService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackEcsService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp               *map[string]interface{}
		results            interface{}
		next               string
		ok                 bool
		ecsInstance        map[string]interface{}
		networkInterfaces  []interface{}
		networkInterfaceId string
	)
	data, err = bp.WithNextTokenQuery(condition, "MaxResults", "NextToken", 20, nil, func(m map[string]interface{}) ([]interface{}, string, error) {
		ecs := s.Client.EcsClient
		action := "DescribeInstances"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = ecs.DescribeInstancesCommon(nil)
			if err != nil {
				return data, next, err
			}
		} else {
			resp, err = ecs.DescribeInstancesCommon(&condition)
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
		data, err = RemoveSystemTags(data)
		return data, next, err
	})

	if err != nil {
		return nil, err
	}

	for _, v := range data {
		if ecsInstance, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		} else {
			// query primary network interface info of the ecs instance
			if networkInterfaces, ok = ecsInstance["NetworkInterfaces"].([]interface{}); !ok {
				return data, errors.New("Instances.NetworkInterfaces is not Slice")
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
				continue
			}
			logger.Debug(logger.RespFormat, action, condition, *res)

			networkInterfaceInfos, err := bp.ObtainSdkValue("Result.NetworkInterfaceSets", *res)
			if err != nil {
				logger.Info("ObtainSdkValue Result.NetworkInterfaceSets error:", err)
				continue
			}
			if ipv6Sets, ok := networkInterfaceInfos.([]interface{})[0].(map[string]interface{})["IPv6Sets"].([]interface{}); ok {
				ecsInstance["Ipv6Addresses"] = ipv6Sets
				ecsInstance["Ipv6AddressCount"] = len(ipv6Sets)
			}
		}
	}

	return data, err
}

func (s *VestackEcsService) ReadResource(resourceData *schema.ResourceData, instanceId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if instanceId == "" {
		instanceId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"InstanceIds.1": instanceId,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, fmt.Errorf("Value is not map ")
		}
	}

	if len(data) == 0 {
		return data, fmt.Errorf("Ecs Instance %s not exist ", instanceId)
	}
	return data, nil
}

func (s *VestackEcsService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				data       map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "ERROR")

			if err = resource.Retry(20*time.Minute, func() *resource.RetryError {
				data, err = s.ReadResource(resourceData, id)
				if err != nil {
					if bp.ResourceNotFoundError(err) {
						return resource.RetryableError(err)
					} else {
						return resource.NonRetryableError(err)
					}
				}
				return nil
			}); err != nil {
				return nil, "", err
			}

			status, err = bp.ObtainSdkValue("Status", data)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("Ecs Instance  status  error, status:%s", status.(string))
				}
			}
			project, err := bp.ObtainSdkValue("ProjectName", data)
			if err != nil {
				return nil, "", err
			}
			if resourceData.Get("project_name") != nil && resourceData.Get("project_name").(string) != "" {
				if project != resourceData.Get("project_name") {
					return data, "", err
				}
			}
			return data, status.(string), err
		},
	}
}

func (s *VestackEcsService) WithResourceResponseHandlers(ecs map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		var (
			typeErr             error
			ebsErr              error
			userDataErr         error
			networkInterfaceErr error
			errorStr            string
			wg                  sync.WaitGroup
			syncMap             sync.Map
		)
		// 使用小写的 Hostname
		delete(ecs, "HostName")

		//计算period
		if ecs["InstanceChargeType"].(string) == "PrePaid" {
			ct, _ := time.Parse("2006-01-02T15:04:05", ecs["CreatedAt"].(string)[0:strings.Index(ecs["CreatedAt"].(string), "+")])
			et, _ := time.Parse("2006-01-02T15:04:05", ecs["ExpiredAt"].(string)[0:strings.Index(ecs["ExpiredAt"].(string), "+")])
			y := et.Year() - ct.Year()
			m := et.Month() - ct.Month()
			ecs["Period"] = y*12 + int(m)
		}

		wg.Add(4)
		instanceId := ecs["InstanceId"]
		//read instance type
		go func() {
			defer func() {
				if _err := recover(); _err != nil {
					logger.Debug(logger.ReqFormat, "DescribeInstancesType", _err)
				}
				wg.Done()
			}()
			temp := map[string]interface{}{
				"InstanceTypeId": ecs["InstanceTypeId"],
			}
			_, typeErr = s.readInstanceTypes([]interface{}{temp})
			if typeErr != nil {
				return
			}
			syncMap.Store("GpuDevices", temp["GpuDevices"])
			syncMap.Store("IsGpu", temp["IsGpu"])
		}()
		//read ebs data
		go func() {
			defer func() {
				if _err := recover(); _err != nil {
					logger.Debug(logger.ReqFormat, "DescribeVolumes", _err)
				}
				wg.Done()
			}()
			temp := map[string]interface{}{
				"InstanceId": ecs["InstanceId"],
			}
			_, ebsErr = s.readEbsVolumes([]interface{}{temp})
			if ebsErr != nil {
				return
			}
			syncMap.Store("Volumes", temp["Volumes"])
		}()
		//read user_data
		go func() {
			defer func() {
				if _err := recover(); _err != nil {
					logger.Debug(logger.ReqFormat, "DescribeUserData", _err)
				}
				bp.Release()
				wg.Done()
			}()
			bp.Acquire()
			var (
				userDataParam *map[string]interface{}
				userDataResp  *map[string]interface{}
				userData      interface{}
			)
			userDataParam = &map[string]interface{}{
				"InstanceId": instanceId,
			}
			userDataResp, userDataErr = s.Client.EcsClient.DescribeUserDataCommon(userDataParam)
			if userDataErr != nil {
				return
			}
			userData, userDataErr = bp.ObtainSdkValue("Result.UserData", *userDataResp)
			if userDataErr != nil {
				return
			}
			syncMap.Store("UserData", userData)
		}()
		//read network_interfaces
		go func() {
			defer func() {
				if _err := recover(); _err != nil {
					logger.Debug(logger.ReqFormat, "DescribeNetworkInterfaces", _err)
				}
				bp.Release()
				wg.Done()
			}()
			bp.Acquire()
			var (
				networkInterfaceParam *map[string]interface{}
				networkInterfaceResp  *map[string]interface{}
				networkInterface      interface{}
			)
			networkInterfaceParam = &map[string]interface{}{
				"InstanceId": instanceId,
			}
			networkInterfaceResp, networkInterfaceErr = s.Client.VpcClient.DescribeNetworkInterfacesCommon(networkInterfaceParam)
			if networkInterfaceErr != nil {
				return
			}
			networkInterface, networkInterfaceErr = bp.ObtainSdkValue("Result.NetworkInterfaceSets", *networkInterfaceResp)
			if networkInterfaceErr != nil {
				return
			}
			syncMap.Store("NetworkInterfaces", networkInterface)
		}()
		wg.Wait()
		//error processed
		if ebsErr != nil {
			errorStr = errorStr + ebsErr.Error() + ";"
		}
		if userDataErr != nil {
			errorStr = errorStr + userDataErr.Error() + ";"
		}
		if networkInterfaceErr != nil {
			errorStr = errorStr + networkInterfaceErr.Error() + ";"
		}
		if len(errorStr) > 0 {
			return ecs, s.CommonResponseConvert(), fmt.Errorf(errorStr)
		}
		//clean something
		delete(ecs, "Volumes")
		delete(ecs, "UserData")
		delete(ecs, "NetworkInterfaces")
		//merge extra data
		syncMap.Range(func(key, value interface{}) bool {
			ecs[key.(string)] = value
			return true
		})

		//split primary vif and secondary vif
		if networkInterfaces, ok1 := ecs["NetworkInterfaces"].([]interface{}); ok1 {
			var dataNetworkInterfaces []interface{}
			for _, vif := range networkInterfaces {
				if v1, ok2 := vif.(map[string]interface{}); ok2 {
					if v1["Type"] == "primary" {
						ecs["SubnetId"] = v1["SubnetId"]
						ecs["SecurityGroupIds"] = v1["SecurityGroupIds"]
						ecs["NetworkInterfaceId"] = v1["NetworkInterfaceId"]
						ecs["PrimaryIpAddress"] = v1["PrimaryIpAddress"]
					} else {
						dataNetworkInterfaces = append(dataNetworkInterfaces, vif)
					}
				}
			}
			if len(dataNetworkInterfaces) > 0 {
				ecs["SecondaryNetworkInterfaces"] = dataNetworkInterfaces
			}
		}

		//split System volume and Data volumes
		if volumes, ok1 := ecs["Volumes"].([]interface{}); ok1 {
			var dataVolumes []interface{}
			for _, volume := range volumes {
				if v1, ok2 := volume.(map[string]interface{}); ok2 {
					if v1["Kind"] == "system" {
						ecs["SystemVolumeType"] = v1["VolumeType"]
						ecs["SystemVolumeSize"] = v1["Size"]
						ecs["SystemVolumeId"] = v1["VolumeId"]
					} else {
						dataVolumes = append(dataVolumes, volume)
					}
				}
			}
			if len(dataVolumes) > 0 {
				v1 := volumeInfo{
					list: dataVolumes,
				}
				sort.Sort(&v1)
				ecs["DataVolumes"] = v1.list
			}
		}
		return ecs, s.CommonResponseConvert(), nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackEcsService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "RunInstances",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"zone_id": {
					ConvertType: bp.ConvertDefault,
					ForceGet:    true,
				},
				"system_volume_type": {
					ConvertType: bp.ConvertDefault,
					TargetField: "Volumes.1.VolumeType",
				},
				"system_volume_size": {
					ConvertType: bp.ConvertDefault,
					TargetField: "Volumes.1.Size",
				},
				"subnet_id": {
					ConvertType: bp.ConvertDefault,
					TargetField: "NetworkInterfaces.1.SubnetId",
				},
				"security_group_ids": {
					ConvertType: bp.ConvertWithN,
					TargetField: "NetworkInterfaces.1.SecurityGroupIds",
				},
				"data_volumes": {
					ConvertType: bp.ConvertListN,
					TargetField: "Volumes",
					StartIndex:  1,
					NextLevelConvert: map[string]bp.RequestConvert{
						"delete_with_instance": {
							TargetField: "DeleteWithInstance",
							ForceGet:    true,
						},
					},
				},
				"cpu_options": {
					ConvertType: bp.ConvertListUnique,
					TargetField: "CpuOptions",
					NextLevelConvert: map[string]bp.RequestConvert{
						"threads_per_core": {
							TargetField: "ThreadsPerCore",
						},
					},
				},
				"secondary_network_interfaces": {
					ConvertType: bp.ConvertListN,
					TargetField: "NetworkInterfaces",
					NextLevelConvert: map[string]bp.RequestConvert{
						"security_group_ids": {
							ConvertType: bp.ConvertWithN,
						},
					},
					StartIndex: 1,
				},
				"user_data": {
					ConvertType: bp.ConvertDefault,
					TargetField: "UserData",
					Convert: func(data *schema.ResourceData, i interface{}) interface{} {
						_, base64DecodeError := base64.StdEncoding.DecodeString(i.(string))
						if base64DecodeError == nil {
							return i.(string)
						} else {
							return base64.StdEncoding.EncodeToString([]byte(i.(string)))
						}
					},
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertListN,
				},
				"ipv6_address_count": {
					Ignore: true,
				},
				"ipv6_addresses": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				(*call.SdkParam)["Volumes.1.DeleteWithInstance"] = true
				(*call.SdkParam)["Count"] = 1

				if _, ok := (*call.SdkParam)["ZoneId"]; !ok || (*call.SdkParam)["ZoneId"] == "" {
					var (
						vnet map[string]interface{}
						err  error
						zone interface{}
					)
					vnet, err = s.SubnetService.ReadResource(d, (*call.SdkParam)["NetworkInterfaces.1.SubnetId"].(string))
					if err != nil {
						return false, err
					}
					zone, err = bp.ObtainSdkValue("ZoneId", vnet)
					if err != nil {
						return false, err
					}
					(*call.SdkParam)["ZoneId"] = zone
				}

				if (*call.SdkParam)["InstanceChargeType"] == "PrePaid" {
					if (*call.SdkParam)["Period"] == nil || (*call.SdkParam)["Period"].(int) < 1 {
						return false, fmt.Errorf("Instance Charge Type is PrePaid.Must set Period more than 1. ")
					}
					(*call.SdkParam)["PeriodUnit"] = "Month"
				}
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//创建ECS
				return s.Client.EcsClient.RunInstancesCommon(call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := bp.ObtainSdkValue("Result.InstanceIds.0", *resp)
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, id)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"RUNNING"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, callback)

	// 分配Ipv6
	ipv6Callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AssignIpv6Addresses",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ipv6AddressCount, ok1 := d.GetOk("ipv6_address_count")
				ipv6Addresses, ok2 := d.GetOk("ipv6_addresses")
				if !ok1 && !ok2 {
					return false, nil
				}

				var (
					networkInterfaceId string
					networkInterfaces  []interface{}
					ok                 bool
				)
				ecsInstance, err := s.ReadResource(resourceData, d.Id())
				if err != nil {
					return false, err
				}
				// query primary network interface info of the ecs instance
				if networkInterfaces, ok = ecsInstance["NetworkInterfaces"].([]interface{}); !ok {
					return false, errors.New("Instances.NetworkInterfaces is not Slice")
				}
				for _, networkInterface := range networkInterfaces {
					if networkInterfaceMap, ok := networkInterface.(map[string]interface{}); ok &&
						networkInterfaceMap["Type"] == "primary" {
						networkInterfaceId = networkInterfaceMap["NetworkInterfaceId"].(string)
					}
				}

				(*call.SdkParam)["NetworkInterfaceId"] = networkInterfaceId
				if ok1 {
					(*call.SdkParam)["Ipv6AddressCount"] = ipv6AddressCount.(int)
				} else if ok2 {
					for index, ipv6Address := range ipv6Addresses.(*schema.Set).List() {
						(*call.SdkParam)["Ipv6Address."+strconv.Itoa(index)] = ipv6Address
					}
				}

				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//分配Ipv6地址
				return s.Client.UniversalClient.DoCall(getVpcUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"RUNNING"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, ipv6Callback)

	return callbacks
}

func (s *VestackEcsService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) (callbacks []bp.Callback) {
	var (
		passwordChange bool
		flag           bool
	)

	if resourceData.HasChange("password") && !resourceData.HasChange("image_id") {
		passwordChange = true
	}

	modifyInstanceAttribute := bp.Callback{
		Call: bp.SdkCall{
			Action:         "ModifyInstanceAttribute",
			ConvertMode:    bp.RequestConvertInConvert,
			RequestIdField: "InstanceId",
			Convert: map[string]bp.RequestConvert{
				"password": {
					ConvertType: bp.ConvertDefault,
				},
				"user_data": {
					ConvertType: bp.ConvertDefault,
				},
				"instance_name": {
					ConvertType: bp.ConvertDefault,
				},
				"description": {
					ConvertType: bp.ConvertDefault,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				//if image changed ,password change in replaceSystemVolume,not here
				if _, ok := (*call.SdkParam)["Password"]; ok && d.HasChange("image_id") {
					delete(*call.SdkParam, "Password")
				}
				if len(*call.SdkParam) > 1 {
					delete(*call.SdkParam, "Tags")
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//修改实例属性
				return s.Client.EcsClient.ModifyInstanceAttributeCommon(call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"RUNNING", "STOPPED"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	callbacks = append(callbacks, modifyInstanceAttribute)

	modifyInstanceChargeType := bp.Callback{
		Call: bp.SdkCall{
			Action:         "ModifyInstanceChargeType",
			ConvertMode:    bp.RequestConvertInConvert,
			RequestIdField: "InstanceIds.1",
			Convert: map[string]bp.RequestConvert{
				"instance_charge_type": {
					ConvertType: bp.ConvertDefault,
				},
				"include_data_volumes": {
					ConvertType: bp.ConvertDefault,
					ForceGet:    true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) > 2 {
					(*call.SdkParam)["AutoPay"] = true
					if (*call.SdkParam)["InstanceChargeType"].(string) == "PostPaid" {
						//后付费
						return true, nil
					} else {
						//预付费
						period := d.Get("period")
						if period.(int) <= 0 {
							return false, fmt.Errorf("period must set and more than 0 ")
						}
						(*call.SdkParam)["Period"] = period
						//(*call.SdkParam)["PeriodUnit"] = d.Get("period_unit")
						(*call.SdkParam)["PeriodUnit"] = "Month"
						return true, nil
					}
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//修改实例计费方式
				return s.Client.EcsClient.ModifyInstanceChargeTypeCommon(call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"RUNNING", "STOPPED"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}

	callbacks = append(callbacks, modifyInstanceChargeType)

	//primary vif sg change
	if resourceData.HasChange("security_group_ids") {
		modifyNetworkInterfaceAttributes := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyNetworkInterfaceAttributes",
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"network_interface_id": {
						ConvertType: bp.ConvertDefault,
						ForceGet:    true,
					},
					"security_group_ids": {
						ConvertType: bp.ConvertWithN,
						ForceGet:    true,
					},
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					return s.Client.VpcClient.ModifyNetworkInterfaceAttributesCommon(call.SdkParam)
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					return nil
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"RUNNING", "STOPPED"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, modifyNetworkInterfaceAttributes)
	}
	//system_volume change
	extendVolume := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ExtendVolume",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"system_volume_id": {
					ConvertType: bp.ConvertDefault,
					TargetField: "VolumeId",
					ForceGet:    true,
				},
				"system_volume_size": {
					ConvertType: bp.ConvertDefault,
					TargetField: "NewSize",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) == 2 {
					o, n := d.GetChange("system_volume_size")
					if o.(int) > n.(int) {
						return false, fmt.Errorf("SystemVolumeSize only support extend. ")
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.EbsClient.ExtendVolumeCommon(call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"RUNNING", "STOPPED"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	callbacks = append(callbacks, extendVolume)

	if !resourceData.HasChange("instance_charge_type") && resourceData.Get("instance_charge_type").(string) == "PrePaid" {
		//只有当没执行实例状态变更才生效并且是预付费
		renewInstance := bp.Callback{
			Call: bp.SdkCall{
				Action:         "RenewInstance",
				ConvertMode:    bp.RequestConvertInConvert,
				RequestIdField: "InstanceId",
				Convert: map[string]bp.RequestConvert{
					"period": {
						ConvertType: bp.ConvertDefault,
						Convert: func(data *schema.ResourceData, i interface{}) interface{} {
							o, n := data.GetChange("period")
							return n.(int) - o.(int)
						},
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) > 1 {
						if (*call.SdkParam)["Period"].(int) <= 0 {
							return false, fmt.Errorf("period must set and more than 0 ")
						}
						//(*call.SdkParam)["PeriodUnit"] = d.Get("period_unit")
						(*call.SdkParam)["PeriodUnit"] = "Month"
						return true, nil
					}
					return false, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					//续费实例
					return s.Client.EcsClient.RenewInstanceCommon(call.SdkParam)
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					return nil
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"RUNNING", "STOPPED"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, renewInstance)
	}
	//only password changed need stop
	if passwordChange {
		stopInstance := s.StartOrStopInstanceCallback(resourceData, true, &flag)
		callbacks = append(callbacks, stopInstance)
	}
	//instance_type
	if resourceData.HasChange("instance_type") {
		//need stop before ModifyInstanceSpec

		stopInstance := s.StartOrStopInstanceCallback(resourceData, true, &flag)
		callbacks = append(callbacks, stopInstance)

		modifyInstanceSpec := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyInstanceSpec",
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"instance_type": {
						ConvertType: bp.ConvertDefault,
					},
				},
				RequestIdField: "InstanceId",
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) > 1 {
						return true, nil
					}
					return false, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					//修改实例规格
					return s.Client.EcsClient.ModifyInstanceSpecCommon(call.SdkParam)
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					return nil
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"RUNNING", "STOPPED"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, modifyInstanceSpec)
	}
	//image change
	if resourceData.HasChange("image_id") {
		//need stop before ReplaceSystemVolume
		stopInstance := s.StartOrStopInstanceCallback(resourceData, true, &flag)
		callbacks = append(callbacks, stopInstance)
		replaceSystemVolume := bp.Callback{
			Call: bp.SdkCall{
				Action:         "ReplaceSystemVolume",
				ConvertMode:    bp.RequestConvertInConvert,
				RequestIdField: "InstanceId",
				Convert: map[string]bp.RequestConvert{
					"image_id": {
						ConvertType: bp.ConvertDefault,
					},
					"system_volume_size": {
						ConvertType: bp.ConvertDefault,
						ForceGet:    true,
					},
					"key_pair_name": {
						ConvertType: bp.ConvertDefault,
						ForceGet:    true,
					},
					"password": {
						ConvertType: bp.ConvertDefault,
						ForceGet:    true,
					},
					"user_data": {
						ConvertType: bp.ConvertDefault,
						ForceGet:    true,
					},
					"keep_image_credential": {
						ConvertType: bp.ConvertDefault,
						ForceGet:    true,
					},
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					return s.Client.EcsClient.ReplaceSystemVolumeCommon(call.SdkParam)
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					return nil
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"RUNNING"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, replaceSystemVolume)
	}

	if resourceData.HasChange("deployment_set_id") {
		stopInstance := s.StartOrStopInstanceCallback(resourceData, true, &flag)
		callbacks = append(callbacks, stopInstance)
		deploymentSet := bp.Callback{
			Call: bp.SdkCall{
				Action:         "ModifyInstanceDeployment",
				ConvertMode:    bp.RequestConvertInConvert,
				Convert:        map[string]bp.RequestConvert{},
				RequestIdField: "InstanceId",
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["DeploymentSetId"] = resourceData.Get("deployment_set_id")
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					return nil
				},
			},
		}
		refresh := map[bp.ResourceService]*bp.StateRefresh{
			ecs_deployment_set_associate.NewEcsDeploymentSetAssociateService(s.Client): {
				Target:     []string{"success"},
				Timeout:    resourceData.Timeout(schema.TimeoutCreate),
				ResourceId: resourceData.Get("deployment_set_id").(string) + ":" + resourceData.Id(),
			},
		}

		if resourceData.Get("deployment_set_id").(string) != "" {
			deploymentSet.Call.ExtraRefresh = refresh
		}

		callbacks = append(callbacks, deploymentSet)
	}

	startInstance := s.StartOrStopInstanceCallback(resourceData, false, &flag)
	callbacks = append(callbacks, startInstance)

	// 更新Tags
	setResourceTagsCallbacks := bp.SetResourceTags(s.Client, "CreateTags", "DeleteTags", "instance", resourceData, getUniversalInfo)
	callbacks = append(callbacks, setResourceTagsCallbacks...)

	return callbacks
}

func (s *VestackEcsService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteInstance",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"InstanceId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//删除ECS
				return s.Client.EcsClient.DeleteInstanceCommon(call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					ecs, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading ecs on delete %q, %w", d.Id(), callErr))
						}
					}

					//if ecs["InstanceChargeType"] == "PrePaid" {
					//	return resource.NonRetryableError(fmt.Errorf("PrePaid instance charge type not support remove,Please change instance charge type to PostPaid. "))
					//}

					if ecs["InstanceChargeType"] == "PrePaid" {
						logger.Debug(logger.RespFormat, call.Action, "PrePaid instance charge type not support remove,Only Remove from State")
						return nil
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

func (s *VestackEcsService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "InstanceIds",
				ConvertType: bp.ConvertWithN,
			},
			"tags": {
				TargetField: "TagFilters",
				ConvertType: bp.ConvertListN,
				NextLevelConvert: map[string]bp.RequestConvert{
					"value": {
						TargetField: "Values.1",
					},
				},
			},
			"deployment_set_ids": {
				TargetField: "DeploymentSetIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:        "InstanceName",
		IdField:          "InstanceId",
		CollectField:     "instances",
		ResponseConverts: s.CommonResponseConvert(),
		ExtraData: func(sourceData []interface{}) (extraData []interface{}, err error) {
			sourceData, err = s.readInstanceTypes(sourceData)
			if err != nil {
				return extraData, err
			}
			sourceData, err = s.readEbsVolumes(sourceData)
			if err != nil {
				return extraData, err
			}
			return sourceData, err
		},
	}
}

func (s *VestackEcsService) CommonResponseConvert() map[string]bp.ResponseConvert {
	return map[string]bp.ResponseConvert{
		"Id": {
			TargetField: "instance_id",
		},
		"Hostname": {
			TargetField: "host_name",
		},
		"InstanceTypeId": {
			TargetField: "instance_type",
		},
		"InstanceType": {
			Ignore: true,
		},
		"SystemVolumeSize": {
			TargetField: "system_volume_size",
			Convert: func(i interface{}) interface{} {
				size, _ := strconv.Atoi(i.(string))
				return size
			},
		},
		"UserData": {
			TargetField: "user_data",
			Convert: func(i interface{}) interface{} {
				v, base64DecodeError := base64.StdEncoding.DecodeString(i.(string))
				if base64DecodeError != nil {
					v = []byte(i.(string))
				}
				return string(v)
			},
		},
		"DataVolumes": {
			TargetField: "data_volumes",
			Convert: func(i interface{}) interface{} {
				var results []interface{}
				if dd, ok := i.([]interface{}); ok {
					for _, _data := range dd {
						if v, ok1 := _data.(map[string]interface{}); ok1 {
							if reflect.TypeOf(v["Size"]).Kind() == reflect.String {
								v["Size"], _ = strconv.Atoi(v["Size"].(string))
							}
							results = append(results, v)
						}
					}
				}
				return results
			},
		},
		"Volumes": {
			TargetField: "volumes",
			Convert: func(i interface{}) interface{} {
				var results []interface{}
				if dd, ok := i.([]interface{}); ok {
					for _, _data := range dd {
						if v, ok1 := _data.(map[string]interface{}); ok1 {
							if reflect.TypeOf(v["Size"]).Kind() == reflect.String {
								v["Size"], _ = strconv.Atoi(v["Size"].(string))
							}
							results = append(results, v)
						}
					}
				}
				return results
			},
		},
		"GpuDevices": {
			TargetField: "gpu_devices",
			Convert: func(i interface{}) interface{} {
				var results []interface{}
				if dd, ok := i.([]interface{}); ok {
					for _, _data := range dd {
						if v, ok1 := _data.(map[string]interface{}); ok1 {
							memorySize, _ := bp.ObtainSdkValue("Memory.Size", v)
							encryptedMemorySize, _ := bp.ObtainSdkValue("Memory.EncryptedSize", v)
							delete(v, "Memory")
							v["MemorySize"] = memorySize
							v["EncryptedMemorySize"] = encryptedMemorySize
							results = append(results, v)
						}
					}
				}
				return results
			},
		},
	}
}

func (s *VestackEcsService) StartOrStopInstanceCallback(resourceData *schema.ResourceData, isStop bool, flag *bool) bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"InstanceId": resourceData.Id(),
			},
		},
	}
	if isStop {
		callback.Call.Action = "StopInstance"
		callback.Call.BeforeCall = func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
			instance, err := s.ReadResource(resourceData, resourceData.Id())
			if err != nil {
				return false, err
			}
			status, err := bp.ObtainSdkValue("Status", instance)
			if err != nil {
				return false, err
			}
			if status.(string) == "RUNNING" {
				return true, nil
			}
			return false, nil
		}
		callback.Call.ExecuteCall = func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
			logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
			return s.Client.EcsClient.StopInstanceCommon(call.SdkParam)
		}
		callback.Call.AfterCall = func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
			*flag = true
			return nil
		}
		callback.Call.Refresh = &bp.StateRefresh{
			Target:  []string{"STOPPED"},
			Timeout: resourceData.Timeout(schema.TimeoutUpdate),
		}
	} else {
		callback.Call.Action = "StartInstance"
		callback.Call.BeforeCall = func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
			instance, err := s.ReadResource(resourceData, resourceData.Id())
			if err != nil {
				return false, err
			}
			status, err := bp.ObtainSdkValue("Status", instance)
			if err != nil {
				return false, err
			}
			if status.(string) == "RUNNING" {
				return false, nil
			}
			return *flag, nil
		}
		callback.Call.ExecuteCall = func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
			logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
			return s.Client.EcsClient.StartInstanceCommon(call.SdkParam)
		}
		callback.Call.Refresh = &bp.StateRefresh{
			Target:  []string{"RUNNING"},
			Timeout: resourceData.Timeout(schema.TimeoutUpdate),
		}
	}
	return callback
}

func (s *VestackEcsService) ReadResourceId(id string) string {
	return id
}

func (s *VestackEcsService) readInstanceTypes(sourceData []interface{}) (extraData []interface{}, err error) {
	//merge instance_type_info
	var (
		wg      sync.WaitGroup
		syncMap sync.Map
	)
	if len(sourceData) == 0 {
		return sourceData, err
	}
	wg.Add(len(sourceData))
	for _, data := range sourceData {
		instance := data
		var (
			instanceTypeId interface{}
			action         string
			resp           *map[string]interface{}
			results        interface{}
			_err           error
		)
		go func() {
			defer func() {
				if e := recover(); e != nil {
					logger.Debug(logger.ReqFormat, action, e)
				}
				bp.Release()
				wg.Done()
			}()
			bp.Acquire()

			instanceTypeId, _err = bp.ObtainSdkValue("InstanceTypeId", instance)
			if _err != nil {
				syncMap.Store(instanceTypeId, err)
				return
			}
			//if exist continue
			if _, ok := syncMap.Load(instanceTypeId); ok {
				return
			}

			action = "DescribeInstanceTypes"
			logger.Debug(logger.ReqFormat, action, instanceTypeId)
			instanceTypeCondition := map[string]interface{}{
				"InstanceTypeIds.1": instanceTypeId,
			}
			logger.Debug(logger.ReqFormat, action, instanceTypeCondition)
			resp, _err = s.Client.EcsClient.DescribeInstanceTypesCommon(&instanceTypeCondition)
			if _err != nil {
				syncMap.Store(instanceTypeId, err)
				return
			}
			logger.Debug(logger.RespFormat, action, instanceTypeCondition, *resp)
			results, _err = bp.ObtainSdkValue("Result.InstanceTypes.0", *resp)
			if _err != nil {
				syncMap.Store(instanceTypeId, err)
				return
			}
			syncMap.Store(instanceTypeId, results)
		}()
	}
	wg.Wait()
	var errorStr string
	for _, instance := range sourceData {
		var (
			instanceTypeId interface{}
			gpu            interface{}
			gpuDevices     interface{}
		)
		instanceTypeId, err = bp.ObtainSdkValue("InstanceTypeId", instance)
		if err != nil {
			return
		}
		if v, ok := syncMap.Load(instanceTypeId); ok {
			if e1, ok1 := v.(error); ok1 {
				errorStr = errorStr + e1.Error() + ";"
			}
			gpu, _ = bp.ObtainSdkValue("Gpu", v)
			if gpu != nil {
				gpuDevices, _ = bp.ObtainSdkValue("Gpu.GpuDevices", v)
				instance.(map[string]interface{})["GpuDevices"] = gpuDevices
				instance.(map[string]interface{})["IsGpu"] = true
			} else {
				instance.(map[string]interface{})["GpuDevices"] = []interface{}{}
				instance.(map[string]interface{})["IsGpu"] = false
			}
		}
		extraData = append(extraData, instance)
	}
	if len(errorStr) > 0 {
		return extraData, fmt.Errorf(errorStr)
	}
	return extraData, err
}

func (s *VestackEcsService) readEbsVolumes(sourceData []interface{}) (extraData []interface{}, err error) {
	//merge ebs
	var (
		wg      sync.WaitGroup
		syncMap sync.Map
	)
	if len(sourceData) == 0 {
		return sourceData, err
	}
	wg.Add(len(sourceData))
	for _, data := range sourceData {
		instance := data
		var (
			instanceId interface{}
			action     string
			resp       *map[string]interface{}
			results    interface{}
			_err       error
		)
		go func() {
			defer func() {
				if e := recover(); e != nil {
					logger.Debug(logger.ReqFormat, action, e)
				}
				bp.Release()
				wg.Done()
			}()
			bp.Acquire()

			instanceId, _err = bp.ObtainSdkValue("InstanceId", instance)
			if _err != nil {
				syncMap.Store(instanceId, err)
				return
			}
			action = "DescribeVolumes"
			logger.Debug(logger.ReqFormat, action, instanceId)
			volumeCondition := map[string]interface{}{
				"InstanceId": instanceId,
			}
			logger.Debug(logger.ReqFormat, action, volumeCondition)
			resp, _err = s.Client.EbsClient.DescribeVolumesCommon(&volumeCondition)
			if _err != nil {
				syncMap.Store(instanceId, err)
				return
			}
			logger.Debug(logger.RespFormat, action, volumeCondition, *resp)
			results, _err = bp.ObtainSdkValue("Result.Volumes", *resp)
			if _err != nil {
				syncMap.Store(instanceId, err)
				return
			}
			syncMap.Store(instanceId, results)
		}()
	}
	wg.Wait()
	var errorStr string
	for _, instance := range sourceData {
		var (
			instanceId interface{}
		)
		instanceId, err = bp.ObtainSdkValue("InstanceId", instance)
		if err != nil {
			return
		}
		if v, ok := syncMap.Load(instanceId); ok {
			if e1, ok1 := v.(error); ok1 {
				errorStr = errorStr + e1.Error() + ";"
			}
			instance.(map[string]interface{})["Volumes"] = v
		}
		extraData = append(extraData, instance)
	}
	if len(errorStr) > 0 {
		return extraData, fmt.Errorf(errorStr)
	}
	return extraData, err
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
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
		Action:      actionName,
	}
}

type volumeInfo struct {
	list []interface{}
}

func (v *volumeInfo) Len() int {
	return len(v.list)
}

func (v *volumeInfo) Less(i, j int) bool {
	return v.list[i].(map[string]interface{})["VolumeName"].(string) < v.list[j].(map[string]interface{})["VolumeName"].(string)
}

func (v *volumeInfo) Swap(i, j int) {
	v.list[i], v.list[j] = v.list[j], v.list[i]
}

func (s *VestackEcsService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "ecs",
		ResourceType:         "instance",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}

func (s *VestackEcsService) UnsubscribeInfo(resourceData *schema.ResourceData, resource *schema.Resource) (*bp.UnsubscribeInfo, error) {
	info := bp.UnsubscribeInfo{
		InstanceId: s.ReadResourceId(resourceData.Id()),
	}
	if resourceData.Get("instance_charge_type") == "PrePaid" {
		//查询实例类型的配置
		action := "DescribeInstanceTypes"
		input := map[string]interface{}{
			"InstanceTypeIds.1": resourceData.Get("instance_type"),
		}
		var (
			output *map[string]interface{}
			err    error
			t      interface{}
		)
		output, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &input)
		if err != nil {
			return &info, err
		}
		t, err = bp.ObtainSdkValue("Result.InstanceTypes.0", *output)
		if err != nil {
			return &info, err
		}
		if tt, ok := t.(map[string]interface{}); ok {
			if tt["Gpu"] != nil && tt["Rdma"] != nil {
				info.Products = []string{"HPC_GPU", "ECS", "ECS_BareMetal", "GPU_Server"}
			} else if tt["Gpu"] != nil && tt["Rdma"] == nil {
				info.Products = []string{"GPU_Server", "ECS", "ECS_BareMetal", "HPC_GPU"}
			} else {
				info.Products = []string{"ECS", "ECS_BareMetal", "GPU_Server", "HPC_GPU"}
			}
			info.NeedUnsubscribe = true
		}

	}
	return &info, nil
}
