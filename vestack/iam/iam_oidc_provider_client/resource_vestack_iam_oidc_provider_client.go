package iam_oidc_provider_client

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Iam OidcProvider key don't support import

*/

func ResourceVestackIamOidcProviderClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackIamOidcProviderClientCreate,
		Read:   resourceVestackIamOidcProviderClientRead,
		Delete: resourceVestackIamOidcProviderClientDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"oidc_provider_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the OIDC provider.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The client id of the OIDC provider.",
			},
		},
	}
}

func resourceVestackIamOidcProviderClientCreate(d *schema.ResourceData, meta interface{}) error {
	service := NewIamOidcProviderClientService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Create(service, d, ResourceVestackIamOidcProviderClient())
}

func resourceVestackIamOidcProviderClientRead(d *schema.ResourceData, meta interface{}) error {
	service := NewIamOidcProviderClientService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Read(service, d, ResourceVestackIamOidcProviderClient())
}

func resourceVestackIamOidcProviderClientDelete(d *schema.ResourceData, meta interface{}) error {
	service := NewIamOidcProviderClientService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Delete(service, d, ResourceVestackIamOidcProviderClient())
}
