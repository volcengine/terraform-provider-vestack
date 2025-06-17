package vpc

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
VPC can be imported using the id, e.g.
```
$ terraform import vestack_vpc.default vpc-mizl7m1kqccg5smt1bdpijuj
```

*/

func ResourceVestackVpc() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackVpcCreate,
		Read:   resourceVestackVpcRead,
		Update: resourceVestackVpcUpdate,
		Delete: resourceVestackVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "A network address block which should be a subnet of the three internal network segments (10.0.0.0/16, 172.16.0.0/12 and 192.168.0.0/16).",
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the VPC.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the VPC.",
			},
			"dns_servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPAddress,
				},
				Set:         schema.HashString,
				Description: "The DNS server list of the VPC. And you can specify 0 to 5 servers to this list.",
			},
			"enable_ipv6": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether to enable the IPv6 CIDR block of the VPC.",
			},
			"ipv6_cidr_block_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The IPv6 CIDR block type of the VPC..",
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.HasChange("enable_ipv6") && d.Get("enable_ipv6").(bool) {
						return false
					}
					return true
				},
				Description: "The IPv6 CIDR block of the VPC.",
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				//ForceNew:    true,
				Description: "The ProjectName of the VPC.",
			},
			"tags": bp.TagsSchema(),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of VPC.",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of VPC.",
			},
		},
	}
	bp.MergeDateSourceToResource(DataSourceVestackVpcs().Schema["vpcs"].Elem.(*schema.Resource).Schema, &resource.Schema)
	return resource
}

func resourceVestackVpcCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := NewVpcService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(vpcService, d, ResourceVestackVpc())
	if err != nil {
		return fmt.Errorf("error on creating vpc  %q, %s", d.Id(), err)
	}
	return resourceVestackVpcRead(d, meta)
}

func resourceVestackVpcRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := NewVpcService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(vpcService, d, ResourceVestackVpc())
	if err != nil {
		return fmt.Errorf("error on reading vpc %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackVpcUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := NewVpcService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(vpcService, d, ResourceVestackVpc())
	if err != nil {
		return fmt.Errorf("error on updating vpc  %q, %s", d.Id(), err)
	}
	return resourceVestackVpcRead(d, meta)
}

func resourceVestackVpcDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := NewVpcService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(vpcService, d, ResourceVestackVpc())
	if err != nil {
		return fmt.Errorf("error on deleting vpc %q, %s", d.Id(), err)
	}
	return err
}
