package iam_user_group_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_user_group"
)

const testAccVestackIamUserGroupsDatasourceConfig = `
resource "vestack_iam_user_group" "foo" {
  user_group_name = "acc-test-group"
  description = "acc-test"
  display_name = "acc-test"
}

data "vestack_iam_user_groups" "foo"{
    query = "acc-test-group"
}
`

func TestAccVestackIamUserGroupsDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_iam_user_groups.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return iam_user_group.NewIamUserGroupService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamUserGroupsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(acc.ResourceId, "user_groups"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "user_groups.0.update_date"),
				),
			},
		},
	})
}
