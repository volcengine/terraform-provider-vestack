package clb_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/clb"
	"testing"
)

const testAccVestackClbsDatasourceConfig = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
	vpc_name   = "acc-test-vpc"
	cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
	subnet_name = "acc-test-subnet"
	cidr_block = "172.16.0.0/24"
	zone_id = "${data.vestack_zones.foo.zones[0].id}"
	vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_clb" "foo" {
	type = "public"
  	subnet_id = "${vestack_subnet.foo.id}"
  	load_balancer_spec = "small_1"
  	description = "acc-test-demo"
  	load_balancer_name = "acc-test-clb-${count.index}"
	load_balancer_billing_type = "PostPaid"
  	eip_billing_config {
    	isp = "BGP"
    	eip_billing_type = "PostPaidByBandwidth"
    	bandwidth = 1
  	}
	tags {
		key = "k1"
		value = "v1"
	}
	count = 3
}

data "vestack_clbs" "foo"{
    ids = vestack_clb.foo[*].id
}
`

func TestAccVestackClbsDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_clbs.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &clb.VestackClbService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackClbsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "clbs.#", "3"),
				),
			},
		},
	})
}
