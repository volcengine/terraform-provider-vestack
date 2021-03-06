package server_group

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
ServerGroup can be imported using the id, e.g.
```
$ terraform import vestack_server_group.default rsp-273yv0kir1vk07fap8tt9jtwg
```

*/

func ResourceVestackServerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackServerGroupCreate,
		Read:   resourceVestackServerGroupRead,
		Update: resourceVestackServerGroupUpdate,
		Delete: resourceVestackServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"server_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the ServerGroup.",
			},
			"load_balancer_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the Clb.",
			},
			"server_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the ServerGroup.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of ServerGroup.",
			},
		},
	}
}

func resourceVestackServerGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	serverGroupService := NewServerGroupService(meta.(*ve.SdkClient))
	err = serverGroupService.Dispatcher.Create(serverGroupService, d, ResourceVestackServerGroup())
	if err != nil {
		return fmt.Errorf("error on creating serverGroup  %q, %w", d.Id(), err)
	}
	return resourceVestackServerGroupRead(d, meta)
}

func resourceVestackServerGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	serverGroupService := NewServerGroupService(meta.(*ve.SdkClient))
	err = serverGroupService.Dispatcher.Read(serverGroupService, d, ResourceVestackServerGroup())
	if err != nil {
		return fmt.Errorf("error on reading serverGroup %q, %w", d.Id(), err)
	}
	return err
}

func resourceVestackServerGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	serverGroupService := NewServerGroupService(meta.(*ve.SdkClient))
	err = serverGroupService.Dispatcher.Update(serverGroupService, d, ResourceVestackServerGroup())
	if err != nil {
		return fmt.Errorf("error on updating serverGroup  %q, %w", d.Id(), err)
	}
	return resourceVestackServerGroupRead(d, meta)
}

func resourceVestackServerGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	serverGroupService := NewServerGroupService(meta.(*ve.SdkClient))
	err = serverGroupService.Dispatcher.Delete(serverGroupService, d, ResourceVestackServerGroup())
	if err != nil {
		return fmt.Errorf("error on deleting serverGroup %q, %w", d.Id(), err)
	}
	return err
}
