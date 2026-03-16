package tos_bucket_request_payment

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
TosBucketRequestPayment can be imported using the bucketName, e.g.
```
$ terraform import vestack_tos_bucket_request_payment.default bucket_name
```

*/

func ResourceVestackTosBucketRequestPayment() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackTosBucketRequestPaymentCreate,
		Read:   resourceVestackTosBucketRequestPaymentRead,
		Delete: resourceVestackTosBucketRequestPaymentDelete,
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

func resourceVestackTosBucketRequestPaymentCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketRequestPaymentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackTosBucketRequestPayment())
	if err != nil {
		return fmt.Errorf("error on creating tos_bucket_request_payment %q, %s", d.Id(), err)
	}
	return resourceVestackTosBucketRequestPaymentRead(d, meta)
}

func resourceVestackTosBucketRequestPaymentRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketRequestPaymentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackTosBucketRequestPayment())
	if err != nil {
		return fmt.Errorf("error on reading tos_bucket_request_payment %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackTosBucketRequestPaymentDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketRequestPaymentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackTosBucketRequestPayment())
	if err != nil {
		return fmt.Errorf("error on deleting tos_bucket_request_payment %q, %s", d.Id(), err)
	}
	return nil
}
