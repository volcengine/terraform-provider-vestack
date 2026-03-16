package tos_bucket_rename

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
TosBucketRename can be imported using the bucketName, e.g.
```
$ terraform import vestack_tos_bucket_rename.default bucket_name
```

*/

func ResourceVestackTosBucketRename() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackTosBucketRenameCreate,
		Read:   resourceVestackTosBucketRenameRead,
		Delete: resourceVestackTosBucketRenameDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				bucketName := d.Id()
				if err := d.Set("bucket_name", bucketName); err != nil {
					return nil, err
				}
				d.SetId(bucketName)
				return []*schema.ResourceData{d}, nil
			},
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
				Description: "The name of the TOS bucket to configure rename functionality for.",
			},
		},
	}
	return resource
}

func resourceVestackTosBucketRenameCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketRenameService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackTosBucketRename())
	if err != nil {
		return fmt.Errorf("error on creating tos_bucket_rename %q, %s", d.Id(), err)
	}
	return resourceVestackTosBucketRenameRead(d, meta)
}

func resourceVestackTosBucketRenameRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketRenameService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackTosBucketRename())
	if err != nil {
		return fmt.Errorf("error on reading tos_bucket_rename %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackTosBucketRenameUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketRenameService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackTosBucketRename())
	if err != nil {
		return fmt.Errorf("error on updating tos_bucket_rename %q, %s", d.Id(), err)
	}
	return resourceVestackTosBucketRenameRead(d, meta)
}

func resourceVestackTosBucketRenameDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketRenameService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackTosBucketRename())
	if err != nil {
		return fmt.Errorf("error on deleting tos_bucket_rename %q, %s", d.Id(), err)
	}
	return nil
}
