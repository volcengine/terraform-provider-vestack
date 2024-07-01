package iam_role_policy_attachment_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_role_policy_attachment"
)

const testAccVestackIamRolePolicyAttachmentCreateConfig = `
resource "vestack_iam_role" "foo" {
	role_name = "acc-test-role"
    display_name = "acc-test"
	description = "acc-test"
    trust_policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"auto_scaling\"]}}]}"
	max_session_duration = 3600
}

resource "vestack_iam_policy" "foo" {
    policy_name = "acc-test-policy"
	description = "acc-test"
	policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}

resource "vestack_iam_role_policy_attachment" "foo" {
    policy_name = vestack_iam_policy.foo.policy_name
    policy_type = "Custom"
    role_name = vestack_iam_role.foo.role_name
}
`

func TestAccVestackIamRolePolicyAttachmentResource_Basic(t *testing.T) {
	resourceName := "vestack_iam_role_policy_attachment.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_role_policy_attachment.VestackIamRolePolicyAttachmentService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamRolePolicyAttachmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_type", "Custom"),
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
