package prefix_list_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/prefix_list"
)

const testAccVestackVpcPrefixListCreateConfig = `
resource "vestack_vpc_prefix_list" "foo" {
	prefix_list_name = "acc-test-prefix"
    max_entries = 3
	description = "acc test description"
    ip_version = "IPv4"
	prefix_list_entries {
		cidr = "192.168.4.0/28"
		description = "acc-test-1"
	}
	prefix_list_entries {
		cidr = "192.168.5.0/28"
		description = "acc-test-2"
	}
	tags {
		key = "tf-key1"
		value = "tf-value1"
	}
}
`

const testAccVestackVpcPrefixListUpdateConfig = `
resource "vestack_vpc_prefix_list" "foo" {
    prefix_list_name = "acc-test-prefix-modify"
    max_entries = 4
	description = "acc test description"
    ip_version = "IPv4"
	prefix_list_entries {
		cidr = "192.168.4.0/28"
		description = "acc-test-1"
	}
	prefix_list_entries {
		cidr = "192.168.7.0/28"
		description = "acc-test-3"
	}
	tags {
		key = "tf-key1"
		value = "tf-value1"
	}
}
`

func TestAccVestackVpcPrefixListResource_Basic(t *testing.T) {
	resourceName := "vestack_vpc_prefix_list.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return prefix_list.NewVpcPrefixListService(client)
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
				Config: testAccVestackVpcPrefixListCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc test description"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ip_version", "IPv4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "max_entries", "3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "prefix_list_entries.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "prefix_list_name", "acc-test-prefix"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
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

func TestAccVestackVpcPrefixListResource_Update(t *testing.T) {
	resourceName := "vestack_vpc_prefix_list.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return prefix_list.NewVpcPrefixListService(client)
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
				Config: testAccVestackVpcPrefixListCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc test description"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ip_version", "IPv4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "max_entries", "3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "prefix_list_entries.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "prefix_list_name", "acc-test-prefix"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
				),
			},
			{
				Config: testAccVestackVpcPrefixListUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc test description"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ip_version", "IPv4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "max_entries", "4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "prefix_list_entries.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "prefix_list_name", "acc-test-prefix-modify"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-key1",
						"value": "tf-value1",
					}),
				),
			},
			{
				Config:             testAccVestackVpcPrefixListUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
