package nat_gateway

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
NatGateway can be imported using the id, e.g.
```
$ terraform import vestack_nat_gateway.default ngw-vv3t043k05sm****
```

*/

func ResourceVestackNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackNatGatewayCreate,
		Read:   resourceVestackNatGatewayRead,
		Update: resourceVestackNatGatewayUpdate,
		Delete: resourceVestackNatGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the VPC.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the Subnet.",
			},
			"spec": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The specification of the NatGateway. Optional choice contains `Small`(default), `Medium`, `Large`.",
			},
			"nat_gateway_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the NatGateway.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the NatGateway.",
			},
			"billing_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "PostPaid",
				Description:  "The billing type of the NatGateway.",
				ValidateFunc: validation.StringInSlice([]string{"PostPaid"}, false),
			},
			//"period_unit": {
			//	Type:         schema.TypeString,
			//	Optional:     true,
			//	ForceNew:     true,
			//	Description:  "The period unit of the NatGateway.",
			//	ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
			//	DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			//		return d.Id() != ""
			//	},
			//},
			//"period": {
			//	Type:        schema.TypeInt,
			//	Optional:    true,
			//	ForceNew:    true,
			//	Description: "The period of the NatGateway.",
			//	DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			//		return d.Id() != ""
			//	},
			//},
		},
	}
}

func resourceVestackNatGatewayCreate(d *schema.ResourceData, meta interface{}) (err error) {
	natGatewayService := NewNatGatewayService(meta.(*ve.SdkClient))
	err = natGatewayService.Dispatcher.Create(natGatewayService, d, ResourceVestackNatGateway())
	if err != nil {
		return fmt.Errorf("error on creating nat gateway  %q, %w", d.Id(), err)
	}
	return resourceVestackNatGatewayRead(d, meta)
}

func resourceVestackNatGatewayRead(d *schema.ResourceData, meta interface{}) (err error) {
	natGatewayService := NewNatGatewayService(meta.(*ve.SdkClient))
	err = natGatewayService.Dispatcher.Read(natGatewayService, d, ResourceVestackNatGateway())
	if err != nil {
		return fmt.Errorf("error on reading nat gateway %q, %w", d.Id(), err)
	}
	return err
}

func resourceVestackNatGatewayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	natGatewayService := NewNatGatewayService(meta.(*ve.SdkClient))
	err = natGatewayService.Dispatcher.Update(natGatewayService, d, ResourceVestackNatGateway())
	if err != nil {
		return fmt.Errorf("error on updating nat gateway  %q, %w", d.Id(), err)
	}
	return resourceVestackNatGatewayRead(d, meta)
}

func resourceVestackNatGatewayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	natGatewayService := NewNatGatewayService(meta.(*ve.SdkClient))
	err = natGatewayService.Dispatcher.Delete(natGatewayService, d, ResourceVestackNatGateway())
	if err != nil {
		return fmt.Errorf("error on deleting nat gateway %q, %w", d.Id(), err)
	}
	return err
}
