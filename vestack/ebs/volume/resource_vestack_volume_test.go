package volume_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/ebs/volume"
	"testing"
)

const testAccVestackVolumeCreateConfig = `
data "vestack_zones" "foo"{
}

resource "vestack_volume" "foo" {
	volume_name = "acc-test-volume"
    volume_type = "ESSD_PL0"
	description = "acc-test"
    kind = "data"
    size = 40
    zone_id = "${data.vestack_zones.foo.zones[0].id}"
	volume_charge_type = "PostPaid"
	project_name = "default"
}
`

func TestAccVestackVolumeResource_Basic(t *testing.T) {
	resourceName := "vestack_volume.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &volume.VestackVolumeService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackVolumeCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_name", "acc-test-volume"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_with_instance", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kind", "data"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_id", ""),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "created_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trade_status"),
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

const testAccVestackVolumeUpdateBasicAttributeConfig = `
data "vestack_zones" "foo"{
}

resource "vestack_volume" "foo" {
	volume_name = "acc-test-volume-new"
    volume_type = "ESSD_PL0"
	description = "acc-test-new"
    kind = "data"
    size = 40
    zone_id = "${data.vestack_zones.foo.zones[0].id}"
	volume_charge_type = "PostPaid"
	project_name = "default"
	delete_with_instance = true
}
`

func TestAccVestackVolumeResource_UpdateBasicAttribute(t *testing.T) {
	resourceName := "vestack_volume.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &volume.VestackVolumeService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackVolumeCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_name", "acc-test-volume"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_with_instance", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kind", "data"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_id", ""),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "created_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trade_status"),
				),
			},
			{
				Config: testAccVestackVolumeUpdateBasicAttributeConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_name", "acc-test-volume-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kind", "data"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_id", ""),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "created_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trade_status"),
				),
			},
			{
				Config:             testAccVestackVolumeUpdateBasicAttributeConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccVestackVolumeUpdateVolumeSizeConfig = `
data "vestack_zones" "foo"{
}

resource "vestack_volume" "foo" {
	volume_name = "acc-test-volume"
    volume_type = "ESSD_PL0"
	description = "acc-test"
    kind = "data"
    size = 60
    zone_id = "${data.vestack_zones.foo.zones[0].id}"
	volume_charge_type = "PostPaid"
	project_name = "default"
}
`

func TestAccVestackVolumeResource_UpdateVolumeSize(t *testing.T) {
	resourceName := "vestack_volume.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &volume.VestackVolumeService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackVolumeCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_name", "acc-test-volume"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_with_instance", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kind", "data"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_id", ""),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "created_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trade_status"),
				),
			},
			{
				Config: testAccVestackVolumeUpdateVolumeSizeConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_name", "acc-test-volume"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_with_instance", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kind", "data"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_id", ""),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "created_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trade_status"),
				),
			},
			{
				Config:             testAccVestackVolumeUpdateVolumeSizeConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
