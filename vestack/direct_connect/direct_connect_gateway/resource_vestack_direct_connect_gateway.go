package direct_connect_gateway

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
DirectConnectGateway can be imported using the id, e.g.
```
$ terraform import vestack_direct_connect_gateway.default resource_id
```

*/

func ResourceVestackDirectConnectGateway() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackDirectConnectGatewayCreate,
		Read:   resourceVestackDirectConnectGatewayRead,
		Update: resourceVestackDirectConnectGatewayUpdate,
		Delete: resourceVestackDirectConnectGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"direct_connect_gateway_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of direct connect gateway.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The account id of direct connect gateway.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of direct connect gateway.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The direct connect gateway tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The tag value.",
						},
					},
				},
			},
		},
	}
	return resource
}

func resourceVestackDirectConnectGatewayCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectGatewayService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackDirectConnectGateway())
	if err != nil {
		return fmt.Errorf("error on creating direct_connect_gateway %q, %s", d.Id(), err)
	}
	return resourceVestackDirectConnectGatewayRead(d, meta)
}

func resourceVestackDirectConnectGatewayRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectGatewayService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackDirectConnectGateway())
	if err != nil {
		return fmt.Errorf("error on reading direct_connect_gateway %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackDirectConnectGatewayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectGatewayService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackDirectConnectGateway())
	if err != nil {
		return fmt.Errorf("error on updating direct_connect_gateway %q, %s", d.Id(), err)
	}
	return resourceVestackDirectConnectGatewayRead(d, meta)
}

func resourceVestackDirectConnectGatewayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectGatewayService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackDirectConnectGateway())
	if err != nil {
		return fmt.Errorf("error on deleting direct_connect_gateway %q, %s", d.Id(), err)
	}
	return err
}
