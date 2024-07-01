package zone_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/zone"
	"testing"
)

const testAccVestackZonesDatasourceConfig = `
data "vestack_zones" "foo"{
    ids = ["cn-chengdu-a"]
}
`

func TestAccVestackZonesDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_zones.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &zone.VestackZoneService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackZonesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "zones.#", "1"),
				),
			},
		},
	})
}
