package subnet

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Subnet can be imported using the id, e.g.
```
$ terraform import vestack_subnet.default subnet-274oj9a8rs9a87fap8sf9515b
```

*/

func ResourceVestackSubnet() *schema.Resource {
	return &schema.Resource{
		Delete: resourceVestackSubnetDelete,
		Create: resourceVestackSubnetCreate,
		Read:   resourceVestackSubnetRead,
		Update: resourceVestackSubnetUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "A network address block which should be a subnet of the three internal network segments (10.0.0.0/16, 172.16.0.0/12 and 192.168.0.0/16).",
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the Subnet.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Subnet.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of Subnet.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the VPC.",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the Zone.",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of Subnet.",
			},
		},
	}
}

func resourceVestackSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	subnetService := NewSubnetService(meta.(*ve.SdkClient))
	if err := subnetService.Dispatcher.Create(subnetService, d, ResourceVestackSubnet()); err != nil {
		return fmt.Errorf("error on creating subnet  %q, %w", d.Id(), err)
	}
	return resourceVestackSubnetRead(d, meta)
}

func resourceVestackSubnetRead(d *schema.ResourceData, meta interface{}) error {
	subnetService := NewSubnetService(meta.(*ve.SdkClient))
	if err := subnetService.Dispatcher.Read(subnetService, d, ResourceVestackSubnet()); err != nil {
		return fmt.Errorf("error on reading subnet %q, %w", d.Id(), err)
	}
	return nil
}

func resourceVestackSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	subnetService := NewSubnetService(meta.(*ve.SdkClient))
	if err := subnetService.Dispatcher.Update(subnetService, d, ResourceVestackSubnet()); err != nil {
		return fmt.Errorf("error on updating subnet %q, %w", d.Id(), err)
	}
	return resourceVestackSubnetRead(d, meta)
}

func resourceVestackSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	subnetService := NewSubnetService(meta.(*ve.SdkClient))
	if err := subnetService.Dispatcher.Delete(subnetService, d, ResourceVestackSubnet()); err != nil {
		return fmt.Errorf("error on deleting subnet %q, %w", d.Id(), err)
	}
	return nil
}
