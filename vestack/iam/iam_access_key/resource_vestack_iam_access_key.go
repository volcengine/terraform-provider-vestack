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
				Description: "The user name. If not specified, the current user is used.",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "active",
				ValidateFunc: validation.StringInSlice([]string{"active", "inactive"}, false),
				Description:  "The status of the access key, Optional choice contains `active` or `inactive`.",
			},
			"access_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The access key id.",
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The secret access key.",
			},
			"update_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update date of the access key.",
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
