package cr_registry

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
CR Instance can be imported using the name, e.g.
```
$ terraform import vestack_cr_instance.default enterprise-x
```

*/

func ResourceVestackCrRegistry() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackCrRegistryCreate,
		Read:   resourceVestackCrRegistryRead,
		Update: resourceVestackCrRegistryUpdate,
		Delete: resourceVestackCrRegistryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of registry.",
			},
			"delete_immediately": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether delete registry immediately. Only effected in delete action.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password of registry user.",
			},
		},
	}
	dataSource := DataSourceVestackCrRegistries().Schema["registries"].Elem.(*schema.Resource).Schema
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceVestackCrRegistryCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRegistryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceVestackCrRegistry())
	if err != nil {
		return fmt.Errorf("error on creating CrRegistry %q,%s", d.Id(), err)
	}
	return resourceVestackCrRegistryRead(d, meta)
}

func resourceVestackCrRegistryUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRegistryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceVestackCrRegistry())
	if err != nil {
		return fmt.Errorf("error on updating CrRegistry  %q, %s", d.Id(), err)
	}
	return resourceVestackCrRegistryRead(d, meta)
}

func resourceVestackCrRegistryDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRegistryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceVestackCrRegistry())
	if err != nil {
		return fmt.Errorf("error on deleting CrRegistry %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackCrRegistryRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRegistryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceVestackCrRegistry())
	if err != nil {
		return fmt.Errorf("Error on reading CrRegistry %q,%s", d.Id(), err)
	}
	return err
}
