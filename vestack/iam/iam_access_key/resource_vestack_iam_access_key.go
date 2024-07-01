package iam_access_key

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Iam access key don't support import

*/

func ResourceVestackIamAccessKey() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackIamAccessKeyCreate,
		Read:   resourceVestackIamAccessKeyRead,
		Update: resourceVestackIamAccessKeyUpdate,
		Delete: resourceVestackIamAccessKeyDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The user name.",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "active",
				ValidateFunc: validation.StringInSlice([]string{"active", "inactive"}, false),
				Description:  "The status of the access key, Optional choice contains `active` or `inactive`.",
			},
			"pgp_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Either a base-64 encoded PGP public key, or a keybase username in the form `keybase:some_person_that_exists`.",
			},
			"secret_file": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The file to save the access id and secret. Strongly suggest you to specified it when you creating access key, otherwise, you wouldn't get its secret ever.",
			},
			"secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The secret of the access key.",
			},
			"encrypted_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The encrypted secret of the access key by pgp key, base64 encoded.",
			},
			"key_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key fingerprint of the encrypted secret.",
			},
			"create_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create date of the access key.",
			},
		},
	}
	return resource
}

func resourceVestackIamAccessKeyCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamAccessKeyService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceVestackIamAccessKey())
	if err != nil {
		return fmt.Errorf("error on creating access key  %q, %s", d.Id(), err)
	}
	return resourceVestackIamAccessKeyRead(d, meta)
}

func resourceVestackIamAccessKeyRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamAccessKeyService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceVestackIamAccessKey())
	if err != nil {
		return fmt.Errorf("error on reading access key %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackIamAccessKeyUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamAccessKeyService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceVestackIamAccessKey())
	if err != nil {
		return fmt.Errorf("error on updating access key %q, %s", d.Id(), err)
	}
	return resourceVestackIamAccessKeyRead(d, meta)
}

func resourceVestackIamAccessKeyDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamAccessKeyService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceVestackIamAccessKey())
	if err != nil {
		return fmt.Errorf("error on deleting access key %q, %s", d.Id(), err)
	}
	return err
}
