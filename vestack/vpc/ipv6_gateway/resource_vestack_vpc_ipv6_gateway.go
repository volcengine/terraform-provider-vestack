package ipv6_gateway

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Ipv6Gateway can be imported using the id, e.g.
```
$ terraform import vestack_vpc_ipv6_gateway.default ipv6gw-12bcapllb5ukg17q7y2sd3thx
```

*/

func ResourceVestackIpv6Gateway() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackIpv6GatewayCreate,
		Read:   resourceVestackIpv6GatewayRead,
		Update: resourceVestackIpv6GatewayUpdate,
		Delete: resourceVestackIpv6GatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the VPC which the Ipv6Gateway belongs to.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the Ipv6Gateway.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the Ipv6Gateway.",
			},
		},
	}
	dataSource := DataSourceVestackIpv6Gateways().Schema["ipv6_gateways"].Elem.(*schema.Resource).Schema
	delete(dataSource, "id")
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceVestackIpv6GatewayCreate(d *schema.ResourceData, meta interface{}) (err error) {
	ipv6GatewayService := NewIpv6GatewayService(meta.(*bp.SdkClient))
	err = ipv6GatewayService.Dispatcher.Create(ipv6GatewayService, d, ResourceVestackIpv6Gateway())
	if err != nil {
		return fmt.Errorf("error on creating Ipv6Gateway %q, %w", d.Id(), err)
	}
	return resourceVestackIpv6GatewayRead(d, meta)
}

func resourceVestackIpv6GatewayRead(d *schema.ResourceData, meta interface{}) (err error) {
	ipv6GatewayService := NewIpv6GatewayService(meta.(*bp.SdkClient))
	err = ipv6GatewayService.Dispatcher.Read(ipv6GatewayService, d, ResourceVestackIpv6Gateway())
	if err != nil {
		return fmt.Errorf("error on reading Ipv6Gateway %q, %w", d.Id(), err)
	}
	return err
}

func resourceVestackIpv6GatewayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	ipv6GatewayService := NewIpv6GatewayService(meta.(*bp.SdkClient))
	err = ipv6GatewayService.Dispatcher.Update(ipv6GatewayService, d, ResourceVestackIpv6Gateway())
	if err != nil {
		return fmt.Errorf("error on updating Ipv6Gateway %q, %w", d.Id(), err)
	}
	return resourceVestackIpv6GatewayRead(d, meta)
}

func resourceVestackIpv6GatewayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	ipv6GatewayService := NewIpv6GatewayService(meta.(*bp.SdkClient))
	err = ipv6GatewayService.Dispatcher.Delete(ipv6GatewayService, d, ResourceVestackIpv6Gateway())
	if err != nil {
		return fmt.Errorf("error on deleting Ipv6Gateway %q, %w", d.Id(), err)
	}
	return err
}
