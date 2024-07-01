package route_table_associate_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/route_table_associate"
	"testing"
)

const testAccRouteTableAssociateForCreate = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc-attach"
  cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet-attach"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.vestack_zones.foo.zones[0].id}"
  vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_subnet" "foo1" {
  subnet_name = "acc-test-subnet-attach1"
  cidr_block = "172.16.16.0/20"
  zone_id = "${data.vestack_zones.foo.zones[0].id}"
  vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_subnet" "foo2" {
  subnet_name = "acc-test-subnet-attach2"
  cidr_block = "172.16.6.0/23"
  zone_id = "${data.vestack_zones.foo.zones[0].id}"
  vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_subnet" "foo3" {
  subnet_name = "acc-test-subnet-attach3"
  cidr_block = "172.16.14.0/26"
  zone_id = "${data.vestack_zones.foo.zones[0].id}"
  vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_route_table" "foo" {
  vpc_id = "${vestack_vpc.foo.id}"
  route_table_name = "acc-test-route-table-attach"
}

resource "vestack_route_table" "foo1" {
  vpc_id = "${vestack_vpc.foo.id}"
  route_table_name = "acc-test-route-table-attach1"
}

resource "vestack_route_table" "foo2" {
  vpc_id = "${vestack_vpc.foo.id}"
  route_table_name = "acc-test-route-table-attach2"
}

resource "vestack_route_table" "foo3" {
  vpc_id = "${vestack_vpc.foo.id}"
  route_table_name = "acc-test-route-table-attach3"
}

resource "vestack_route_table_associate" "foo" {
  route_table_id = "${vestack_route_table.foo.id}"
  subnet_id = "${vestack_subnet.foo.id}"
}

resource "vestack_route_table_associate" "foo1" {
  route_table_id = "${vestack_route_table.foo1.id}"
  subnet_id = "${vestack_subnet.foo1.id}"
}

resource "vestack_route_table_associate" "foo2" {
  route_table_id = "${vestack_route_table.foo2.id}"
  subnet_id = "${vestack_subnet.foo2.id}"
}

resource "vestack_route_table_associate" "foo3" {
  route_table_id = "${vestack_route_table.foo3.id}"
  subnet_id = "${vestack_subnet.foo3.id}"
}

`

func TestAccVestackRouteTableAssociateResource_Basic(t *testing.T) {
	resourceName := "vestack_route_table_associate.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &route_table_associate.VestackRouteTableAssociateService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableAssociateForCreate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
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
