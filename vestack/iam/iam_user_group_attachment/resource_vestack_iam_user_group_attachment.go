package iam_user_group_attachment

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
IamUserGroupAttachment can be imported using the id, e.g.
```
$ terraform import vestack_iam_user_group_attachment.default user_group_id:user_id
```

*/

func ResourceVestackIamUserGroupAttachment() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackIamUserGroupAttachmentCreate,
		Read:   resourceVestackIamUserGroupAttachmentRead,
		Delete: resourceVestackIamUserGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: importIamUserGroupAttachment,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"user_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the user group.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the user.",
			},
		},
	}
	return resource
}

func resourceVestackIamUserGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackIamUserGroupAttachment())
	if err != nil {
		return fmt.Errorf("error on creating iam_user_group_attachment %q, %s", d.Id(), err)
	}
	return resourceVestackIamUserGroupAttachmentRead(d, meta)
}

func resourceVestackIamUserGroupAttachmentRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackIamUserGroupAttachment())
	if err != nil {
		return fmt.Errorf("error on reading iam_user_group_attachment %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackIamUserGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackIamUserGroupAttachment())
	if err != nil {
		return fmt.Errorf("error on deleting iam_user_group_attachment %q, %s", d.Id(), err)
	}
	return err
}

func importIamUserGroupAttachment(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must be of the form user_group_id:user_id")
	}
	err = data.Set("user_group_name", items[0])
	if err != nil {
		return []*schema.ResourceData{data}, err
	}
	err = data.Set("user_name", items[1])
	if err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
