package certificate

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ve "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackCertificateService struct {
	Client     *ve.SdkClient
	Dispatcher *ve.Dispatcher
}

func NewCertificateService(c *ve.SdkClient) *VestackCertificateService {
	return &VestackCertificateService{
		Client:     c,
		Dispatcher: &ve.Dispatcher{},
	}
}

func (s *VestackCertificateService) GetClient() *ve.SdkClient {
	return s.Client
}

func (s *VestackCertificateService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return ve.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		clbClient := s.Client.ClbClient
		action := "DescribeCertificates"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = clbClient.DescribeCertificatesCommon(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = clbClient.DescribeCertificatesCommon(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = ve.ObtainSdkValue("Result.Certificates", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Certificates is not Slice")
		}
		return data, err
	})
}

func (s *VestackCertificateService) ReadResource(resourceData *schema.ResourceData, certificateId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if certificateId == "" {
		certificateId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"CertificateIds.1": certificateId,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Certificate %s not exist ", certificateId)
	}
	return data, err
}

func (s *VestackCertificateService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (VestackCertificateService) WithResourceResponseHandlers(certificate map[string]interface{}) []ve.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]ve.ResponseConvert, error) {
		return certificate, nil, nil
	}
	return []ve.ResourceResponseHandler{handler}

}

func (s *VestackCertificateService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []ve.Callback {
	callback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "UploadCertificate",
			ConvertMode: ve.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//??????certificate
				return s.Client.ClbClient.UploadCertificateCommon(call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *ve.SdkClient, resp *map[string]interface{}, call ve.SdkCall) error {
				//?????? ???????????? ??????????????????????????? ???????????????
				id, _ := ve.ObtainSdkValue("Result.CertificateId", *resp)
				d.SetId(id.(string))
				return nil
			},
		},
	}
	return []ve.Callback{callback}

}

func (s *VestackCertificateService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []ve.Callback {
	return []ve.Callback{}
}

func (s *VestackCertificateService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []ve.Callback {
	callback := ve.Callback{
		Call: ve.SdkCall{
			Action:      "DeleteCertificate",
			ConvertMode: ve.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"CertificateId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//??????Certificate
				return s.Client.ClbClient.DeleteCertificateCommon(call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *ve.SdkClient, call ve.SdkCall, baseErr error) error {
				//?????????????????????
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if ve.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading certificate on delete %q, %w", d.Id(), callErr))
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
	return []ve.Callback{callback}
}

func (s *VestackCertificateService) DatasourceResources(*schema.ResourceData, *schema.Resource) ve.DataSourceInfo {
	return ve.DataSourceInfo{
		RequestConverts: map[string]ve.RequestConvert{
			"ids": {
				TargetField: "CertificateIds",
				ConvertType: ve.ConvertWithN,
			},
		},
		NameField:    "CertificateName",
		IdField:      "CertificateId",
		CollectField: "certificates",
		ResponseConverts: map[string]ve.ResponseConvert{
			"CertificateId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *VestackCertificateService) ReadResourceId(id string) string {
	return id
}
