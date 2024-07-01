package ecs_deployment_set_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_deployment_set"
	"testing"
)

const testAccVestackEcsDeploymentSetsDatasourceConfig = `
resource "vestack_ecs_deployment_set" "foo" {
    deployment_set_name = "acc-test-ecs-ds-${count.index}"
	description = "acc-test"
    granularity = "switch"
    strategy = "Availability"
	count = 3
}

data "vestack_ecs_deployment_sets" "foo"{
    granularity = "switch"
    ids = vestack_ecs_deployment_set.foo[*].id
}
`

func TestAccVestackEcsDeploymentSetsDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_ecs_deployment_sets.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_deployment_set.VestackEcsDeploymentSetService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEcsDeploymentSetsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_sets.#", "3"),
				),
			},
		},
	})
}
