package object

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackTosObjectService struct {
	Client *bp.SdkClient
}

func NewTosObjectService(c *bp.SdkClient) *VestackTosObjectService {
	return &VestackTosObjectService{
		Client: c,
	}
}

func (s *VestackTosObjectService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *VestackTosObjectService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	tos := s.Client.BypassSvcClient
	var (
		action  string
		resp    *map[string]interface{}
		results interface{}
	)
	action = "ListObjects"
	logger.Debug(logger.ReqFormat, action, nil)
	resp, err = tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     condition[bp.BypassDomain].(string),
	}, nil)
	if err != nil {
		return data, err
	}
	results, err = bp.ObtainSdkValue(bp.BypassResponse+".Contents", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VestackTosObjectService) ReadResource(resourceData *schema.ResourceData, instanceId string) (data map[string]interface{}, err error) {
	tos := s.Client.BypassSvcClient
	bucketName := resourceData.Get("bucket_name").(string)
	var (
		action        string
		resp          *map[string]interface{}
		respBody      *map[string]interface{}
		ok            bool
		header        http.Header
		acl           map[string]interface{}
		bucketVersion map[string]interface{}
	)

	if instanceId == "" {
		instanceId = s.ReadResourceId(resourceData.Id())
	} else {
		instanceId = s.ReadResourceId(instanceId)
	}

	action = "HeadObject"
	logger.Debug(logger.ReqFormat, action, bucketName+":"+instanceId)
	resp, err = tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.HEAD,
		Domain:     bucketName,
		Path:       []string{instanceId},
	}, nil)
	if err != nil {
		return data, err
	}
	data = make(map[string]interface{})

	if header, ok = (*resp)[bp.BypassHeader].(http.Header); ok {
		if header.Get("X-Tos-Storage-Class") != "" {
			data["StorageClass"] = header.Get("x-tos-storage-class")
		}
		if header.Get("Content-Type") != "" {
			data["ContentType"] = header.Get("Content-Type")
			if strings.Contains(strings.ToLower(header.Get("Content-Type")), "application/json") ||
				strings.Contains(strings.ToLower(header.Get("Content-Type")), "application/xml") ||
				strings.Contains(strings.ToLower(header.Get("Content-Type")), "text/plain") {
				action = "GetObject"
				logger.Debug(logger.ReqFormat, action, bucketName+":"+instanceId)
				respBody, err = tos.DoBypassSvcCall(bp.BypassSvcInfo{
					HttpMethod: bp.GET,
					Domain:     bucketName,
					Path:       []string{instanceId},
				}, nil)
				if err != nil {
					return data, err
				}
				data["Content"] = (*respBody)[bp.BypassResponseData]
			}
		}
		if header.Get("X-Tos-Server-Side-Encryption") != "" {
			data["Encryption"] = header.Get("X-Tos-Server-Side-Encryption")
		}
		if header.Get("x-tos-meta-content-md5") != "" {
			data["ContentMd5"] = strings.Replace(header.Get("x-tos-meta-content-md5"), "\"", "", -1)
		}

		if header.Get("X-Tos-Version-Id") != "" {
			action = "ListObjects"
			logger.Debug(logger.ReqFormat, action, bucketName+":"+instanceId)

			var (
				nextVersionIdMarker string
				versionIds          []string
			)

			for {
				urlParam := map[string]string{
					"prefix":   instanceId,
					"max-keys": "100",
					"versions": "",
				}
				if nextVersionIdMarker != "" {
					urlParam["key-marker"] = instanceId
					urlParam["version-id-marker"] = nextVersionIdMarker
				}

				resp, err = tos.DoBypassSvcCall(bp.BypassSvcInfo{
					HttpMethod: bp.GET,
					Domain:     bucketName,
					UrlParam:   urlParam,
				}, nil)

				if err != nil {
					return data, err
				}
				versions, _ := bp.ObtainSdkValue(bp.BypassResponse+".Versions", *resp)
				next, _ := bp.ObtainSdkValue(bp.BypassResponse+".NextVersionIdMarker", *resp)

				if versions == nil || len(versions.([]interface{})) == 0 {
					break
				}

				if next == nil || next.(string) == "" {
					nextVersionIdMarker = ""
				} else {
					nextVersionIdMarker = next.(string)
				}

				for _, version := range versions.([]interface{}) {
					versionId, _ := bp.ObtainSdkValue("VersionId", version)
					versionIds = append(versionIds, versionId.(string))
				}

				if nextVersionIdMarker == "" {
					break
				}
			}
			logger.Debug(logger.ReqFormat, action, versionIds)
			data["VersionIds"] = versionIds
		}
	}

	action = "GetObjectAcl"
	req := map[string]interface{}{
		"acl": "",
	}
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = tos.DoBypassSvcCall(bp.BypassSvcInfo{
		HttpMethod: bp.GET,
		Domain:     bucketName,
		Path:       []string{instanceId},
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
		Domain:     bucketName,
	}, &req)
	if err != nil {
		return data, err
	}
	if bucketVersion, ok = (*resp)[bp.BypassResponse].(map[string]interface{}); ok {
		data["EnableVersion"] = bucketVersion
	}

	if len(data) == 0 {
		return data, fmt.Errorf("object %s not exist ", instanceId)
	}
	return data, nil
}

func (s *VestackTosObjectService) RefreshResourceState(data *schema.ResourceData, target []string, timeout time.Duration, instanceId string) *resource.StateChangeConf {
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

func (VestackTosObjectService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
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

func (s *VestackTosObjectService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	//create object
	callback := s.createOrReplaceObject(resourceData, resource, false)
	//acl
	callbackAcl := s.createOrUpdateObjectAcl(resourceData, resource, false)
	return []bp.Callback{callback, callbackAcl}
}

func (s *VestackTosObjectService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	if data.HasChange("file_path") || data.HasChanges("content_md5") {
		callbacks = append(callbacks, s.createOrReplaceObject(data, resource, true))
		callbacks = append(callbacks, s.createOrUpdateObjectAcl(data, resource, true))
	} else {
		var grant = []string{
			"public_acl",
			"account_acl",
		}
		for _, v := range grant {
			if data.HasChange(v) {
				callbackAcl := s.createOrUpdateObjectAcl(data, resource, true)
				callbacks = append(callbacks, callbackAcl)
				break
			}
		}
	}

	return callbacks
}

func (s *VestackTosObjectService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteObject",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"BucketName": resourceData.Get("bucket_name"),
				"ObjectName": s.ReadResourceId(resourceData.Id()),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {

				if d.Get("version_ids") != nil && len(d.Get("version_ids").(*schema.Set).List()) > 0 {
					for _, vv := range d.Get("version_ids").(*schema.Set).List() {
						condition := make(map[string]interface{})
						condition["versionId"] = vv
						//remove Object-with-version
						logger.Debug(logger.RespFormat, call.Action, condition)
						_, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
							HttpMethod: bp.DELETE,
							Domain:     (*call.SdkParam)["BucketName"].(string),
							Path:       []string{(*call.SdkParam)["ObjectName"].(string)},
						}, &condition)
						if err != nil {
							return nil, err
						}
					}
				} else {
					//remove Object-no-version
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					return s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
						HttpMethod: bp.DELETE,
						Domain:     (*call.SdkParam)["BucketName"].(string),
						Path:       []string{(*call.SdkParam)["ObjectName"].(string)},
					}, nil)
				}

				return nil, nil
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading tos object on delete %q, %w", s.ReadResourceId(d.Id()), callErr))
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

func (s *VestackTosObjectService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	name, ok := data.GetOk("object_name")
	bucketName, _ := data.GetOk("bucket_name")
	return bp.DataSourceInfo{
		ServiceCategory: bp.ServiceBypass,
		RequestConverts: map[string]bp.RequestConvert{
			"bucket_name": {
				ConvertType: bp.ConvertDefault,
				SpecialParam: &bp.SpecialParam{
					Type: bp.DomainParam,
				},
			},
			"object_name": {
				Ignore: true,
			},
		},
		NameField:    "Key",
		IdField:      "ObjectId",
		CollectField: "objects",
		ResponseConverts: map[string]bp.ResponseConvert{
			"Key": {
				TargetField: "name",
			},
		},
		ExtraData: func(sourceData []interface{}) (extraData []interface{}, err error) {
			for _, v := range sourceData {
				if ok {
					if name.(string) == v.(map[string]interface{})["Key"].(string) {
						v.(map[string]interface{})["ObjectId"] = bucketName.(string) + ":" + v.(map[string]interface{})["Key"].(string)
						extraData = append(extraData, v)
						break
					} else {
						continue
					}
				} else {
					v.(map[string]interface{})["ObjectId"] = bucketName.(string) + ":" + v.(map[string]interface{})["Key"].(string)
					extraData = append(extraData, v)
				}

			}
			return extraData, err
		},

		EachResource: func(sourceData []interface{}, d *schema.ResourceData) ([]interface{}, error) {
			var newSourceData []interface{}
			for _, v := range sourceData {
				var (
					key     interface{}
					newData map[string]interface{}
					err     error
				)
				key, err = bp.ObtainSdkValue("Key", v)
				if err != nil {
					return nil, err
				}

				if str, ok1 := key.(string); ok1 {
					newData, err = s.ReadResource(d, str)
					if err != nil {
						return nil, err
					}
				}

				if v1, ok1 := v.(map[string]interface{}); ok1 {
					for k, value := range newData {
						if _, ok2 := v1[k]; !ok2 {
							v1[k] = value
						}
					}
					newSourceData = append(newSourceData, v1)
				}
			}
			return newSourceData, nil
		},
	}
}

func (s *VestackTosObjectService) ReadResourceId(id string) string {
	return id[strings.Index(id, ":")+1:]
}

func (s *VestackTosObjectService) beforePutObjectAcl() bp.BeforeCallFunc {
	return func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
		logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
		data, err := s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
			HttpMethod: bp.GET,
			Domain:     (*call.SdkParam)[bp.BypassDomain].(string),
			Path:       (*call.SdkParam)[bp.BypassPath].([]string),
			UrlParam: map[string]string{
				"acl": "",
			},
		}, nil)
		return bp.BeforeTosPutAcl(d, call, data, err)
	}
}

func (s *VestackTosObjectService) executePutObjectAcl() bp.ExecuteCallFunc {
	return func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
		logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
		//PutAcl
		param := (*call.SdkParam)[bp.BypassParam].(map[string]interface{})
		return s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
			HttpMethod:  bp.PUT,
			ContentType: bp.ApplicationJSON,
			Domain:      (*call.SdkParam)[bp.BypassDomain].(string),
			Path:        (*call.SdkParam)[bp.BypassPath].([]string),
			UrlParam: map[string]string{
				"acl": "",
			},
		}, &param)
	}
}

func (s *VestackTosObjectService) createOrUpdateObjectAcl(resourceData *schema.ResourceData, resource *schema.Resource, isUpdate bool) bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "PutObjectAcl",
			ConvertMode:     bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"bucket_name": {
					ConvertType: bp.ConvertDefault,
					TargetField: "BucketName",
					SpecialParam: &bp.SpecialParam{
						Type: bp.DomainParam,
					},
					ForceGet: isUpdate,
				},
				"object_name": {
					ConvertType: bp.ConvertDefault,
					TargetField: "ObjectName",
					SpecialParam: &bp.SpecialParam{
						Type:  bp.PathParam,
						Index: 0,
					},
					ForceGet: isUpdate,
				},
				"account_acl": {
					ConvertType: bp.ConvertListN,
					TargetField: "Grants",
					ForceGet:    isUpdate,
					NextLevelConvert: map[string]bp.RequestConvert{
						"account_id": {
							ConvertType: bp.ConvertDefault,
							TargetField: "Grantee.ID",
							ForceGet:    isUpdate,
						},
						"acl_type": {
							ConvertType: bp.ConvertDefault,
							TargetField: "Grantee.Type",
							ForceGet:    isUpdate,
						},
						"permission": {
							ConvertType: bp.ConvertDefault,
							TargetField: "Permission",
							ForceGet:    isUpdate,
						},
					},
				},
			},
			BeforeCall:  s.beforePutObjectAcl(),
			ExecuteCall: s.executePutObjectAcl(),
			//Refresh: &bp.StateRefresh{
			//	Target:  []string{"Success"},
			//	Timeout: resourceData.Timeout(schema.TimeoutCreate),
			//},
		},
	}
	//如果出现acl缓存的问题 这里再打开 暂时去掉 不再等待60s
	//if isUpdate && !resourceData.HasChange("file_path") && !resourceData.HasChanges("content_md5") {
	//	callback.Call.Refresh = &bp.StateRefresh{
	//		Target:  []string{"Success"},
	//		Timeout: resourceData.Timeout(schema.TimeoutCreate),
	//	}
	//}
	return callback
}

func (s *VestackTosObjectService) createOrReplaceObject(resourceData *schema.ResourceData, resource *schema.Resource, isUpdate bool) bp.Callback {
	return bp.Callback{
		Call: bp.SdkCall{
			ServiceCategory: bp.ServiceBypass,
			Action:          "PutObject",
			ConvertMode:     bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"bucket_name": {
					ConvertType: bp.ConvertDefault,
					TargetField: "BucketName",
					SpecialParam: &bp.SpecialParam{
						Type: bp.DomainParam,
					},
					ForceGet: isUpdate,
				},
				"object_name": {
					ConvertType: bp.ConvertDefault,
					TargetField: "ObjectName",
					SpecialParam: &bp.SpecialParam{
						Type:  bp.PathParam,
						Index: 0,
					},
					ForceGet: isUpdate,
				},
				"public_acl": {
					ConvertType: bp.ConvertDefault,
					TargetField: "x-tos-acl",
					SpecialParam: &bp.SpecialParam{
						Type: bp.HeaderParam,
					},
					ForceGet: isUpdate,
				},
				"storage_class": {
					ConvertType: bp.ConvertDefault,
					TargetField: "x-tos-storage-class",
					SpecialParam: &bp.SpecialParam{
						Type: bp.HeaderParam,
					},
					ForceGet: isUpdate,
				},
				"content_type": {
					ConvertType: bp.ConvertDefault,
					TargetField: "Content-Type",
					SpecialParam: &bp.SpecialParam{
						Type: bp.HeaderParam,
					},
					ForceGet: isUpdate,
				},
				"content_md5": {
					ConvertType: bp.ConvertDefault,
					TargetField: "Content-MD5",
					Convert: func(data *schema.ResourceData, i interface{}) interface{} {
						b, _ := hex.DecodeString(i.(string))
						return base64.StdEncoding.EncodeToString(b)
					},
					SpecialParam: &bp.SpecialParam{
						Type: bp.HeaderParam,
					},
					ForceGet: isUpdate,
				},
				"file_path": {
					ConvertType: bp.ConvertDefault,
					TargetField: "file-path",
					SpecialParam: &bp.SpecialParam{
						Type: bp.FilePathParam,
					},
					ForceGet: isUpdate,
				},
				"encryption": {
					ConvertType: bp.ConvertDefault,
					TargetField: "x-tos-server-side-encryption",
					SpecialParam: &bp.SpecialParam{
						Type: bp.HeaderParam,
					},
					ForceGet: isUpdate,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if _, ok := (*call.SdkParam)[bp.BypassHeader].(map[string]string)["Content-MD5"]; ok {
					(*call.SdkParam)[bp.BypassHeader].(map[string]string)["x-tos-meta-content-md5"] = d.Get("content_md5").(string)
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//创建Object
				return s.Client.BypassSvcClient.DoBypassSvcCall(bp.BypassSvcInfo{
					HttpMethod:  bp.PUT,
					Domain:      (*call.SdkParam)[bp.BypassDomain].(string),
					Header:      (*call.SdkParam)[bp.BypassHeader].(map[string]string),
					Path:        (*call.SdkParam)[bp.BypassPath].([]string),
					ContentPath: (*call.SdkParam)[bp.BypassFilePath].(string),
				}, nil)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId((*call.SdkParam)[bp.BypassDomain].(string) + ":" + (*call.SdkParam)[bp.BypassPath].([]string)[0])
				return nil
			},
		},
	}
}
