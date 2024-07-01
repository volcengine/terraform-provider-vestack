package cr_endpoint

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
CR endpoints can be imported using the endpoint:registryName, e.g.
```
$ terraform import vestack_cr_endpoint.default endpoint:cr-basic
```

*/

func crEndpointImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(d.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{d}, fmt.Errorf("the format of import id must start with 'endpoint:',eg: 'endpoint:[registry-1]'")
	}
	if err := d.Set("registry", items[1]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}

func ResourceVestackCrEndpoint() *schema.Resource {
	resource := &schema.Resource{
		Read:   resourceVestackCrEndpointRead,
		Create: resourceVestackCrEndpointCreate,
		Update: resourceVestackCrEndpointUpdate,
		Delete: resourceVestackCrEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: crEndpointImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"registry": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The CrRegistry name.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether enable public endpoint.",
			},
		},
	}
	dataSource := DataSourceVestackCrEndpoints().Schema["endpoints"].Elem.(*schema.Resource).Schema
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceVestackCrEndpointCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrEndpointService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceVestackCrEndpoint())
	if err != nil {
		return fmt.Errorf("Error on creating CrEndpoint %q,%s", d.Id(), err)
	}
	return resourceVestackCrEndpointRead(d, meta)
}

func resourceVestackCrEndpointUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrEndpointService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceVestackCrEndpoint())
	if err != nil {
		return fmt.Errorf("error on updating CrEndpoint  %q, %s", d.Id(), err)
	}
	return resourceVestackCrEndpointRead(d, meta)
}

func resourceVestackCrEndpointDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrEndpointService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceVestackCrEndpoint())
	if err != nil {
		return fmt.Errorf("error on deleting CrEndpoint %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackCrEndpointRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrEndpointService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceVestackCrEndpoint())
	if err != nil {
		return fmt.Errorf("Error on reading CrEndpoint %q,%s", d.Id(), err)
	}
	return err
}
