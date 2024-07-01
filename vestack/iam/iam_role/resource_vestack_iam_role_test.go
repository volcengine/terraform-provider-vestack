package iam_role_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_role"
)

const testAccVestackIamRoleCreateConfig = `
resource "vestack_iam_role" "foo" {
	role_name = "acc-test-role"
    display_name = "acc-test"
	description = "acc-test"
    trust_policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"auto_scaling\"]}}]}"
	max_session_duration = 3600
}
`

func TestAccVestackIamRoleResource_Basic(t *testing.T) {
	resourceName := "vestack_iam_role.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_role.VestackIamRoleService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamRoleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "display_name", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "max_session_duration", "3600"),
					resource.TestCheckResourceAttr(acc.ResourceId, "role_name", "acc-test-role"),
					resource.TestCheckResourceAttr(acc.ResourceId, "trust_policy_document", "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"auto_scaling\"]}}]}"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trn"),
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

const testAccVestackIamRoleUpdateConfig = `
resource "vestack_iam_role" "foo" {
    role_name = "acc-test-role-new"
    display_name = "acc-test-new"
	description = "acc-test-new"
    trust_policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"ecs\"]}}]}"
	max_session_duration = 3700
}
`

func TestAccVestackIamRoleResource_Update(t *testing.T) {
	resourceName := "vestack_iam_role.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_role.VestackIamRoleService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamRoleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "display_name", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "max_session_duration", "3600"),
					resource.TestCheckResourceAttr(acc.ResourceId, "role_name", "acc-test-role"),
					resource.TestCheckResourceAttr(acc.ResourceId, "trust_policy_document", "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"auto_scaling\"]}}]}"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trn"),
				),
			},
			{
				Config: testAccVestackIamRoleUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "display_name", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "max_session_duration", "3700"),
					resource.TestCheckResourceAttr(acc.ResourceId, "role_name", "acc-test-role-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "trust_policy_document", "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"ecs\"]}}]}"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trn"),
				),
			},
			{
				Config:             testAccVestackIamRoleUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
