package prefix_list_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/prefix_list"
)

const testAccVestackVpcPrefixListsDatasourceConfig = `
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

data "vestack_vpc_prefix_lists" "foo" {
  ids = [vestack_vpc_prefix_list.foo.id]
}
`

func TestAccVestackVpcPrefixListsDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_vpc_prefix_lists.foo"

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
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackVpcPrefixListsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "prefix_lists.#", "1"),
				),
			},
		},
	})
}
