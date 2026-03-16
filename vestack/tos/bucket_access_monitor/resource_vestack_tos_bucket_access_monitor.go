package tos_bucket_access_monitor

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
TosBucketAccessMonitor can be imported using the bucketName, e.g.
```
$ terraform import vestack_tos_bucket_access_monitor.default bucket_name
```

*/

func ResourceVestackTosBucketAccessMonitor() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackTosBucketAccessMonitorCreate,
		Read:   resourceVestackTosBucketAccessMonitorRead,
		Delete: resourceVestackTosBucketAccessMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the TOS bucket.",
			},
		},
	}
	return resource
}

func resourceVestackTosBucketAccessMonitorCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketAccessMonitorService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackTosBucketAccessMonitor())
	if err != nil {
		return fmt.Errorf("error on creating tos_bucket_access_monitor %q, %s", d.Id(), err)
	}
	return resourceVestackTosBucketAccessMonitorRead(d, meta)
}

func resourceVestackTosBucketAccessMonitorRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketAccessMonitorService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackTosBucketAccessMonitor())
	if err != nil {
		return fmt.Errorf("error on reading tos_bucket_access_monitor %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackTosBucketAccessMonitorUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketAccessMonitorService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackTosBucketAccessMonitor())
	if err != nil {
		return fmt.Errorf("error on updating tos_bucket_access_monitor %q, %s", d.Id(), err)
	}
	return resourceVestackTosBucketAccessMonitorRead(d, meta)
}

func resourceVestackTosBucketAccessMonitorDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketAccessMonitorService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackTosBucketAccessMonitor())
	if err != nil {
		return fmt.Errorf("error on deleting tos_bucket_access_monitor %q, %s", d.Id(), err)
	}
	return nil
}
