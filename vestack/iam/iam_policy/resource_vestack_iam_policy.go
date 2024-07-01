package iam_policy

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Iam policy can be imported using the id, e.g.
```
$ terraform import vestack_iam_policy.default TerraformTestPolicy
```

*/

func ResourceVestackIamPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackIamPolicyCreate,
		Read:   resourceVestackIamPolicyRead,
		Update: resourceVestackIamPolicyUpdate,
		Delete: resourceVestackIamPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Policy.",
			},
			"policy_document": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The document of the Policy.",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Policy.",
			},
			"policy_trn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource name of the Policy.",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the Policy.",
			},
			"create_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create time of the Policy.",
			},
			"update_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time of the Policy.",
			},
		},
	}
}

func resourceVestackIamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	iamPolicyService := NewIamPolicyService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(iamPolicyService, d, ResourceVestackIamPolicy()); err != nil {
		return fmt.Errorf("error on creating iam policy %q, %w", d.Id(), err)
	}
	return resourceVestackIamPolicyRead(d, meta)
}

func resourceVestackIamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	iamPolicyService := NewIamPolicyService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(iamPolicyService, d, ResourceVestackIamPolicy()); err != nil {
		return fmt.Errorf("error on reading iam policy %q, %w", d.Id(), err)
	}
	return nil
}

func resourceVestackIamPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	iamPolicyService := NewIamPolicyService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Update(iamPolicyService, d, ResourceVestackIamPolicy()); err != nil {
		return fmt.Errorf("error on updating iam policy %q, %w", d.Id(), err)
	}
	return resourceVestackIamPolicyRead(d, meta)
}

func resourceVestackIamPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	iamPolicyService := NewIamPolicyService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(iamPolicyService, d, ResourceVestackIamPolicy()); err != nil {
		return fmt.Errorf("error on deleting iam policy %q, %w", d.Id(), err)
	}
	return nil
}
