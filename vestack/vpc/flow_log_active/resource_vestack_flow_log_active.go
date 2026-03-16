package flow_log_active

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
FlowLogActive can be imported using the id, e.g.
```
$ terraform import vestack_flow_log_active.default resource_id
```

*/

func ResourceVestackFlowLogActive() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackFlowLogActiveCreate,
		Read:   resourceVestackFlowLogActiveRead,
		Delete: resourceVestackFlowLogActiveDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"flow_log_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of flow log.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of flow log.",
			},
		},
	}
	return resource
}

func resourceVestackFlowLogActiveCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewFlowLogActiveService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackFlowLogActive())
	if err != nil {
		return fmt.Errorf("error on creating flow_log_active %q, %s", d.Id(), err)
	}
	return resourceVestackFlowLogActiveRead(d, meta)
}

func resourceVestackFlowLogActiveRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewFlowLogActiveService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackFlowLogActive())
	if err != nil {
		return fmt.Errorf("error on reading flow_log_active %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackFlowLogActiveUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewFlowLogActiveService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackFlowLogActive())
	if err != nil {
		return fmt.Errorf("error on updating flow_log_active %q, %s", d.Id(), err)
	}
	return resourceVestackFlowLogActiveRead(d, meta)
}

func resourceVestackFlowLogActiveDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewFlowLogActiveService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackFlowLogActive())
	if err != nil {
		return fmt.Errorf("error on deleting flow_log_active %q, %s", d.Id(), err)
	}
	return err
}
