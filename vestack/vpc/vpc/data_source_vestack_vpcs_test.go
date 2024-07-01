package vpc_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/vpc"
	"testing"
)

const testAccVpcDatasourceConfig = `
resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "vestack_vpc" "foo1" {
  vpc_name   = "acc-test-vpc1"
  cidr_block = "172.16.0.0/16"
}

data "vestack_vpcs" "foo"{
  ids = ["${vestack_vpc.foo1.id}", "${vestack_vpc.foo.id}"]
}
`

func TestAccVestackVpcDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_vpcs.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &vpc.VestackVpcService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "vpcs.#", "2"),
				),
			},
		},
	})
}
