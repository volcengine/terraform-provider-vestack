package route_table_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/route_table"
	"testing"
)

const testAccRouteTableDatasourceConfig = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "vestack_route_table" "foo" {
  vpc_id = "${vestack_vpc.foo.id}"
  route_table_name = "acc-test-route-table"
  count = 3
}

data "vestack_route_tables" "foo" {
  ids = ["${vestack_route_table.foo[0].id}", "${vestack_route_table.foo[1].id}", "${vestack_route_table.foo[2].id}"]
}
`

func TestAccVestackRouteTableDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_route_tables.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &route_table.VestackRouteTableService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "route_tables.#", "3"),
				),
			},
		},
	})
}
