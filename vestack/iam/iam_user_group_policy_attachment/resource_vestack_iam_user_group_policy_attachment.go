package iam_user_group_policy_attachment

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
IamUserGroupPolicyAttachment can be imported using the user group name and policy name, e.g.
```
$ terraform import vestack_iam_user_group_policy_attachment.default userGroupName:policyName
```

*/

func ResourceVestackIamUserGroupPolicyAttachment() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackIamUserGroupPolicyAttachmentCreate,
		Read:   resourceVestackIamUserGroupPolicyAttachmentRead,
		Delete: resourceVestackIamUserGroupPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				items := strings.Split(data.Id(), ":")
				if len(items) != 2 {
					return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
				}
				if err := data.Set("user_group_name", items[0]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				if err := data.Set("policy_name", items[1]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				return []*schema.ResourceData{data}, nil
			},
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
				Description: "The user group name.",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The policy name.",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Strategy types, System strategy, Custom strategy.",
			},
		},
	}
	return resource
}

func resourceVestackIamUserGroupPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupPolicyAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackIamUserGroupPolicyAttachment())
	if err != nil {
		return fmt.Errorf("error on creating iam_user_group_policy_attachment %q, %s", d.Id(), err)
	}
	return resourceVestackIamUserGroupPolicyAttachmentRead(d, meta)
}

func resourceVestackIamUserGroupPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupPolicyAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackIamUserGroupPolicyAttachment())
	if err != nil {
		return fmt.Errorf("error on reading iam_user_group_policy_attachment %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackIamUserGroupPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupPolicyAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackIamUserGroupPolicyAttachment())
	if err != nil {
		return fmt.Errorf("error on deleting iam_user_group_policy_attachment %q, %s", d.Id(), err)
	}
	return err
}
