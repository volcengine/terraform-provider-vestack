package nat_ip

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
NatIp can be imported using the id, e.g.
```
$ terraform import vestack_nat_ip.default resource_id
```

*/

func ResourceVestackNatIp() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackNatIpCreate,
		Read:   resourceVestackNatIpRead,
		Update: resourceVestackNatIpUpdate,
		Delete: resourceVestackNatIpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"nat_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the nat gateway to which the Nat Ip belongs.",
			},
			"nat_ip_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the Nat Ip.",
			},
			"nat_ip_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the Nat Ip.",
			},
			"nat_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ip address of the Nat Ip.",
			},

			// computed fields
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the Nat Ip.",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the Ip is the default Nat Ip.",
			},
			"using_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The using status of the Nat Ip.",
			},
		},
	}
	return resource
}

func resourceVestackNatIpCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewNatIpService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackNatIp())
	if err != nil {
		return fmt.Errorf("error on creating nat_ip %q, %s", d.Id(), err)
	}
	return resourceVestackNatIpRead(d, meta)
}

func resourceVestackNatIpRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewNatIpService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackNatIp())
	if err != nil {
		return fmt.Errorf("error on reading nat_ip %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackNatIpUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewNatIpService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackNatIp())
	if err != nil {
		return fmt.Errorf("error on updating nat_ip %q, %s", d.Id(), err)
	}
	return resourceVestackNatIpRead(d, meta)
}

func resourceVestackNatIpDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewNatIpService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackNatIp())
	if err != nil {
		return fmt.Errorf("error on deleting nat_ip %q, %s", d.Id(), err)
	}
	return err
}
