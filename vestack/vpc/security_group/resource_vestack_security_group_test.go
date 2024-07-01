package security_group_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/security_group"
	"testing"
)

const testAccSecurityGroupForCreate = `
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
`

const testAccSecurityGroupForUpdate = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "vestack_security_group" "foo" {
  description = "tfdesc"
  vpc_id = "${vestack_vpc.foo.id}"
  security_group_name = "acc-test-security-group"

  tags {
    key = "k2"
    value = "v2"
  }

  tags {
    key = "k1"
    value = "v1"
  }
}
`

func TestAccVestackSecurityGroupResource_Basic(t *testing.T) {
	resourceName := "vestack_security_group.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &security_group.VestackSecurityGroupService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupForCreate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_name", "acc-test-security-group"),
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

func TestAccVestackSecurityGroupResource_Update(t *testing.T) {
	resourceName := "vestack_security_group.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &security_group.VestackSecurityGroupService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupForCreate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_name", "acc-test-security-group"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config: testAccSecurityGroupForUpdate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_name", "acc-test-security-group"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "tfdesc"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "2"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config:             testAccSecurityGroupForUpdate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
