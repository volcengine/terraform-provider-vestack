package volume_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/ebs/volume"
	"testing"
)

const testAccVestackVolumesDatasourceConfig = `
data "vestack_zones" "foo"{
}

resource "vestack_volume" "foo" {
	volume_name = "acc-test-volume-${count.index}"
    volume_type = "ESSD_PL0"
	description = "acc-test"
    kind = "data"
    size = 60
    zone_id = "${data.vestack_zones.foo.zones[0].id}"
	volume_charge_type = "PostPaid"
	project_name = "default"
	count = 3
}

data "vestack_volumes" "foo"{
    ids = vestack_volume.foo[*].id
}
`

func TestAccVestackVolumesDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_volumes.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &volume.VestackVolumeService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackVolumesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.#", "3"),
				),
			},
		},
	})
}
