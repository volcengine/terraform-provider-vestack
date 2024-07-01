package ecs_deployment_set_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_deployment_set"
	"testing"
)

const testAccVestackEcsDeploymentSetCreateConfig = `
resource "vestack_ecs_deployment_set" "foo" {
    deployment_set_name = "acc-test-ecs-ds"
	description = "acc-test"
    granularity = "switch"
    strategy = "Availability"
}
`

func TestAccVestackEcsDeploymentSetResource_Basic(t *testing.T) {
	resourceName := "vestack_ecs_deployment_set.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_deployment_set.VestackEcsDeploymentSetService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEcsDeploymentSetCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_name", "acc-test-ecs-ds"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "granularity", "switch"),
					resource.TestCheckResourceAttr(acc.ResourceId, "strategy", "Availability"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
		},
	})
}

const testAccVestackEcsDeploymentSetUpdateConfig = `
resource "vestack_ecs_deployment_set" "foo" {
    deployment_set_name = "acc-test-ecs-ds-new"
	description = "acc-test"
    granularity = "switch"
    strategy = "Availability"
}
`

func TestAccVestackEcsDeploymentSetResource_Update(t *testing.T) {
	resourceName := "vestack_ecs_deployment_set.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_deployment_set.VestackEcsDeploymentSetService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEcsDeploymentSetCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_name", "acc-test-ecs-ds"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "granularity", "switch"),
					resource.TestCheckResourceAttr(acc.ResourceId, "strategy", "Availability"),
				),
			},
			{
				Config: testAccVestackEcsDeploymentSetUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_name", "acc-test-ecs-ds-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "granularity", "switch"),
					resource.TestCheckResourceAttr(acc.ResourceId, "strategy", "Availability"),
				),
			},
			{
				Config:             testAccVestackEcsDeploymentSetUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
