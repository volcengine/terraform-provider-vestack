package access_log

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
The AccessLog is not support import.

*/

func ResourceVestackAccessLog() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackAccessLogCreate,
		Read:   resourceVestackAccessLogRead,
		Delete: resourceVestackAccessLogDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the CLB instance.",
			},
			"delivery_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tos",
				ForceNew:    true,
				Description: "The type of log delivery. Valid values: 'tos', 'tls'. Default: 'tos'.",
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the TOS bucket for storing access logs. Required when delivery_type is 'tos'.",
			},
			"tls_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the TLS project. Required when delivery_type is 'tls'.",
			},
			"tls_topic_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the TLS topic. Required when delivery_type is 'tls'.",
			},
		},
	}
	return resource
}

func resourceVestackAccessLogCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewAccessLogService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackAccessLog())
	if err != nil {
		return fmt.Errorf("error on creating access_log %q, %s", d.Id(), err)
	}
	return resourceVestackAccessLogRead(d, meta)
}

func resourceVestackAccessLogRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewAccessLogService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackAccessLog())
	if err != nil {
		return fmt.Errorf("error on reading access_log %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackAccessLogDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewAccessLogService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackAccessLog())
	if err != nil {
		return fmt.Errorf("error on deleting access_log %q, %s", d.Id(), err)
	}
	return err
}
