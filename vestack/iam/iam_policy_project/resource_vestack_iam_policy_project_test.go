package iam_policy_project_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_policy_project"
)

func TestAccVestackIamPolicyProjectResource_Basic(t *testing.T) {
	resourceName := "vestack_iam_policy_project.foo"
	userName := "acc-test-user-project"
	policyName := "acc-test-policy-project"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(&vestack.AccTestResource{
			ResourceId: resourceName,
			Svc:        &iam_policy_project.VestackIamPolicyProjectService{},
		}),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamPolicyProjectResourceConfig(userName, policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "project_names.0", "default"),
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

func TestAccVestackIamPolicyProjectResource_SystemPolicy(t *testing.T) {
	resourceName := "vestack_iam_policy_project.foo"
	userName := "acc-test-user-project-sys"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(&vestack.AccTestResource{
			ResourceId: resourceName,
			Svc:        &iam_policy_project.VestackIamPolicyProjectService{},
		}),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamPolicyProjectResourceConfigSystem(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "project_names.0", "default"),
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

func testAccVestackIamPolicyProjectResourceConfigSystem(userName string) string {
	return fmt.Sprintf(`
resource "vestack_iam_user" "foo" {
  user_name = "%s"
  display_name = "acc-test"
}

resource "vestack_iam_policy_project" "foo" {
  principal_type = "User"
  principal_name = vestack_iam_user.foo.user_name
  policy_type = "System"
  policy_name = "AdministratorAccess"
  project_names = ["default"]
}
`, userName)
}

func testAccVestackIamPolicyProjectResourceConfig(userName, policyName string) string {
	return fmt.Sprintf(`
resource "vestack_iam_user" "foo" {
  user_name = "%s"
  display_name = "acc-test"
}

resource "vestack_iam_policy" "foo" {
  policy_name = "%s"
  description = "acc-test"
  policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}

resource "vestack_iam_policy_project" "foo" {
  principal_type = "User"
  principal_name = vestack_iam_user.foo.user_name
  policy_type = "Custom"
  policy_name = vestack_iam_policy.foo.policy_name
  project_names = ["default"]
}
`, userName, policyName)
}
