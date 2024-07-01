package ipv6_address_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/ipv6_address"
)

const testAccV6AddressDatasourceConfig = `
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
	  enable_ipv6 = true
	}

	resource "vestack_subnet" "foo" {
	  subnet_name = "acc-test-subnet"
	  cidr_block = "172.16.0.0/24"
	  zone_id = "${data.vestack_zones.foo.zones[0].id}"
	  vpc_id = "${vestack_vpc.foo.id}"
	  ipv6_cidr_block = 1
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
	  ipv6_address_count = 1
	}

	data "vestack_vpc_ipv6_addresses" "foo"{
	  associated_instance_id = "${vestack_ecs_instance.foo.id}"
	}
`

func TestAccVestackV6AddressDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_vpc_ipv6_addresses.foo"
	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ipv6_address.VestackIpv6AddressService{},
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccV6AddressDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "1"),
				),
			},
		},
	})
}
