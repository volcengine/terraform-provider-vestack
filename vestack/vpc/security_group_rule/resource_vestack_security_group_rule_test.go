package security_group_rule_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/security_group_rule"
	"testing"
)

const testAccSecurityGroupRuleForCreate = `
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
`

const testAccSecurityGroupRuleForUpdate = `
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
  description       = "tfdesc"
}
`

func TestAccVestackSecurityGroupRuleResource_Basic(t *testing.T) {
	resourceName := "vestack_security_group_rule.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &security_group_rule.VestackSecurityGroupRuleService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleForCreate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "direction", "egress"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "tcp"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_start", "8000"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_end", "9003"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cidr_ip", "172.16.0.0/24"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config:             testAccSecurityGroupRuleForCreate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVestackSubnetResource_Update(t *testing.T) {
	resourceName := "vestack_security_group_rule.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &security_group_rule.VestackSecurityGroupRuleService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleForCreate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "direction", "egress"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "tcp"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_start", "8000"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_end", "9003"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cidr_ip", "172.16.0.0/24"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config: testAccSecurityGroupRuleForUpdate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "direction", "egress"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "tcp"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_start", "8000"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_end", "9003"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cidr_ip", "172.16.0.0/24"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "tfdesc"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config:             testAccSecurityGroupRuleForUpdate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
