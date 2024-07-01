package network_interface_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/network_interface"
)

const testAccVestackNetworkInterfacesDatasourceConfig = `
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

resource "vestack_security_group" "foo" {
  security_group_name = "acc-test-sg"
  vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_network_interface" "foo" {
  network_interface_name = "acc-test-eni-${count.index}"
  subnet_id = "${vestack_subnet.foo.id}"
  security_group_ids = ["${vestack_security_group.foo.id}"]
  count = 3
}

data "vestack_network_interfaces" "foo"{
    ids = vestack_network_interface.foo[*].id
}
`

func TestAccVestackNetworkInterfacesDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_network_interfaces.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_interface.VestackNetworkInterfaceService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackNetworkInterfacesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "network_interfaces.#", "3"),
				),
			},
		},
	})
}
