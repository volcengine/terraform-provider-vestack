package iam_role_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_role"
)

const testAccVestackIamRolesDatasourceConfig = `
resource "vestack_iam_role" "foo1" {
	role_name = "acc-test-role1"
    display_name = "acc-test1"
	description = "acc-test1"
    trust_policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"auto_scaling\"]}}]}"
	max_session_duration = 3600
}

resource "vestack_iam_role" "foo2" {
    role_name = "acc-test-role2"
    display_name = "acc-test2"
	description = "acc-test2"
    trust_policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"ecs\"]}}]}"
	max_session_duration = 3600
}

data "vestack_iam_roles" "foo"{
    role_name = "${vestack_iam_role.foo1.role_name},${vestack_iam_role.foo2.role_name}"
}
`

func TestAccVestackIamRolesDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_iam_roles.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_role.VestackIamRoleService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamRolesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "roles.#", "2"),
				),
			},
		},
	})
}
