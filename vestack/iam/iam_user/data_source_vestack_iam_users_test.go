package iam_user_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_user"
	"testing"
)

const testAccVestackIamUsersDatasourceConfig = `
resource "vestack_iam_user" "foo" {
  user_name = "acc-test-user"
  description = "acc test"
  display_name = "name"
}
data "vestack_iam_users" "foo"{
    user_names = [vestack_iam_user.foo.user_name]
}
`

func TestAccVestackIamUsersDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_iam_users.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_user.VestackIamUserService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamUsersDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "users.#", "1"),
				),
			},
		},
	})
}
