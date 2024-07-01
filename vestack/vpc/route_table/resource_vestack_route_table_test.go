package route_table_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/route_table"
	"testing"
)

const testAccRouteTableForCreate = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "vestack_route_table" "foo" {
  vpc_id = "${vestack_vpc.foo.id}"
  route_table_name = "acc-test-route-table"
}
`

func TestAccVestackRouteTableResource_Basic(t *testing.T) {
	resourceName := "vestack_route_table.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &route_table.VestackRouteTableService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableForCreate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "route_table_name", "acc-test-route-table"),
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

const testAccRouteTableForUpdate = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "vestack_route_table" "foo" {
  vpc_id = "${vestack_vpc.foo.id}"
  route_table_name = "acc-test-route-table-new"
  description = "tfdesc"
}
`

func TestAccVestackRouteTableResource_Update(t *testing.T) {
	resourceName := "vestack_route_table.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &route_table.VestackRouteTableService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableForCreate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "route_table_name", "acc-test-route-table"),
				),
			},
			{
				Config: testAccRouteTableForUpdate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "tfdesc"),
					resource.TestCheckResourceAttr(acc.ResourceId, "route_table_name", "acc-test-route-table-new"),
				),
			},
			{
				Config:             testAccRouteTableForUpdate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
