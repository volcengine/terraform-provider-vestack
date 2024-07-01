package ipv6_gateway_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/ipv6_gateway"
)

const testAccIpv6GatewayConfig = `
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

	resource "vestack_vpc_ipv6_gateway" "foo" {
	  vpc_id = "${vestack_vpc.foo.id}"
	  name = "acc-test-1"
	  description = "test"
	}

	data "vestack_vpc_ipv6_gateways" "foo" {
		ids = ["${vestack_vpc_ipv6_gateway.foo.id}"]
	}
`

func TestAccVestackIpv6GatewayDataSource_Basic(t *testing.T) {
	resourceName := "data.vestack_vpc_ipv6_gateways.foo"
	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ipv6_gateway.VestackIpv6GatewayService{},
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccIpv6GatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_gateways.#", "1"),
				),
			},
		},
	})
}
