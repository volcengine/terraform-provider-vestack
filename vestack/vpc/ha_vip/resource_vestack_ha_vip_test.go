package ha_vip_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/ha_vip"
)

const testAccVestackHaVipCreateConfig = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block = "172.16.0.0/24"
  zone_id = data.vestack_zones.foo.zones[0].id
  vpc_id = vestack_vpc.foo.id
}

resource "vestack_ha_vip" "foo" {
  ha_vip_name = "acc-test-ha-vip"
  description = "acc-test"
  subnet_id = vestack_subnet.foo.id
  ip_address = "172.16.0.5"
}

resource "vestack_eip_address" "foo" {
  billing_type = "PostPaidByTraffic"
}

resource "vestack_eip_associate" "foo" {
  allocation_id = vestack_eip_address.foo.id
  instance_id = vestack_ha_vip.foo.id
  instance_type = "HaVip"
}
`

func TestAccVestackHaVipResource_Basic(t *testing.T) {
	resourceName := "vestack_ha_vip.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return ha_vip.NewHaVipService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackHaVipCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "ha_vip_name", "acc-test-ha-vip"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ip_address", "172.16.0.5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "associated_instance_ids.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "associated_instance_type", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "associated_eip_address", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "associated_eip_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "master_instance_id", ""),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "vpc_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "created_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "updated_at"),
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

const testAccVestackHaVipUpdateConfig = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block = "172.16.0.0/24"
  zone_id = data.vestack_zones.foo.zones[0].id
  vpc_id = vestack_vpc.foo.id
}

resource "vestack_ha_vip" "foo" {
  ha_vip_name = "acc-test-ha-vip-new"
  description = "acc-test-new"
  subnet_id = vestack_subnet.foo.id
  ip_address = "172.16.0.5"
}

resource "vestack_eip_address" "foo" {
  billing_type = "PostPaidByTraffic"
}

resource "vestack_eip_associate" "foo" {
  allocation_id = vestack_eip_address.foo.id
  instance_id = vestack_ha_vip.foo.id
  instance_type = "HaVip"
}
`

func TestAccVestackHaVipResource_Update(t *testing.T) {
	resourceName := "vestack_ha_vip.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return ha_vip.NewHaVipService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackHaVipCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "ha_vip_name", "acc-test-ha-vip"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ip_address", "172.16.0.5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "associated_instance_ids.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "associated_instance_type", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "associated_eip_address", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "associated_eip_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "master_instance_id", ""),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "vpc_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "created_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "updated_at"),
				),
			},
			{
				Config: testAccVestackHaVipUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "ha_vip_name", "acc-test-ha-vip-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ip_address", "172.16.0.5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "associated_instance_ids.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "associated_instance_type", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "master_instance_id", ""),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "vpc_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "created_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "updated_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "associated_eip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "associated_eip_id"),
				),
			},
			{
				Config:             testAccVestackHaVipUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
