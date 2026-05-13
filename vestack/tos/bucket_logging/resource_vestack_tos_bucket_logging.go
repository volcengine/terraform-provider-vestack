package tos_bucket_logging

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
TosBucketLogging can be imported using the bucketName, e.g.
```
$ terraform import vestack_tos_bucket_logging.default bucket_name
```

*/

func ResourceVestackTosBucketLogging() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackTosBucketLoggingCreate,
		Read:   resourceVestackTosBucketLoggingRead,
		Update: resourceVestackTosBucketLoggingUpdate,
		Delete: resourceVestackTosBucketLoggingDelete,
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
			"logging_enabled": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The name of the TOS bucket.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_bucket": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the target bucket where the access logs are stored.",
						},
						"target_prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The prefix for the log object keys.",
						},
						"role": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The role that is assumed by TOS to write log objects to the target bucket.",
						},
					},
				},
			},
		},
	}
	return resource
}

func resourceVestackTosBucketLoggingCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketLoggingService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackTosBucketLogging())
	if err != nil {
		return fmt.Errorf("error on creating tos_bucket_logging %q, %s", d.Id(), err)
	}
	return resourceVestackTosBucketLoggingRead(d, meta)
}

func resourceVestackTosBucketLoggingRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketLoggingService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackTosBucketLogging())
	if err != nil {
		return fmt.Errorf("error on reading tos_bucket_logging %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackTosBucketLoggingUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketLoggingService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackTosBucketLogging())
	if err != nil {
		return fmt.Errorf("error on updating tos_bucket_logging %q, %s", d.Id(), err)
	}
	return resourceVestackTosBucketLoggingRead(d, meta)
}

func resourceVestackTosBucketLoggingDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketLoggingService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackTosBucketLogging())
	if err != nil {
		return fmt.Errorf("error on deleting tos_bucket_logging %q, %s", d.Id(), err)
	}
	return nil
}
