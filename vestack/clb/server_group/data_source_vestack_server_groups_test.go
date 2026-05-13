package server_group_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/server_group"
	"testing"
)

const testAccVestackServerGroupsDatasourceConfig = `
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
  description = "acc0Demo"
  load_balancer_name = "acc-test-create"
  eip_billing_config {
    isp = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth = 1
  }
}

resource "vestack_server_group" "foo" {
  load_balancer_id = "${vestack_clb.foo.id}"
  server_group_name = "acc-test-create"
  description = "hello demo11"
}

data "vestack_server_groups" "foo"{
    ids = ["${vestack_server_group.foo.id}"]
}
`

func TestAccVestackServerGroupsDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_server_groups.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &server_group.VestackServerGroupService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackServerGroupsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "groups.#", "1"),
				),
			},
		},
	})
}
