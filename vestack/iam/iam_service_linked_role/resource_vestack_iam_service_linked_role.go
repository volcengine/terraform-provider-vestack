package iam_service_linked_role

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Iam service linked role can be imported using the servicx name and the service linked role name, e.g.
```
$ terraform import vestack_iam_service_linked_role.default ecs:ServiceRoleForEcs
```

*/

func ResourceVestackIamServiceLinkedRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackIamServiceLinkedRoleCreate,
		Read:   resourceVestackIamServiceLinkedRoleRead,
		Delete: resourceVestackIamServiceLinkedRoleDelete,
		Importer: &schema.ResourceImporter{
			State: iamServiceLinkedRoleImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the service.",
			},
			"role_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the service linked role.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the service linked Role.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the service linked Role.",
			},
			"max_session_duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The max session duration of the service linked Role.",
			},
			"trn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource name of the service linked Role.",
			},
			"trust_policy_document": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The trust policy document of the service linked Role.",
			},
		},
	}
}

func resourceVestackIamServiceLinkedRoleCreate(d *schema.ResourceData, meta interface{}) error {
	IamServiceLinkedRoleService := NewIamServiceLinkedRoleService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(IamServiceLinkedRoleService, d, ResourceVestackIamServiceLinkedRole()); err != nil {
		return fmt.Errorf("error on creating iam service linked role %q, %w", d.Id(), err)
	}
	return resourceVestackIamServiceLinkedRoleRead(d, meta)
}

func resourceVestackIamServiceLinkedRoleRead(d *schema.ResourceData, meta interface{}) error {
	IamServiceLinkedRoleService := NewIamServiceLinkedRoleService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(IamServiceLinkedRoleService, d, ResourceVestackIamServiceLinkedRole()); err != nil {
		return fmt.Errorf("error on reading iam service linked role %q, %w", d.Id(), err)
	}
	return nil
}

func resourceVestackIamServiceLinkedRoleDelete(d *schema.ResourceData, meta interface{}) error {
	IamServiceLinkedRoleService := NewIamServiceLinkedRoleService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(IamServiceLinkedRoleService, d, ResourceVestackIamServiceLinkedRole()); err != nil {
		return fmt.Errorf("error on deleting iam service linked role %q, %w", d.Id(), err)
	}
	return nil
}

func iamServiceLinkedRoleImporter(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
	}
	if err := data.Set("service_name", items[0]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	if err := data.Set("role_name", items[1]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
