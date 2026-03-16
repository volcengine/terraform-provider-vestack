package vpc_cidr_block_associate

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
The VpcCidrBlockAssociate is not support import.

*/

func ResourceVestackVpcCidrBlockAssociate() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackVpcCidrBlockAssociateCreate,
		Read:   resourceVestackVpcCidrBlockAssociateRead,
		Delete: resourceVestackVpcCidrBlockAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the VPC.",
			},
			"secondary_cidr_block": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The secondary cidr block of the VPC.",
			},
		},
	}
	return resource
}

func resourceVestackVpcCidrBlockAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewVpcCidrBlockAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackVpcCidrBlockAssociate())
	if err != nil {
		return fmt.Errorf("error on creating vpc_cidr_block_associate %q, %s", d.Id(), err)
	}
	return resourceVestackVpcCidrBlockAssociateRead(d, meta)
}

func resourceVestackVpcCidrBlockAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewVpcCidrBlockAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackVpcCidrBlockAssociate())
	if err != nil {
		return fmt.Errorf("error on reading vpc_cidr_block_associate %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackVpcCidrBlockAssociateUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewVpcCidrBlockAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackVpcCidrBlockAssociate())
	if err != nil {
		return fmt.Errorf("error on updating vpc_cidr_block_associate %q, %s", d.Id(), err)
	}
	return resourceVestackVpcCidrBlockAssociateRead(d, meta)
}

func resourceVestackVpcCidrBlockAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewVpcCidrBlockAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackVpcCidrBlockAssociate())
	if err != nil {
		return fmt.Errorf("error on deleting vpc_cidr_block_associate %q, %s", d.Id(), err)
	}
	return err
}
