package ecs_available_resource_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_available_resource"
)

const testAccVestackAvailableResourcesDatasourceConfig = `
data "vestack_ecs_available_resources" "foo"{
    destination_resource = "InstanceType"
}
`

func TestAccVestackAvailableResourcesDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_ecs_available_resources.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return ecs_available_resource.NewEcsAvailableResourceService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackAvailableResourcesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "available_zones.#", "3"),
				),
			},
		},
	})
}
