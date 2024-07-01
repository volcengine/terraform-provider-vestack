package security_group

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
SecurityGroup can be imported using the id, e.g.
```
$ terraform import vestack_security_group.default sg-273ycgql3ig3k7fap8t3dyvqx
```

*/

func ResourceVestackSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackSecurityGroupCreate,
		Read:   resourceVestackSecurityGroupRead,
		Update: resourceVestackSecurityGroupUpdate,
		Delete: resourceVestackSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the VPC.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of SecurityGroup.",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of SecurityGroup.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of SecurityGroup.",
			},
			"security_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of SecurityGroup.",
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				//ForceNew:    true,
				Description: "The ProjectName of SecurityGroup.",
			},
			"tags": bp.TagsSchema(),
		},
	}
}

func resourceVestackSecurityGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	securityGroupService := NewSecurityGroupService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(securityGroupService, d, ResourceVestackSecurityGroup())
	if err != nil {
		return fmt.Errorf("error on creating securityGroupService  %q, %w", d.Id(), err)
	}
	return resourceVestackSecurityGroupRead(d, meta)
}

func resourceVestackSecurityGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	securityGroupService := NewSecurityGroupService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(securityGroupService, d, ResourceVestackSecurityGroup())
	if err != nil {
		return fmt.Errorf("error on reading securityGroupService %q, %w", d.Id(), err)
	}
	return err
}

func resourceVestackSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	securityGroupService := NewSecurityGroupService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(securityGroupService, d, ResourceVestackSecurityGroup())
	if err != nil {
		return fmt.Errorf("error on updating securityGroupService  %q, %w", d.Id(), err)
	}
	return resourceVestackSecurityGroupRead(d, meta)
}

func resourceVestackSecurityGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	securityGroupService := NewSecurityGroupService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(securityGroupService, d, ResourceVestackSecurityGroup())
	if err != nil {
		return fmt.Errorf("error on deleting securityGroupService %q, %w", d.Id(), err)
	}
	return err
}
