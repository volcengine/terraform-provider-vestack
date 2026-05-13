package server_group_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/server_group"
	"testing"
)

const testAccVestackServerGroupCreateConfig = `
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

resource "vestack_clb" "foo" {
  type = "public"
  subnet_id = "${vestack_subnet.foo.id}"
  load_balancer_spec = "small_1"
  description = "acc0Demo"
  load_balancer_name = "acc-test-create"
  eip_billing_config {
    isp = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth = 1
  }
}

resource "vestack_server_group" "foo" {
  load_balancer_id = "${vestack_clb.foo.id}"
  server_group_name = "acc-test-create"
  description = "hello demo11"
}

`

const testAccVestackServerGroupUpdateConfig = `
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

resource "vestack_clb" "foo" {
  type = "public"
  subnet_id = "${vestack_subnet.foo.id}"
  load_balancer_spec = "small_1"
  description = "acc0Demo"
  load_balancer_name = "acc-test-create"
  eip_billing_config {
    isp = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth = 1
  }
}

resource "vestack_server_group" "foo" {
  load_balancer_id = "${vestack_clb.foo.id}"
  server_group_name = "acc-test-create1"
  description = "acc hello demo11"
}

`

func TestAccVestackServerGroupResource_Basic(t *testing.T) {
	resourceName := "vestack_server_group.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &server_group.VestackServerGroupService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackServerGroupCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "hello demo11"),
					resource.TestCheckResourceAttr(acc.ResourceId, "server_group_name", "acc-test-create"),
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

func TestAccVestackServerGroupResource_Update(t *testing.T) {
	resourceName := "vestack_server_group.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &server_group.VestackServerGroupService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackServerGroupCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "hello demo11"),
					resource.TestCheckResourceAttr(acc.ResourceId, "server_group_name", "acc-test-create"),
				),
			},
			{
				Config: testAccVestackServerGroupUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc hello demo11"),
					resource.TestCheckResourceAttr(acc.ResourceId, "server_group_name", "acc-test-create1"),
				),
			},
			{
				Config:             testAccVestackServerGroupUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccVestackServerGroupCreateConfigIpv6 = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "vpc_ipv6" {
  vpc_name = "acc-test-vpc-ipv6"
  cidr_block = "172.16.0.0/16"
  enable_ipv6 = true
}

resource "vestack_subnet" "subnet_ipv6" {
  subnet_name = "acc-test-subnet-ipv6"
  cidr_block = "172.16.0.0/24"
  zone_id = data.vestack_zones.foo.zones[1].id
  vpc_id = vestack_vpc.vpc_ipv6.id
  ipv6_cidr_block = 1
}

resource "vestack_clb" "private_clb_ipv6" {
  type = "private"
  subnet_id = vestack_subnet.subnet_ipv6.id
  load_balancer_name = "acc-test-clb-ipv6"
  load_balancer_spec = "small_1"
  description = "acc-test-demo"
  project_name = "default"
  address_ip_version = "DualStack"
  tags {
    key = "k1"
    value = "v1"
  }
}

resource "vestack_server_group" "foo" {
  load_balancer_id = "${vestack_clb.private_clb_ipv6.id}"
  server_group_name = "acc-test-sg-ipv6"
  description = "acc-test"
  address_ip_version = "ipv6"
}
`

func TestAccVestackServerGroupResource_CreateIpv6(t *testing.T) {
	resourceName := "vestack_server_group.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &server_group.VestackServerGroupService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackServerGroupCreateConfigIpv6,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "server_group_name", "acc-test-sg-ipv6"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "address_ip_version", "ipv6"),
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
