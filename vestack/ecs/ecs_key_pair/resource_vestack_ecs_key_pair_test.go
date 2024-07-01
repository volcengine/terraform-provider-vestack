package ecs_key_pair_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_key_pair"
	"testing"
)

const testAccVestackEcsKeyPairCreateConfig = `
resource "vestack_ecs_key_pair" "foo" {
  key_pair_name = "acc-test-key-name"
  description ="acc-test"
}
`

const testAccVestackEcsKeyPairUpdateConfig = `
resource "vestack_ecs_key_pair" "foo" {
    description = "acc-test-2"
    key_pair_name = "acc-test-key-name"
}
`

func TestAccVestackEcsKeyPairResource_Basic(t *testing.T) {
	resourceName := "vestack_ecs_key_pair.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_key_pair.VestackEcsKeyPairService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEcsKeyPairCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", "acc-test-key-name"),
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

func TestAccVestackEcsKeyPairResource_Update(t *testing.T) {
	resourceName := "vestack_ecs_key_pair.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_key_pair.VestackEcsKeyPairService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEcsKeyPairCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", "acc-test-key-name"),
				),
			},
			{
				Config: testAccVestackEcsKeyPairUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", "acc-test-key-name"),
				),
			},
			{
				Config:             testAccVestackEcsKeyPairUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
