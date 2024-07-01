package route_entry_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/route_entry"
	"testing"
)

const testAccRouteEntryForCreate = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc-rn"
  cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet-rn"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.vestack_zones.foo.zones[0].id}"
  vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_nat_gateway" "foo" {
  vpc_id = "${vestack_vpc.foo.id}"
  subnet_id = "${vestack_subnet.foo.id}"
  spec = "Small"
  nat_gateway_name = "acc-test-nat-rn"
}

resource "vestack_route_table" "foo" {
  vpc_id = "${vestack_vpc.foo.id}"
  route_table_name = "acc-test-route-table"
}

resource "vestack_route_entry" "foo" {
  route_table_id = "${vestack_route_table.foo.id}"
  destination_cidr_block = "172.16.1.0/24"
  next_hop_type = "NatGW"
  next_hop_id = "${vestack_nat_gateway.foo.id}"
  route_entry_name = "acc-test-route-entry"
}
`

func TestAccVestackRouteEntryResource_Basic(t *testing.T) {
	resourceName := "vestack_route_entry.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &route_entry.VestackRouteEntryService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteEntryForCreate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "route_entry_name", "acc-test-route-entry"),
					// compute status
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

const testAccRouteEntryForUpdate = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc-rn"
  cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet-rn"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.vestack_zones.foo.zones[0].id}"
  vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_nat_gateway" "foo" {
  vpc_id = "${vestack_vpc.foo.id}"
  subnet_id = "${vestack_subnet.foo.id}"
  spec = "Small"
  nat_gateway_name = "acc-test-nat-rn"
}

resource "vestack_route_table" "foo" {
  vpc_id = "${vestack_vpc.foo.id}"
  route_table_name = "acc-test-route-table"
}

resource "vestack_route_entry" "foo" {
  route_table_id = "${vestack_route_table.foo.id}"
  destination_cidr_block = "172.16.1.0/24"
  next_hop_type = "NatGW"
  next_hop_id = "${vestack_nat_gateway.foo.id}"
  route_entry_name = "acc-test-route-entry-new"
  description = "tfdesc"
}
`

func TestAccVestackRouteEntryResource_Update(t *testing.T) {
	resourceName := "vestack_route_entry.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &route_entry.VestackRouteEntryService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteEntryForCreate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "route_entry_name", "acc-test-route-entry"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config: testAccRouteEntryForUpdate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "tfdesc"),
					resource.TestCheckResourceAttr(acc.ResourceId, "route_entry_name", "acc-test-route-entry-new"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config:             testAccRouteEntryForUpdate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
