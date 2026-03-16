package snat_entry_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/nat/snat_entry"
	"testing"
)

const testAccVestackSnatEntryCreateConfig = `
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

resource "vestack_snat_entry" "foo" {
	snat_entry_name = "acc-test-snat-entry"
    nat_gateway_id = "${vestack_nat_gateway.foo.id}"
	eip_id = "${vestack_eip_address.foo.id}"
	subnet_id = "${vestack_subnet.foo.id}"
	depends_on = [vestack_eip_associate.foo]
}
`

func TestAccVestackSnatEntryResource_Basic(t *testing.T) {
	resourceName := "vestack_snat_entry.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &snat_entry.VestackSnatEntryService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackSnatEntryCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "snat_entry_name", "acc-test-snat-entry"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "nat_gateway_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "source_cidr"),
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

const testAccVestackSnatEntryCreateSourceCidrConfig = `
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

resource "vestack_snat_entry" "foo" {
	snat_entry_name = "acc-test-snat-entry"
    nat_gateway_id = "${vestack_nat_gateway.foo.id}"
	eip_id = "${vestack_eip_address.foo.id}"
	source_cidr = "172.16.0.0/24"
	depends_on = [vestack_eip_associate.foo]
}
`

func TestAccVestackSnatEntryResource_SourceCidr(t *testing.T) {
	resourceName := "vestack_snat_entry.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &snat_entry.VestackSnatEntryService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackSnatEntryCreateSourceCidrConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "snat_entry_name", "acc-test-snat-entry"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "source_cidr", "172.16.0.0/24"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "nat_gateway_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
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

const testAccVestackSnatEntryUpdateConfig = `
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
	name = "acc-test-eip1"
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

resource "vestack_snat_entry" "foo" {
	snat_entry_name = "acc-test-snat-entry-new"
    nat_gateway_id = "${vestack_nat_gateway.foo.id}"
	eip_id = "${vestack_eip_address.foo1.id}"
	subnet_id = "${vestack_subnet.foo.id}"
	depends_on = [vestack_eip_associate.foo1]
}
`

func TestAccVestackSnatEntryResource_Update(t *testing.T) {
	resourceName := "vestack_snat_entry.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &snat_entry.VestackSnatEntryService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackSnatEntryCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "snat_entry_name", "acc-test-snat-entry"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "nat_gateway_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "source_cidr"),
				),
			},
			{
				Config: testAccVestackSnatEntryUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "snat_entry_name", "acc-test-snat-entry-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "nat_gateway_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "source_cidr"),
				),
			},
			{
				Config:             testAccVestackSnatEntryUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
