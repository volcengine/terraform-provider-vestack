package network_acl_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/network_acl"
)

const testAccNetworkAclCreate = `
	data "vestack_zones" "foo"{
	}

	resource "vestack_vpc" "foo" {
	  vpc_name   = "acc-test-vpc"
	  cidr_block = "172.16.0.0/16"
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
`

const testAccNetworkAclUpdate = `
	data "vestack_zones" "foo"{
	}

	resource "vestack_vpc" "foo" {
	  vpc_name   = "acc-test-vpc"
	  cidr_block = "172.16.0.0/16"
	}

	resource "vestack_network_acl" "foo" {
	  vpc_id = "${vestack_vpc.foo.id}"
	  network_acl_name = "acc-test-acl2"
	  ingress_acl_entries {
		network_acl_entry_name = "acc-ingress1"
		policy = "accept"
		protocol = "all"
		source_cidr_ip = "192.168.1.0/24"
	  }
	  ingress_acl_entries {
		network_acl_entry_name = "acc-ingress2"
		policy = "accept"
		protocol = "all"
		source_cidr_ip = "192.168.0.0/24"
	  }
	  egress_acl_entries {
		network_acl_entry_name = "acc-egress3"
		policy = "accept"
		protocol = "all"
		destination_cidr_ip = "192.168.0.0/16"
	  }
	  egress_acl_entries {
		network_acl_entry_name = "acc-egress4"
		policy = "accept"
		protocol = "all"
		destination_cidr_ip = "192.168.0.0/20"
	  }
	}
`

func TestAccVpcNetworkAclResource_Basic(t *testing.T) {
	resourceName := "vestack_network_acl.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_acl.VestackNetworkAclService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAclCreate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_acl_name", "acc-test-acl"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ingress_acl_entries.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "ingress_acl_entries.*", map[string]string{
						"network_acl_entry_name": "acc-ingress1",
						"policy":                 "accept",
						"protocol":               "all",
						"source_cidr_ip":         "192.168.0.0/24",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "egress_acl_entries.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "egress_acl_entries.*", map[string]string{
						"network_acl_entry_name": "acc-egress2",
						"policy":                 "accept",
						"protocol":               "all",
						"destination_cidr_ip":    "192.168.0.0/16",
					}),
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

func TestAccVpcNetworkAclResource_Update(t *testing.T) {
	resourceName := "vestack_network_acl.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_acl.VestackNetworkAclService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAclCreate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_acl_name", "acc-test-acl"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ingress_acl_entries.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "ingress_acl_entries.*", map[string]string{
						"network_acl_entry_name": "acc-ingress1",
						"policy":                 "accept",
						"protocol":               "all",
						"source_cidr_ip":         "192.168.0.0/24",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "egress_acl_entries.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "egress_acl_entries.*", map[string]string{
						"network_acl_entry_name": "acc-egress2",
						"policy":                 "accept",
						"protocol":               "all",
						"destination_cidr_ip":    "192.168.0.0/16",
					}),
				),
			},
			{
				Config: testAccNetworkAclUpdate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_acl_name", "acc-test-acl2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ingress_acl_entries.#", "2"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "ingress_acl_entries.*", map[string]string{
						"network_acl_entry_name": "acc-ingress1",
						"policy":                 "accept",
						"protocol":               "all",
						"source_cidr_ip":         "192.168.1.0/24",
					}),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "ingress_acl_entries.*", map[string]string{
						"network_acl_entry_name": "acc-ingress2",
						"policy":                 "accept",
						"protocol":               "all",
						"source_cidr_ip":         "192.168.0.0/24",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "egress_acl_entries.#", "2"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "egress_acl_entries.*", map[string]string{
						"network_acl_entry_name": "acc-egress3",
						"policy":                 "accept",
						"protocol":               "all",
						"destination_cidr_ip":    "192.168.0.0/16",
					}),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "egress_acl_entries.*", map[string]string{
						"network_acl_entry_name": "acc-egress4",
						"policy":                 "accept",
						"protocol":               "all",
						"destination_cidr_ip":    "192.168.0.0/20",
					}),
				),
			},
			{
				Config:             testAccNetworkAclUpdate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
