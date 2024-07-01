package network_interface_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/network_interface"
)

const testAccVestackNetworkInterfaceCreateConfig = `
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
  security_group_name = "acc-test-sg"
  vpc_id = "${vestack_vpc.foo.id}"
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
  tags {
    key = "k1"
    value = "v1"
  }
}
`

func TestAccVestackNetworkInterfaceResource_Basic(t *testing.T) {
	resourceName := "vestack_network_interface.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_interface.VestackNetworkInterfaceService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackNetworkInterfaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_interface_name", "acc-test-eni"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_security_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.253"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "private_ip_address.#", "1"),
					vestack.TestCheckTypeSetElemAttr(acc.ResourceId, "private_ip_address.*", "172.16.0.2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
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

const testAccVestackNetworkInterfaceUpdateConfig1 = `
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
  security_group_name = "acc-test-sg"
  vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_network_interface" "foo" {
  network_interface_name = "acc-test-eni-new"
  description = "acc-test-new"
  subnet_id = "${vestack_subnet.foo.id}"
  security_group_ids = ["${vestack_security_group.foo.id}"]
  primary_ip_address = "172.16.0.253"
  port_security_enabled = false
  private_ip_address = ["172.16.0.2"]
  project_name = "default"
  tags {
    key = "k1"
    value = "v1"
  }
  tags {
    key = "k2"
    value = "v2"
  }
}
`

func TestAccVestackNetworkInterfaceResource_Update1(t *testing.T) {
	resourceName := "vestack_network_interface.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_interface.VestackNetworkInterfaceService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackNetworkInterfaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_interface_name", "acc-test-eni"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_security_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.253"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "private_ip_address.#", "1"),
					vestack.TestCheckTypeSetElemAttr(acc.ResourceId, "private_ip_address.*", "172.16.0.2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
				),
			},
			{
				Config: testAccVestackNetworkInterfaceUpdateConfig1,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_interface_name", "acc-test-eni-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_security_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.253"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "private_ip_address.#", "1"),
					vestack.TestCheckTypeSetElemAttr(acc.ResourceId, "private_ip_address.*", "172.16.0.2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "2"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k2",
						"value": "v2",
					}),
				),
			},
			{
				Config:             testAccVestackNetworkInterfaceUpdateConfig1,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccVestackNetworkInterfaceUpdateConfig2 = `
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
  security_group_name = "acc-test-sg"
  vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_network_interface" "foo" {
  network_interface_name = "acc-test-eni"
  description = "acc-test"
  subnet_id = "${vestack_subnet.foo.id}"
  security_group_ids = ["${vestack_security_group.foo.id}"]
  primary_ip_address = "172.16.0.253"
  port_security_enabled = false
  private_ip_address = ["172.16.0.3", "172.16.0.4"]
  project_name = "default"
  tags {
    key = "k1"
    value = "v1"
  }
}
`

func TestAccVestackNetworkInterfaceResource_Update2(t *testing.T) {
	resourceName := "vestack_network_interface.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_interface.VestackNetworkInterfaceService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackNetworkInterfaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_interface_name", "acc-test-eni"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_security_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.253"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "private_ip_address.#", "1"),
					vestack.TestCheckTypeSetElemAttr(acc.ResourceId, "private_ip_address.*", "172.16.0.2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
				),
			},
			{
				Config: testAccVestackNetworkInterfaceUpdateConfig2,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_interface_name", "acc-test-eni"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_security_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.253"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "private_ip_address.#", "2"),
					vestack.TestCheckTypeSetElemAttr(acc.ResourceId, "private_ip_address.*", "172.16.0.3"),
					vestack.TestCheckTypeSetElemAttr(acc.ResourceId, "private_ip_address.*", "172.16.0.4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
				),
			},
			{
				Config:             testAccVestackNetworkInterfaceUpdateConfig2,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
