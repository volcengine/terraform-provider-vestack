package eip_associate_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/eip/eip_associate"
	"testing"
)

const testAccVestackEipAssociateCreateConfig = `
	data "vestack_zones" "foo"{
	}

	data "vestack_images" "foo" {
	  os_type = "Linux"
	  visibility = "public"
	  instance_type_id = "ecs.g1.large"
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
	  vpc_id = "${vestack_vpc.foo.id}"
	  security_group_name = "acc-test-security-group"
	}

	resource "vestack_ecs_instance" "foo" {
	  image_id = "${data.vestack_images.foo.images[0].image_id}"
	  instance_type = "ecs.g1.large"
	  instance_name = "acc-test-ecs-name"
	  password = "93f0cb0614Aab12"
	  instance_charge_type = "PostPaid"
	  system_volume_type = "ESSD_PL0"
	  system_volume_size = 40
	  subnet_id = vestack_subnet.foo.id
	  security_group_ids = [vestack_security_group.foo.id]
	}
	
	resource "vestack_eip_address" "foo" {
    billing_type = "PostPaidByTraffic"
	}

	resource "vestack_eip_associate" "foo" {
		allocation_id = "${vestack_eip_address.foo.id}"
		instance_id = "${vestack_ecs_instance.foo.id}"
		instance_type = "EcsInstance"
	}
`

func TestAccVestackEipAssociateResource_Basic(t *testing.T) {
	resourceName := "vestack_eip_associate.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &eip_associate.VestackEipAssociateService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEipAssociateCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
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
