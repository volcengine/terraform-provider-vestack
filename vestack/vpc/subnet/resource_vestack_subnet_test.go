package subnet_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/subnet"
	"testing"
)

const testAccSubnetForCreate = `
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
`

const testAccSubnetForUpdate = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet_new"
  description = "tfdesc"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.vestack_zones.foo.zones[0].id}"
  vpc_id = "${vestack_vpc.foo.id}"
}
`

func TestAccVestackSubnetResource_Basic(t *testing.T) {
	resourceName := "vestack_subnet.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &subnet.VestackSubnetService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetForCreate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "subnet_name", "acc-test-subnet"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cidr_block", "172.16.0.0/24"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
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

func TestAccVestackSubnetResource_Update(t *testing.T) {
	resourceName := "vestack_subnet.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &subnet.VestackSubnetService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetForCreate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "subnet_name", "acc-test-subnet"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cidr_block", "172.16.0.0/24"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config: testAccSubnetForUpdate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "tfdesc"),
					resource.TestCheckResourceAttr(acc.ResourceId, "subnet_name", "acc-test-subnet_new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cidr_block", "172.16.0.0/24"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config:             testAccSubnetForUpdate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
