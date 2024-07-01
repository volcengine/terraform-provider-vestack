package network_acl_associate_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/network_acl_associate"
)

const testAccNetworkAclAssociateConfig = `
	data "vestack_zones" "foo"{
	}

	resource "vestack_vpc" "foo" {
	  vpc_name   = "acc-test-vpc"
	  cidr_block = "172.16.0.0/16"
	}

	resource "vestack_subnet" "foo" {
	  subnet_name = "acc-test-subnet"
	  cidr_block = "172.16.0.0/16"
	  zone_id = "${data.vestack_zones.foo.zones[0].id}"
	  vpc_id = "${vestack_vpc.foo.id}"
	}		

	resource "vestack_network_acl" "foo" {
	  vpc_id = "${vestack_vpc.foo.id}"
	  network_acl_name = "acc-test-acl"
	  ingress_acl_entries {
		network_acl_entry_name = "acc-ingress1"
		policy = "accept"
		protocol = "all"
		source_cidr_ip = "192.168.0.0/24"
	  }
	  egress_acl_entries {
		network_acl_entry_name = "acc-egress2"
		policy = "accept"
		protocol = "all"
		destination_cidr_ip = "192.168.0.0/16"
	  }
	}

	resource "vestack_network_acl_associate" "foo" {
	  network_acl_id = "${vestack_network_acl.foo.id}"
	  resource_id = "${vestack_subnet.foo.id}"
	}

`

func TestAccVestackNetworkAclAssociateResource_Basic(t *testing.T) {
	resourceName := "vestack_network_acl_associate.foo"
	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_acl_associate.VestackNetworkAclAssociateService{},
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAclAssociateConfig,
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
