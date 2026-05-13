package zone_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/zone"
	"testing"
)

const testAccVestackClbZonesDatasourceConfig = `
data "vestack_clb_zones" "foo"{
}
`

func TestAccVestackClbZonesDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_clb_zones.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &zone.VestackClbZoneService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackClbZonesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "master_zones.#", "1"),
				),
			},
		},
	})
}
