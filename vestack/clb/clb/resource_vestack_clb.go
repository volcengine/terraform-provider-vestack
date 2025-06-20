package clb

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
CLB can be imported using the id, e.g.
```
$ terraform import vestack_clb.default clb-273y2ok6ets007fap8txvf6us
```

*/

func ResourceVestackClb() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackClbCreate,
		Read:   resourceVestackClbRead,
		Update: resourceVestackClbUpdate,
		Delete: resourceVestackClbDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region of the request.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The type of the CLB. And optional choice contains `public` or `private`.",
				ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
			},
			"load_balancer_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the CLB.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the CLB.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "The id of the VPC.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the Subnet.",
			},
			"eni_address": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Optional:    true,
				Description: "The eni address of the CLB.",
			},
			"modification_protection_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The status of the console modification protection, the value can be `NonProtection` or `ConsoleProtection`.",
				ValidateFunc: validation.StringInSlice([]string{"NonProtection", "ConsoleProtection"}, false),
			},
			"modification_protection_reason": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The reason of the console modification protection.",
			},
			"load_balancer_spec": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("load_balancer_billing_type") == "PostPaidByLCU"
				},
				ValidateFunc: validation.StringInSlice([]string{
					"small_1", "small_2", "medium_1", "medium_2", "large_1", "large_2",
				}, false),
				Description: "The specification of the CLB, the value can be `small_1`, `small_2`, `medium_1`, `medium_2`, `large_1`, `large_2`. When the value of the `load_balancer_billing_type` is `PostPaidByLCU`, this field does not need to be specified.",
			},
			"load_balancer_billing_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PostPaid", "PrePaid", "PostPaidByLCU"}, false),
				Description:  "The billing type of the CLB, valid values: `PostPaid`, `PrePaid`, `PostPaidByLCU`. Default is `PostPaid`.",
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  12,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36})),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !(d.Get("load_balancer_billing_type").(string) == "PrePaid")
				},
				Description: "The period of the NatGateway, the valid value range in 1~9 or 12 or 24 or 36. Default value is 12. The period unit defaults to `Month`." +
					"This field is only effective when creating a PrePaid NatGateway. When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"renew_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The renew type of the CLB. When the value of the load_balancer_billing_type is `PrePaid`, the query returns this field.",
			},
			"eip_billing_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Description: "The billing configuration of the EIP which automatically associated to CLB. This field is valid when the type of CLB is `public`." +
					"When the type of the CLB is `private`, suggest using a combination of resource `vestack_eip_address` and `vestack_eip_associate` to achieve public network access function.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"isp": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The ISP of the EIP which automatically associated to CLB, the value can be `BGP` or `ChinaMobile` or `ChinaUnicom` or `ChinaTelecom` or `SingleLine_BGP` or `Static_BGP` or `Fusion_BGP`.",
						},
						"eip_billing_type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"PrePaid", "PostPaidByBandwidth", "PostPaidByTraffic"}, false),
							Description: "The billing type of the EIP which automatically assigned to CLB. And optional choice contains `PostPaidByBandwidth` or `PostPaidByTraffic` or `PrePaid`." +
								"When creating a `PrePaid` public CLB, this field must be specified as `PrePaid` simultaneously." +
								"When the LoadBalancerBillingType changes from `PostPaid` to `PrePaid`, please manually modify the value of this field to `PrePaid` simultaneously.",
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							//ValidateFunc: validation.IntBetween(1, 500),
							//Description:  "The peek bandwidth of the EIP which automatically assigned to CLB. The value range in 1~500 for PostPaidByBandwidth, and 1~200 for PostPaidByTraffic.",
							Description: "The peek bandwidth of the EIP which automatically assigned to CLB.",
						},
					},
				},
			},
			"address_ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "ipv4",
				ValidateFunc: validation.StringInSlice([]string{"ipv4", "DualStack"}, false),
				Description: "The address ip version of the Clb. Valid values: `ipv4`, `DualStack`. Default is `ipv4`.\n" +
					"When the value of this field is `DualStack`, the type of the CLB must be `private`, " +
					"and suggest using a combination of resource `vestack_vpc_ipv6_gateway` and `vestack_vpc_ipv6_address_bandwidth` to achieve ipv6 public network access function.",
			},
			"eni_ipv6_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv6Address,
				Description:  "The eni ipv6 address of the Clb.",
			},
			"eip_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Eip ID of the Clb.",
			},
			"eip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Eip address of the Clb.",
			},
			"ipv6_eip_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Ipv6 Eip ID of the Clb.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ProjectName of the CLB.",
			},
			"tags": ve.TagsSchema(),
			"master_zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The master zone ID of the CLB.",
			},
			"slave_zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The slave zone ID of the CLB.",
			},
			//"period_unit": {
			//	Type:         schema.TypeString,
			//	Optional:     true,
			//	Description:  "The period unit of PrePaid billing type.",
			//	ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
			//	DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			//		return d.Id() != ""
			//	},
			//},
			//"period": {
			//	Type:        schema.TypeInt,
			//	Optional:    true,
			//	Description: "The period of PrePaid billing type.",
			//	DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			//		return d.Id() != ""
			//	},
			//},
		},
	}
}

func resourceVestackClbCreate(d *schema.ResourceData, meta interface{}) (err error) {
	clbService := NewClbService(meta.(*ve.SdkClient))
	err = ve.DefaultDispatcher().Create(clbService, d, ResourceVestackClb())
	if err != nil {
		return fmt.Errorf("error on creating clb  %q, %w", d.Id(), err)
	}
	return resourceVestackClbRead(d, meta)
}

func resourceVestackClbRead(d *schema.ResourceData, meta interface{}) (err error) {
	clbService := NewClbService(meta.(*ve.SdkClient))
	err = ve.DefaultDispatcher().Read(clbService, d, ResourceVestackClb())
	if err != nil {
		return fmt.Errorf("error on reading clb %q, %w", d.Id(), err)
	}
	return err
}

func resourceVestackClbUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	clbService := NewClbService(meta.(*ve.SdkClient))
	err = ve.DefaultDispatcher().Update(clbService, d, ResourceVestackClb())
	if err != nil {
		return fmt.Errorf("error on updating clb  %q, %w", d.Id(), err)
	}
	return resourceVestackClbRead(d, meta)
}

func resourceVestackClbDelete(d *schema.ResourceData, meta interface{}) (err error) {
	clbService := NewClbService(meta.(*ve.SdkClient))
	err = ve.DefaultDispatcher().Delete(clbService, d, ResourceVestackClb())
	if err != nil {
		return fmt.Errorf("error on deleting clb %q, %w", d.Id(), err)
	}
	return err
}
