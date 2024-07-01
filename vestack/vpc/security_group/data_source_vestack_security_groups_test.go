package security_group_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/security_group"
	"testing"
)

const testAccSecurityGroupDatasourceConfig = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "vestack_security_group" "foo" {
  vpc_id = "${vestack_vpc.foo.id}"
  count = 3
}

data "vestack_security_groups" "foo"{
  ids = ["${vestack_security_group.foo[0].id}", "${vestack_security_group.foo[1].id}", "${vestack_security_group.foo[2].id}"]
}
`

func TestAccVestackSecurityGroupDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_security_groups.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &security_group.VestackSecurityGroupService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "security_groups.#", "3"),
				),
			},
		},
	})
}
