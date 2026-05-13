package iam_user_group

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
IamUserGroup can be imported using the id, e.g.
```
$ terraform import vestack_iam_user_group.default user_group_name
```

*/

func ResourceVestackIamUserGroup() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackIamUserGroupCreate,
		Read:   resourceVestackIamUserGroupRead,
		Update: resourceVestackIamUserGroupUpdate,
		Delete: resourceVestackIamUserGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"user_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the user group.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the user group.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The display name of the user group.",
			},
		},
	}
	return resource
}

func resourceVestackIamUserGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackIamUserGroup())
	if err != nil {
		return fmt.Errorf("error on creating iam_user_group %q, %s", d.Id(), err)
	}
	return resourceVestackIamUserGroupRead(d, meta)
}

func resourceVestackIamUserGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackIamUserGroup())
	if err != nil {
		return fmt.Errorf("error on reading iam_user_group %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackIamUserGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackIamUserGroup())
	if err != nil {
		return fmt.Errorf("error on updating iam_user_group %q, %s", d.Id(), err)
	}
	return resourceVestackIamUserGroupRead(d, meta)
}

func resourceVestackIamUserGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackIamUserGroup())
	if err != nil {
		return fmt.Errorf("error on deleting iam_user_group %q, %s", d.Id(), err)
	}
	return err
}
