package vpc_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/vpc"
	"testing"
)

const testAccVpcForCreate = `
resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}
`

const testAccVpcForUpdate = `
resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
  dns_servers = ["8.8.8.8", "114.114.114.114"]

  tags {
    key = "k2"
    value = "v2"
  }

  tags {
    key = "k1"
    value = "v1"
  }
}
`

func TestAccVestackVpcResource_Basic(t *testing.T) {
	resourceName := "vestack_vpc.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &vpc.VestackVpcService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcForCreate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "vpc_name", "acc-test-vpc"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cidr_block", "172.16.0.0/16"),
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

func TestAccVestackVpcResource_Update(t *testing.T) {
	resourceName := "vestack_vpc.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &vpc.VestackVpcService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcForCreate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "vpc_name", "acc-test-vpc"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cidr_block", "172.16.0.0/16"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config: testAccVpcForUpdate,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "vpc_name", "acc-test-vpc"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cidr_block", "172.16.0.0/16"),
					// update attr check
					resource.TestCheckResourceAttr(acc.ResourceId, "dns_servers.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "2"),
					vestack.TestCheckTypeSetElemAttr(acc.ResourceId, "dns_servers.*", "8.8.8.8"),
					vestack.TestCheckTypeSetElemAttr(acc.ResourceId, "dns_servers.*", "114.114.114.114"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k2",
						"value": "v2",
					}),
					// compute status check
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config:             testAccVpcForUpdate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
