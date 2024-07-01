package bucket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackTosBucketService struct {
	Client *bp.SdkClient
}

func NewTosBucketService(c *bp.SdkClient) *VestackTosBucketService {
	return &VestackTosBucketService{
		Client: c,
	}
}

func (s *VestackTosBucketService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackTosBucketService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	tos := s.Client.BypassSvcClient
	var (
		action  string
		resp    *map[string]interface{}
		results interface{}
	)
	action = "ListBuckets"
	logger.Debug(logger.ReqFormat, action, nil)
	resp, err = tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
	}, nil)
	if err != nil {
		return data, err
	}
	results, err = bp.ObtainSdkValue(bp.BypassResponse+".Buckets", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VestackTosBucketService) ReadResource(resourceData *schema.ResourceData, instanceId string) (data map[string]interface{}, err error) {
	tos := s.Client.BypassSvcClient
	var (
		action  string
		resp    *map[string]interface{}
		ok      bool
		header  http.Header
		acl     map[string]interface{}
		version map[string]interface{}
		buckets []interface{}
	)

	if instanceId == "" {
		instanceId = s.ReadResourceId(resourceData.Id())
	} else {
		instanceId = s.ReadResourceId(instanceId)
	}

	action = "HeadBucket"
	logger.Debug(logger.ReqFormat, action, instanceId)
	resp, err = tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.HEAD,
		Domain:     instanceId,
	}, nil)
	logger.Debug(logger.ReqFormat, action, *resp)
	logger.Debug(logger.ReqFormat, action, err)
	if err != nil {
		return data, err
	}

	buckets, err = s.ReadResources(nil)
	if err != nil {
		return data, err
	}
	var (
		local interface{}
		name  interface{}
	)
	for _, bucket := range buckets {
		local, err = bp.ObtainSdkValue("Location", bucket)
		if err != nil {
			return data, err
		}
		name, err = bp.ObtainSdkValue("Name", bucket)
		if err != nil {
			return data, err
		}
		if local.(string) == s.Client.Region && name.(string) == instanceId {
			data = bucket.(map[string]interface{})
		}
	}
	if data == nil {
		data = make(map[string]interface{})
	}

	if header, ok = (*resp)[bp.BypassHeader].(http.Header); ok {
		if header.Get("X-Tos-Storage-Class") != "" {
			data["StorageClass"] = header.Get("X-Tos-Storage-Class")
		}
	}

	action = "GetBucketAcl"
	req := map[string]interface{}{
		"acl": "",
	}
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     instanceId,
	}, &req)
	if err != nil {
		return data, err
	}
	if acl, ok = (*resp)[bp.BypassResponse].(map[string]interface{}); ok {
		data["PublicAcl"] = acl
		data["AccountAcl"] = acl
	}

	action = "GetBucketVersioning"
	req = map[string]interface{}{
		"versioning": "",
	}
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     instanceId,
	}, &req)
	if err != nil {
		return data, err
	}
	if version, ok = (*resp)[bp.BypassResponse].(map[string]interface{}); ok {
		data["EnableVersion"] = version
	}

	if len(data) == 0 {
		return data, fmt.Errorf("bucket %s not exist ", instanceId)
	}
	return data, nil
}

func (s *VestackTosBucketService) RefreshResourceState(data *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      60 * time.Second,
		MinTimeout: 60 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			return data, "Success", err
		},
	}
}

//func (s *VestackTosBucketService) getIdPermission(p string, grants []interface{}) []interface{} {
//	var result []interface{}
//	for _, grant := range grants {
//		permission, _ := bp.ObtainSdkValue("Permission", grant)
//		id, _ := bp.ObtainSdkValue("Grantee.ID", grant)
//		t, _ := bp.ObtainSdkValue("Grantee.Type", grant)
//		if id != nil && t.(string) == "CanonicalUser" && p == permission.(string) {
//			result = append(result, "Id="+id.(string))
//		}
//	}
//	return result
//}

func (s *VestackTosBucketService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return m, map[string]bp.ResponseConvert{
			"EnableVersion": {
				Convert: func(i interface{}) interface{} {
					status, _ := bp.ObtainSdkValue("Status", i)
					return status.(string) == "Enabled"
				},
			},
			"AccountAcl": {
				Convert: bp.ConvertTosAccountAcl(),
			},
			"PublicAcl": {
				Convert: bp.ConvertTosPublicAcl(),
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *VestackTosBucketService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	//create bucket
	callback := bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "CreateBucket",
			ConvertMode:     bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"bucket_name": {
					ConvertType: bp.ConvertDefault,
					TargetField: "BucketName",
					SpecialParam: &bp.SpecialParam{
						Type: bp.DomainParam,
					},
				},
				"public_acl": {
					ConvertType: bp.ConvertDefault,
					TargetField: "x-tos-acl",
					SpecialParam: &bp.SpecialParam{
						Type: bp.HeaderParam,
					},
				},
				"storage_class": {
					ConvertType: bp.ConvertDefault,
					TargetField: "x-tos-storage-class",
					SpecialParam: &bp.SpecialParam{
						Type: bp.HeaderParam,
					},
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//创建Bucket
				return s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					HttpMethod: bp.PUT,
					Domain:     (*call.SdkParam)[bp.BypassDomain].(string),
					Header:     (*call.SdkParam)[bp.BypassHeader].(map[string]string),
				}, nil)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId((*call.SdkParam)[bp.BypassDomain].(string))
				return nil
			},
		},
	}
	//version
	callbackVersion := bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "PutBucketVersioning",
			ConvertMode:     bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"bucket_name": {
					ConvertType: bp.ConvertDefault,
					TargetField: "BucketName",
					SpecialParam: &bp.SpecialParam{
						Type: bp.DomainParam,
					},
				},
				"enable_version": {
					ConvertType: bp.ConvertDefault,
					TargetField: "Status",
					Convert: func(data *schema.ResourceData, i interface{}) interface{} {
						if i.(bool) {
							return "Enabled"
						} else {
							return ""
						}
					},
					ForceGet: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				//if disable version,skip this call
				if (*call.SdkParam)[bp.BypassParam].(map[string]interface{})["Status"] == "" {
					return false, nil
				}
				return true, nil
			},
			ExecuteCall: s.executePutBucketVersioning(),
		},
	}
	//acl
	callbackAcl := bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "PutBucketAcl",
			ConvertMode:     bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"bucket_name": {
					ConvertType: bp.ConvertDefault,
					TargetField: "BucketName",
					SpecialParam: &bp.SpecialParam{
						Type: bp.DomainParam,
					},
				},
				"account_acl": {
					ConvertType: bp.ConvertListN,
					TargetField: "Grants",
					NextLevelConvert: map[string]bp.RequestConvert{
						"account_id": {
							ConvertType: bp.ConvertDefault,
							TargetField: "Grantee.ID",
						},
						"acl_type": {
							ConvertType: bp.ConvertDefault,
							TargetField: "Grantee.Type",
						},
						"permission": {
							ConvertType: bp.ConvertDefault,
							TargetField: "Permission",
						},
					},
				},
			},
			BeforeCall:  s.beforePutBucketAcl(),
			ExecuteCall: s.executePutBucketAcl(),
			//Refresh: &bp.StateRefresh{
			//	Target:  []string{"Success"},
			//	Timeout: resourceData.Timeout(schema.TimeoutCreate),
			//},
		},
	}
	return []bp.Callback{callback, callbackVersion, callbackAcl}
}

func (s *VestackTosBucketService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	if data.HasChange("enable_version") {
		//version
		callbackVersion := bp.Callback{
			Call: bp.SdkCall{
				ServiceCategory: bp.ServiceBypass,
				Action:          "PutBucketVersioning",
				ConvertMode:     bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"bucket_name": {
						ConvertType: bp.ConvertDefault,
						TargetField: "BucketName",
						SpecialParam: &bp.SpecialParam{
							Type: bp.DomainParam,
						},
						ForceGet: true,
					},
					"enable_version": {
						ConvertType: bp.ConvertDefault,
						TargetField: "Status",
						Convert: func(data *schema.ResourceData, i interface{}) interface{} {
							if i.(bool) {
								return "Enabled"
							} else {
								return "Suspended"
							}
						},
						ForceGet: true,
					},
				},
				ExecuteCall: s.executePutBucketVersioning(),
			},
		}
		callbacks = append(callbacks, callbackVersion)
	}
	var grant = []string{
		"public_acl",
		"account_acl",
	}
	for _, v := range grant {
		if data.HasChange(v) {
			callbackAcl := bp.Callback{
				Call: bp.SdkCall{
					ServiceCategory: bp.ServiceBypass,
					Action:          "PutBucketAcl",
					ConvertMode:     bp.RequestConvertInConvert,
					Convert: map[string]bp.RequestConvert{
						"bucket_name": {
							ConvertType: bp.ConvertDefault,
							TargetField: "BucketName",
							SpecialParam: &bp.SpecialParam{
								Type: bp.DomainParam,
							},
							ForceGet: true,
						},
						"account_acl": {
							ConvertType: bp.ConvertListN,
							TargetField: "Grants",
							NextLevelConvert: map[string]bp.RequestConvert{
								"account_id": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Grantee.ID",
									ForceGet:    true,
								},
								"acl_type": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Grantee.Type",
									ForceGet:    true,
								},
								"permission": {
									ConvertType: bp.ConvertDefault,
									TargetField: "Permission",
									ForceGet:    true,
								},
							},
							ForceGet: true,
						},
					},
					BeforeCall:  s.beforePutBucketAcl(),
					ExecuteCall: s.executePutBucketAcl(),
					Refresh: &bp.StateRefresh{
						Target:  []string{"Success"},
						Timeout: data.Timeout(schema.TimeoutCreate),
					},
				},
			}
			callbacks = append(callbacks, callbackAcl)
			break
		}
	}

	return callbacks
}

func (s *VestackTosBucketService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteBucket",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"BucketName": s.ReadResourceId(resourceData.Id()),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//删除Bucket
				return s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					HttpMethod: bp.DELETE,
					Domain:     (*call.SdkParam)["BucketName"].(string),
				}, nil)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading tos on delete %q, %w", s.ReadResourceId(d.Id()), callErr))
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

func (s *VestackTosBucketService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {

	name, ok := data.GetOk("bucket_name")
	return bp.DataSourceInfo{
		ServiceCategory: bp.ServiceBypass,
		RequestConverts: map[string]bp.RequestConvert{
			"bucket_name": {
				Ignore: true,
			},
		},
		NameField:        "Name",
		IdField:          "BucketId",
		CollectField:     "buckets",
		ResponseConverts: map[string]bp.ResponseConvert{},
		ExtraData: func(sourceData []interface{}) (extraData []interface{}, err error) {
			for _, v := range sourceData {
				if v.(map[string]interface{})["Location"].(string) != s.Client.Region {
					continue
				}
				if ok {
					if name.(string) == v.(map[string]interface{})["Name"].(string) {
						v.(map[string]interface{})["BucketId"] = v.(map[string]interface{})["Name"].(string)
						extraData = append(extraData, v)
						break
					} else {
						continue
					}
				} else {
					v.(map[string]interface{})["BucketId"] = v.(map[string]interface{})["Name"].(string)
					extraData = append(extraData, v)
				}

			}
			return extraData, err
		},
	}
}

func (s *VestackTosBucketService) ReadResourceId(id string) string {
	return id
}

func (s *VestackTosBucketService) beforePutBucketAcl() bp.BeforeCallFunc {

	return func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
		data, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
			HttpMethod: bp.GET,
			Domain:     (*call.SdkParam)[bp.BypassDomain].(string),
			UrlParam: map[string]string{
				"acl": "",
			},
		}, nil)
		return bp.BeforeTosPutAcl(d, call, data, err)
	}
}

func (s *VestackTosBucketService) executePutBucketAcl() bp.ExecuteCallFunc {
	return func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
		logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
		//PutAcl
		param := (*call.SdkParam)[bp.BypassParam].(map[string]interface{})
		return s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
			HttpMethod:  bp.PUT,
			ContentType: bp.ApplicationJSON,
			Domain:      (*call.SdkParam)[bp.BypassDomain].(string),
			Header:      (*call.SdkParam)[bp.BypassHeader].(map[string]string),
			UrlParam: map[string]string{
				"acl": "",
			},
		}, &param)
	}
}

func (s *VestackTosBucketService) executePutBucketVersioning() bp.ExecuteCallFunc {
	return func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
		logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
		//PutVersion
		condition := (*call.SdkParam)[bp.BypassParam].(map[string]interface{})
		return s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
			ContentType: bp.ApplicationJSON,
			HttpMethod:  bp.PUT,
			Domain:      (*call.SdkParam)[bp.BypassDomain].(string),
			UrlParam: map[string]string{
				"versioning": "",
			},
		}, &condition)
	}
}
