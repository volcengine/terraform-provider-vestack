package vke_permission

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
VkePermission can be imported using the id, e.g.
```
$ terraform import vestack_vke_permission.default resource_id
```

*/

func ResourceVestackVkePermission() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackVkePermissionCreate,
		Read:   resourceVestackVkePermissionRead,
		Delete: resourceVestackVkePermissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"role_domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The types of permissions granted to IAM users or roles. Valid values: `namespace`, `cluster`, `all_clusters`. \n" +
					"When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The cluster ID that needs to be authorized to IAM users or roles.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The namespace that needs to be authorized to IAM users or roles.",
			},
			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of RBAC role. The following RBAC permissions can be granted: custom role name, system preset role names.",
			},
			"is_custom_role": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Whether the RBAC role is a custom role. Default is false.",
			},
			"grantee_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the grantee.",
			},
			"grantee_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"User"}, false),
				Description:  "The type of the grantee. Valid values: `User`.",
			},

			// computed fields
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the RBAC Permission.",
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The message of the RBAC Permission.",
			},
			"granted_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The granted time of the RBAC Permission.",
			},
			"revoked_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The revoked time of the RBAC Permission.",
			},
			"authorizer_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the Authorizer.",
			},
			"authorizer_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Authorizer.",
			},
			"authorizer_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the Authorizer.",
			},
			"authorized_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The authorized time of the RBAC Permission.",
			},
			"kube_role_binding_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Kube Role Binding.",
			},
		},
	}
	return resource
}

func resourceVestackVkePermissionCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewVkePermissionService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackVkePermission())
	if err != nil {
		return fmt.Errorf("error on creating vke_permission %q, %s", d.Id(), err)
	}
	return resourceVestackVkePermissionRead(d, meta)
}

func resourceVestackVkePermissionRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewVkePermissionService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackVkePermission())
	if err != nil {
		return fmt.Errorf("error on reading vke_permission %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackVkePermissionUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewVkePermissionService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackVkePermission())
	if err != nil {
		return fmt.Errorf("error on updating vke_permission %q, %s", d.Id(), err)
	}
	return resourceVestackVkePermissionRead(d, meta)
}

func resourceVestackVkePermissionDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewVkePermissionService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackVkePermission())
	if err != nil {
		return fmt.Errorf("error on deleting vke_permission %q, %s", d.Id(), err)
	}
	return err
}
