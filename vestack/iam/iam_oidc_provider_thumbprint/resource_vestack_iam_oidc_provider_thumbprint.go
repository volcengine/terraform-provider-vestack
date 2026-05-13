package iam_oidc_provider_thumbprint

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Iam OidcProviderThumbprint key don't support import

*/

func ResourceVestackIamOidcProviderThumbprint() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackIamOidcProviderThumbprintCreate,
		Read:   resourceVestackIamOidcProviderThumbprintRead,
		Delete: resourceVestackIamOidcProviderThumbprintDelete,
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
			"thumbprint": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The thumbprint of the OIDC provider.",
			},
		},
	}
}

func resourceVestackIamOidcProviderThumbprintCreate(d *schema.ResourceData, meta interface{}) error {
	service := NewIamOidcProviderThumbprintService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Create(service, d, ResourceVestackIamOidcProviderThumbprint())
}

func resourceVestackIamOidcProviderThumbprintRead(d *schema.ResourceData, meta interface{}) error {
	service := NewIamOidcProviderThumbprintService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Read(service, d, ResourceVestackIamOidcProviderThumbprint())
}

func resourceVestackIamOidcProviderThumbprintDelete(d *schema.ResourceData, meta interface{}) error {
	service := NewIamOidcProviderThumbprintService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Delete(service, d, ResourceVestackIamOidcProviderThumbprint())
}
