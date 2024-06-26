package subnet_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/subnet"
	"testing"
)

const testAccSubnetDatasourceConfig = `
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

resource "vestack_subnet" "foo1" {
  subnet_name = "acc-test-subnet1"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.vestack_zones.foo.zones[0].id}"
  vpc_id = "${vestack_vpc.foo.id}"
}

data "vestack_subnets" "foo"{
  ids = ["${vestack_subnet.foo.id}", "${vestack_subnet.foo1.id}"]
}
`

func TestAccVestackSubnetDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_subnets.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &subnet.VestackSubnetService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "subnets.#", "2"),
				),
			},
		},
	})
}
