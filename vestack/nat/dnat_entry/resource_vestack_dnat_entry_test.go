package dnat_entry_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/nat/dnat_entry"
	"testing"
)

const testAccVestackDnatEntryCreateConfig = `
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

resource "vestack_nat_gateway" "foo" {
	vpc_id = "${vestack_vpc.foo.id}"
    subnet_id = "${vestack_subnet.foo.id}"
	spec = "Small"
	nat_gateway_name = "acc-test-ng"
	description = "acc-test"
	billing_type = "PostPaid"
	project_name = "default"
	tags {
		key = "k1"
		value = "v1"
	}
}

resource "vestack_eip_address" "foo" {
	name = "acc-test-eip"
    description = "acc-test"
    bandwidth = 1
    billing_type = "PostPaidByBandwidth"
    isp = "BGP"
}

resource "vestack_eip_associate" "foo" {
	allocation_id = "${vestack_eip_address.foo.id}"
	instance_id = "${vestack_nat_gateway.foo.id}"
	instance_type = "Nat"
}

resource "vestack_dnat_entry" "foo" {
	dnat_entry_name = "acc-test-dnat-entry"
    external_ip = "${vestack_eip_address.foo.eip_address}"
    external_port = 80
    internal_ip = "172.16.0.10"
    internal_port = 80
    nat_gateway_id = "${vestack_nat_gateway.foo.id}"
    protocol = "tcp"
	depends_on = [vestack_eip_associate.foo]
}
`

func TestAccVestackDnatEntryResource_Basic(t *testing.T) {
	resourceName := "vestack_dnat_entry.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &dnat_entry.VestackDnatEntryService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackDnatEntryCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "dnat_entry_name", "acc-test-dnat-entry"),
					resource.TestCheckResourceAttr(acc.ResourceId, "external_port", "80"),
					resource.TestCheckResourceAttr(acc.ResourceId, "internal_ip", "172.16.0.10"),
					resource.TestCheckResourceAttr(acc.ResourceId, "internal_port", "80"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "tcp"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "external_ip"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "nat_gateway_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "dnat_entry_id"),
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

const testAccVestackDnatEntryUpdateConfig = `
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

resource "vestack_nat_gateway" "foo" {
	vpc_id = "${vestack_vpc.foo.id}"
    subnet_id = "${vestack_subnet.foo.id}"
	spec = "Small"
	nat_gateway_name = "acc-test-ng"
	description = "acc-test"
	billing_type = "PostPaid"
	project_name = "default"
	tags {
		key = "k1"
		value = "v1"
	}
}

resource "vestack_eip_address" "foo" {
	name = "acc-test-eip"
    description = "acc-test"
    bandwidth = 1
    billing_type = "PostPaidByBandwidth"
    isp = "BGP"
}

resource "vestack_eip_associate" "foo" {
	allocation_id = "${vestack_eip_address.foo.id}"
	instance_id = "${vestack_nat_gateway.foo.id}"
	instance_type = "Nat"
}

resource "vestack_eip_address" "foo1" {
	name = "acc-test-eip"
    description = "acc-test"
    bandwidth = 1
    billing_type = "PostPaidByBandwidth"
    isp = "BGP"
}

resource "vestack_eip_associate" "foo1" {
	allocation_id = "${vestack_eip_address.foo1.id}"
	instance_id = "${vestack_nat_gateway.foo.id}"
	instance_type = "Nat"
}

resource "vestack_dnat_entry" "foo" {
	dnat_entry_name = "acc-test-dnat-entry-new"
    external_ip = "${vestack_eip_address.foo1.eip_address}"
    external_port = 90
    internal_ip = "172.16.0.17"
    internal_port = 90
    nat_gateway_id = "${vestack_nat_gateway.foo.id}"
    protocol = "udp"
	depends_on = [vestack_eip_associate.foo1]
}
`

func TestAccVestackDnatEntryResource_Update(t *testing.T) {
	resourceName := "vestack_dnat_entry.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &dnat_entry.VestackDnatEntryService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackDnatEntryCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "dnat_entry_name", "acc-test-dnat-entry"),
					resource.TestCheckResourceAttr(acc.ResourceId, "external_port", "80"),
					resource.TestCheckResourceAttr(acc.ResourceId, "internal_ip", "172.16.0.10"),
					resource.TestCheckResourceAttr(acc.ResourceId, "internal_port", "80"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "tcp"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "external_ip"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "nat_gateway_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "dnat_entry_id"),
				),
			},
			{
				Config: testAccVestackDnatEntryUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "dnat_entry_name", "acc-test-dnat-entry-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "external_port", "90"),
					resource.TestCheckResourceAttr(acc.ResourceId, "internal_ip", "172.16.0.17"),
					resource.TestCheckResourceAttr(acc.ResourceId, "internal_port", "90"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "udp"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "external_ip"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "nat_gateway_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "dnat_entry_id"),
				),
			},
			{
				Config:             testAccVestackDnatEntryUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
