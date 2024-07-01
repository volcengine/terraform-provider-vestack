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
				Description: "The flag of login allowed.",
			},
			"password_reset_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Is required reset password when next time login in.",
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
