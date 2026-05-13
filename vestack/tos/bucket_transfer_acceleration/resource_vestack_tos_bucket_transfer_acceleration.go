package tos_bucket_transfer_acceleration

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
TosBucketTransferAcceleration can be imported using the bucketName, e.g.
```
$ terraform import vestack_tos_bucket_transfer_acceleration.default bucket_name
```

*/

func ResourceVestackTosBucketTransferAcceleration() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackTosBucketTransferAccelerationCreate,
		Read:   resourceVestackTosBucketTransferAccelerationRead,
		Delete: resourceVestackTosBucketTransferAccelerationDelete,
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

func resourceVestackTosBucketTransferAccelerationCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketTransferAccelerationService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackTosBucketTransferAcceleration())
	if err != nil {
		return fmt.Errorf("error on creating tos_bucket_transfer_acceleration %q, %s", d.Id(), err)
	}
	return resourceVestackTosBucketTransferAccelerationRead(d, meta)
}

func resourceVestackTosBucketTransferAccelerationRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketTransferAccelerationService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackTosBucketTransferAcceleration())
	if err != nil {
		return fmt.Errorf("error on reading tos_bucket_transfer_acceleration %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackTosBucketTransferAccelerationUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketTransferAccelerationService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackTosBucketTransferAcceleration())
	if err != nil {
		return fmt.Errorf("error on updating tos_bucket_transfer_acceleration %q, %s", d.Id(), err)
	}
	return resourceVestackTosBucketTransferAccelerationRead(d, meta)
}

func resourceVestackTosBucketTransferAccelerationDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketTransferAccelerationService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackTosBucketTransferAcceleration())
	if err != nil {
		return fmt.Errorf("error on deleting tos_bucket_transfer_acceleration %q, %s", d.Id(), err)
	}
	return nil
}
