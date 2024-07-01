package iam_role

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Iam role can be imported using the id, e.g.
```
$ terraform import vestack_iam_role.default TerraformTestRole
```

*/

func ResourceVestackIamRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackIamRoleCreate,
		Read:   resourceVestackIamRoleRead,
		Update: resourceVestackIamRoleUpdate,
		Delete: resourceVestackIamRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"trust_policy_document": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The trust policy document of the Role.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldMap := make(map[string]interface{})
					newMap := make(map[string]interface{})

					_ = json.Unmarshal([]byte(old), &oldMap)
					_ = json.Unmarshal([]byte(new), &newMap)

					oldStr, _ := json.MarshalIndent(oldMap, "", "\t")
					newStr, _ := json.MarshalIndent(newMap, "", "\t")
					return string(oldStr) == string(newStr)
				},
			},
			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Role.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the Role.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Role.",
			},
			"max_session_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The max session duration of the Role.",
			},
			"trn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource name of the Role.",
			},
		},
	}
}

func resourceVestackIamRoleCreate(d *schema.ResourceData, meta interface{}) error {
	iamRoleService := NewIamRoleService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(iamRoleService, d, ResourceVestackIamRole()); err != nil {
		return fmt.Errorf("error on creating iam role %q, %w", d.Id(), err)
	}
	return resourceVestackIamRoleRead(d, meta)
}

func resourceVestackIamRoleRead(d *schema.ResourceData, meta interface{}) error {
	iamRoleService := NewIamRoleService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(iamRoleService, d, ResourceVestackIamRole()); err != nil {
		return fmt.Errorf("error on reading iam role %q, %w", d.Id(), err)
	}
	return nil
}

func resourceVestackIamRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	iamRoleService := NewIamRoleService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Update(iamRoleService, d, ResourceVestackIamRole()); err != nil {
		return fmt.Errorf("error on updating iam role %q, %w", d.Id(), err)
	}
	return resourceVestackIamRoleRead(d, meta)
}

func resourceVestackIamRoleDelete(d *schema.ResourceData, meta interface{}) error {
	iamRoleService := NewIamRoleService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(iamRoleService, d, ResourceVestackIamRole()); err != nil {
		return fmt.Errorf("error on deleting iam role %q, %w", d.Id(), err)
	}
	return nil
}
