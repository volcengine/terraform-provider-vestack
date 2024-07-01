package security_group_rule_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/security_group_rule"
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
  security_group_name = "acc-test-security-group"
}

resource "vestack_security_group_rule" "foo" {
  direction         = "egress"
  security_group_id = "${vestack_security_group.foo.id}"
  protocol          = "tcp"
  port_start        = 8000
  port_end          = 9003
  cidr_ip           = "172.16.0.0/24"
}

data "vestack_security_group_rules" "foo"{
  security_group_id = "${vestack_security_group.foo.id}"
  direction = "${vestack_security_group_rule.foo.direction}"
  cidr_ip = "${vestack_security_group_rule.foo.cidr_ip}"
}
`

func TestAccVestackSecurityGroupRulesDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_security_group_rules.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &security_group_rule.VestackSecurityGroupRuleService{},
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
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_rules.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_rules.0.direction", "egress"),
				),
			},
		},
	})
}
