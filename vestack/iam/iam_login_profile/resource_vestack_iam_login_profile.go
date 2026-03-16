package iam_login_profile

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Login profile can be imported using the UserName, e.g.
```
$ terraform import vestack_iam_login_profile.default user_name
```

*/

func ResourceVestackIamLoginProfile() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackIamLoginProfileCreate,
		Read:   resourceVestackIamLoginProfileRead,
		Update: resourceVestackIamLoginProfileUpdate,
		Delete: resourceVestackIamLoginProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The user name.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The password.",
			},
			"login_allowed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "The flag of login allowed.",
			},
			"password_reset_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Is required reset password when next time login in.",
			},
			"safe_auth_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "The flag of safe auth.",
			},
			"safe_auth_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The type of safe auth.",
			},
			"safe_auth_exempt_required": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The flag of safe auth exempt required.",
			},
			"safe_auth_exempt_unit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unit of safe auth exempt.",
			},
			"safe_auth_exempt_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The duration of safe auth exempt.",
			},
			"user_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The user id.",
			},
			"password_expire_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The password expire at.",
			},
			"last_reset_password_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The last reset password time.",
			},
			"last_login_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last login date.",
			},
			"last_login_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last login ip.",
			},
			"login_locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag of login locked.",
			},
			"create_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create date.",
			},
			"update_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update date.",
			},
		},
	}
	return resource
}

func resourceVestackIamLoginProfileCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamLoginProfileService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceVestackIamLoginProfile())
	if err != nil {
		return fmt.Errorf("error on creating login profile %q, %s", d.Id(), err)
	}
	return resourceVestackIamLoginProfileRead(d, meta)
}

func resourceVestackIamLoginProfileRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamLoginProfileService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceVestackIamLoginProfile())
	if err != nil {
		return fmt.Errorf("error on reading login profile %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackIamLoginProfileUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamLoginProfileService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceVestackIamLoginProfile())
	if err != nil {
		return fmt.Errorf("error on updating login profile %q, %s", d.Id(), err)
	}
	return resourceVestackIamLoginProfileRead(d, meta)
}

func resourceVestackIamLoginProfileDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamLoginProfileService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceVestackIamLoginProfile())
	if err != nil {
		return fmt.Errorf("error on deleting login profile %q, %s", d.Id(), err)
	}
	return err
}
