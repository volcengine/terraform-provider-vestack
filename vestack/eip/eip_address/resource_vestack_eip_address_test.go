package eip_address_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/eip/eip_address"
	"testing"
)

const testAccVestackEipAddressCreateConfig = `
resource "vestack_eip_address" "foo" {
    billing_type = "PostPaidByTraffic"
}
`

const testAccVestackEipAddressUpdateConfig = `
resource "vestack_eip_address" "foo" {
    bandwidth = 1
    billing_type = "PostPaidByBandwidth"
    description = "acc-test"
    isp = "BGP"
    name = "acc-test-eip"
}
`

func TestAccVestackEipAddressResource_Basic(t *testing.T) {
	resourceName := "vestack_eip_address.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &eip_address.VestackEipAddressService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEipAddressCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "billing_type", "PostPaidByTraffic"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVestackEipAddressResource_Update(t *testing.T) {
	resourceName := "vestack_eip_address.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &eip_address.VestackEipAddressService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEipAddressCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "billing_type", "PostPaidByTraffic"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config: testAccVestackEipAddressUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "billing_type", "PostPaidByBandwidth"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "isp", "BGP"),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-eip"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config:             testAccVestackEipAddressUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
