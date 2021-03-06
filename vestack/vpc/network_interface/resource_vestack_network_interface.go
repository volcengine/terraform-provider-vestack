package network_interface

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Network interface can be imported using the id, e.g.
```
$ terraform import vestack_network_interface.default eni-bp1fgnh68xyz9****
```

*/

func ResourceVestackNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Delete: resourceVestackNetworkInterfaceDelete,
		Create: resourceVestackNetworkInterfaceCreate,
		Read:   resourceVestackNetworkInterfaceRead,
		Update: resourceVestackNetworkInterfaceUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the subnet to which the ENI is connected.",
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The list of the security group id to which the secondary ENI belongs.",
			},
			"primary_ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The primary IP address of the ENI.",
				ValidateFunc: validation.IsIPAddress,
			},
			"network_interface_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the ENI.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the ENI.",
			},
			"port_security_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Set port security enable or disable.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the ENI.",
			},
		},
	}
}

func resourceVestackNetworkInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	networkInterfaceService := NewNetworkInterfaceService(meta.(*ve.SdkClient))
	if err := networkInterfaceService.Dispatcher.Create(networkInterfaceService, d, ResourceVestackNetworkInterface()); err != nil {
		return fmt.Errorf("error on creating network interface  %q, %w", d.Id(), err)
	}
	return resourceVestackNetworkInterfaceRead(d, meta)
}

func resourceVestackNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	networkInterfaceService := NewNetworkInterfaceService(meta.(*ve.SdkClient))
	if err := networkInterfaceService.Dispatcher.Read(networkInterfaceService, d, ResourceVestackNetworkInterface()); err != nil {
		return fmt.Errorf("error on reading network interface %q, %w", d.Id(), err)
	}
	return nil
}

func resourceVestackNetworkInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	networkInterfaceService := NewNetworkInterfaceService(meta.(*ve.SdkClient))
	if err := networkInterfaceService.Dispatcher.Update(networkInterfaceService, d, ResourceVestackNetworkInterface()); err != nil {
		return fmt.Errorf("error on updating network interface %q, %w", d.Id(), err)
	}
	return resourceVestackNetworkInterfaceRead(d, meta)
}

func resourceVestackNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	networkInterfaceService := NewNetworkInterfaceService(meta.(*ve.SdkClient))
	if err := networkInterfaceService.Dispatcher.Delete(networkInterfaceService, d, ResourceVestackNetworkInterface()); err != nil {
		return fmt.Errorf("error on deleting network interface %q, %w", d.Id(), err)
	}
	return nil
}
