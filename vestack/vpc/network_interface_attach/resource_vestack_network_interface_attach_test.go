package network_interface_attach_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/network_interface_attach"
)

const testAccVestackNetworkInterfaceAttachCreateConfig = `
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

resource "vestack_network_interface" "foo" {
  network_interface_name = "acc-test-eni"
  description = "acc-test"
  subnet_id = "${vestack_subnet.foo.id}"
  security_group_ids = ["${vestack_security_group.foo.id}"]
  primary_ip_address = "172.16.0.253"
  port_security_enabled = false
  private_ip_address = ["172.16.0.2"]
  project_name = "default"
}

resource "vestack_network_interface_attach" "foo" {
    instance_id = "${vestack_ecs_instance.foo.id}"
    network_interface_id = "${vestack_network_interface.foo.id}"
}
`

func TestAccVestackNetworkInterfaceAttachResource_Basic(t *testing.T) {
	resourceName := "vestack_network_interface_attach.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_interface_attach.VestackNetworkInterfaceAttachService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackNetworkInterfaceAttachCreateConfig,
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
