package ecs_instance_state_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_instance_state"
	"testing"
)

const testAccVestackEcsInstanceStateCreateConfig = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
	vpc_name   = "acc-test-vpc"
  	cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
  	subnet_name = "acc-test-subnet"
  	cidr_block = "172.16.0.0/24"
  	zone_id = "${data.vestack_zones.foo.zones[0].id}"
	vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_security_group" "foo" {
  	security_group_name = "acc-test-security-group"
  	vpc_id = "${vestack_vpc.foo.id}"
}

data "vestack_images" "foo" {
  	os_type = "Linux"
  	visibility = "public"
  	instance_type_id = "ecs.g1.large"
}

resource "vestack_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
  	image_id = "${data.vestack_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1.large"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	subnet_id = "${vestack_subnet.foo.id}"
	security_group_ids = ["${vestack_security_group.foo.id}"]
}

resource "vestack_ecs_instance_state" "foo" {
  	instance_id = "${vestack_ecs_instance.foo.id}"
  	action = "Stop"
  	stopped_mode = "KeepCharging"
}
`

func TestAccVestackEcsInstanceStateResource_Basic(t *testing.T) {
	resourceName := "vestack_ecs_instance_state.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance_state.VestackInstanceStateService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEcsInstanceStateCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "stopped_mode", "KeepCharging"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "STOPPED"),
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

const testAccVestackEcsInstanceStateUpdateConfig = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
	vpc_name   = "acc-test-vpc"
  	cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
  	subnet_name = "acc-test-subnet"
  	cidr_block = "172.16.0.0/24"
  	zone_id = "${data.vestack_zones.foo.zones[0].id}"
	vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_security_group" "foo" {
  	security_group_name = "acc-test-security-group"
  	vpc_id = "${vestack_vpc.foo.id}"
}

data "vestack_images" "foo" {
  	os_type = "Linux"
  	visibility = "public"
  	instance_type_id = "ecs.g1.large"
}

resource "vestack_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
  	image_id = "${data.vestack_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1.large"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	subnet_id = "${vestack_subnet.foo.id}"
	security_group_ids = ["${vestack_security_group.foo.id}"]
}

resource "vestack_ecs_instance_state" "foo" {
  	instance_id = "${vestack_ecs_instance.foo.id}"
  	action = "Start"
}
`

func TestAccVestackEcsInstanceStateResource_Update(t *testing.T) {
	resourceName := "vestack_ecs_instance_state.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance_state.VestackInstanceStateService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEcsInstanceStateCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "stopped_mode", "KeepCharging"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "STOPPED"),
				),
			},
			{
				Config: testAccVestackEcsInstanceStateUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
				),
			},
			{
				Config:             testAccVestackEcsInstanceStateUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
