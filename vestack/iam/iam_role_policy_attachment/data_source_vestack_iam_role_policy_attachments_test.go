package iam_role_policy_attachment_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_role_policy_attachment"
)

func TestAccVestackIamRolePolicyAttachmentsDataSource_Basic(t *testing.T) {
	resourceName := "data.vestack_iam_role_policy_attachments.foo"
	roleName := "acc-test-role-ds"
	policyName := "acc-test-policy-ds"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_role_policy_attachment.VestackIamRolePolicyAttachmentService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamRolePolicyAttachmentsDataSourceConfig(roleName, policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "policies.#", "1"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "policies.0.attach_date"),
				),
			},
		},
	})
}

func testAccVestackIamRolePolicyAttachmentsDataSourceConfig(roleName, policyName string) string {
	return fmt.Sprintf(`
resource "vestack_iam_role" "foo" {
	role_name = "%s"
    display_name = "acc-test"
	description = "acc-test"
    trust_policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"auto_scaling\"]}}]}"
	max_session_duration = 3600
}

resource "vestack_iam_policy" "foo" {
    policy_name = "%s"
	description = "acc-test"
	policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}

resource "vestack_iam_role_policy_attachment" "foo" {
    policy_name = vestack_iam_policy.foo.policy_name
    policy_type = "Custom"
    role_name = vestack_iam_role.foo.role_name
}

data "vestack_iam_role_policy_attachments" "foo" {
  role_name = vestack_iam_role.foo.role_name
  depends_on = [vestack_iam_role_policy_attachment.foo]
}
`, roleName, policyName)
}
