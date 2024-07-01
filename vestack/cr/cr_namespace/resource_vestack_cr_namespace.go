package cr_namespace

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
CR namespace can be imported using the registry:name, e.g.
```
$ terraform import vestack_cr_namespace.default cr-basic:namespace-1
```

*/

func crNamespaceImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(d.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{d}, fmt.Errorf("the format of import id must be 'registry:namespace'")
	}
	if err := d.Set("registry", items[0]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	if err := d.Set("name", items[1]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}

func ResourceVestackCrNamespace() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackCrNamespaceCreate,
		Read:   resourceVestackCrNamespaceRead,
		Delete: resourceVestackCrNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: crNamespaceImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"registry": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The registry name.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of CrNamespace.",
			},
		},
	}
	dataSource := DataSourceVestackCrNamespaces().Schema["namespaces"].Elem.(*schema.Resource).Schema
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceVestackCrNamespaceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrNamespaceService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceVestackCrNamespace())
	if err != nil {
		return fmt.Errorf("error on creating CrNamespace %q,%s", d.Id(), err)
	}
	return resourceVestackCrNamespaceRead(d, meta)
}

func resourceVestackCrNamespaceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrNamespaceService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceVestackCrNamespace())
	if err != nil {
		return fmt.Errorf("error on deleting CrNamespace %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackCrNamespaceRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrNamespaceService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceVestackCrNamespace())
	if err != nil {
		return fmt.Errorf("error on reading CrNamespace %q,%s", d.Id(), err)
	}
	return err
}
