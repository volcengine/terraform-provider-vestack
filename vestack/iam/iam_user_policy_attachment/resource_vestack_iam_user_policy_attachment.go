package iam_user_policy_attachment

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Iam user policy attachment can be imported using the UserName:PolicyName:PolicyType, e.g.
```
$ terraform import vestack_iam_user_policy_attachment.default TerraformTestUser:TerraformTestPolicy:Custom
```

*/

func ResourceVestackIamUserPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackIamUserPolicyAttachmentCreate,
		Read:   resourceVestackIamUserPolicyAttachmentRead,
		Delete: resourceVestackIamUserPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: iamUserPolicyAttachmentImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the user.",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Policy.",
			},
			"policy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"System", "Custom"}, false),
				Description:  "The type of the Policy.",
			},
		},
	}
}

func resourceVestackIamUserPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	iamUserPolicyAttachmentService := NewIamUserPolicyAttachmentService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(iamUserPolicyAttachmentService, d, ResourceVestackIamUserPolicyAttachment()); err != nil {
		return fmt.Errorf("error on creating iam user policy attachment %q, %w", d.Id(), err)
	}
	return resourceVestackIamUserPolicyAttachmentRead(d, meta)
}

func resourceVestackIamUserPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	iamUserPolicyAttachmentService := NewIamUserPolicyAttachmentService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(iamUserPolicyAttachmentService, d, ResourceVestackIamUserPolicyAttachment()); err != nil {
		return fmt.Errorf("error on reading iam user policy attachment %q, %w", d.Id(), err)
	}
	return nil
}

func resourceVestackIamUserPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	iamUserPolicyAttachmentService := NewIamUserPolicyAttachmentService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(iamUserPolicyAttachmentService, d, ResourceVestackIamUserPolicyAttachment()); err != nil {
		return fmt.Errorf("error on deleting iam user policy attachment %q, %w", d.Id(), err)
	}
	return nil
}
