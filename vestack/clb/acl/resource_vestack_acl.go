package acl

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Acl can be imported using the id, e.g.
```
$ terraform import vestack_acl.default acl-mizl7m1kqccg5smt1bdpijuj
```

*/

func ResourceVestackAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackAclCreate,
		Read:   resourceVestackAclRead,
		Update: resourceVestackAclUpdate,
		Delete: resourceVestackAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of Acl.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Acl.",
			},
			"acl_entries": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "The acl entry set of the Acl.",
				Set:         ve.ClbAclEntryHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The description of the AclEntry.",
						},
						"entry": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The content of the AclEntry.",
						},
					},
				},
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ProjectName of the Acl.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of Acl.",
			},
		},
	}
}

func resourceVestackAclCreate(d *schema.ResourceData, meta interface{}) (err error) {
	aclService := NewAclService(meta.(*ve.SdkClient))
	err = ve.DefaultDispatcher().Create(aclService, d, ResourceVestackAcl())
	if err != nil {
		return fmt.Errorf("error on creating acl %q, %w", d.Id(), err)
	}
	return resourceVestackAclRead(d, meta)
}

func resourceVestackAclRead(d *schema.ResourceData, meta interface{}) (err error) {
	aclService := NewAclService(meta.(*ve.SdkClient))
	err = ve.DefaultDispatcher().Read(aclService, d, ResourceVestackAcl())
	if err != nil {
		return fmt.Errorf("error on reading acl %q, %w", d.Id(), err)
	}
	return err
}

func resourceVestackAclUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	aclService := NewAclService(meta.(*ve.SdkClient))
	err = ve.DefaultDispatcher().Update(aclService, d, ResourceVestackAcl())
	if err != nil {
		return fmt.Errorf("error on updating acl %q, %w", d.Id(), err)
	}
	return resourceVestackAclRead(d, meta)
}

func resourceVestackAclDelete(d *schema.ResourceData, meta interface{}) (err error) {
	aclService := NewAclService(meta.(*ve.SdkClient))
	err = ve.DefaultDispatcher().Delete(aclService, d, ResourceVestackAcl())
	if err != nil {
		return fmt.Errorf("error on deleting acl %q, %w", d.Id(), err)
	}
	return err
}
