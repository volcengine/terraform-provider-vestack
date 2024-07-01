package iam_user_policy_attachment_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_user_policy_attachment"
	"testing"
)

const testAccVestackIamUserPolicyAttachmentCreateConfig = `
resource "vestack_iam_user" "foo" {
  user_name = "acc-test-user"
  description = "acc test"
  display_name = "name"
}
resource "vestack_iam_policy" "foo" {
    policy_name = "acc-test-policy"
	description = "acc-test"
	policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}
resource "vestack_iam_user_policy_attachment" "foo" {
    policy_name = vestack_iam_policy.foo.policy_name
    policy_type = "Custom"
    user_name = vestack_iam_user.foo.user_name
}
`

func TestAccVestackIamUserPolicyAttachmentResource_Basic(t *testing.T) {
	resourceName := "vestack_iam_user_policy_attachment.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_user_policy_attachment.VestackIamUserPolicyAttachmentService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamUserPolicyAttachmentCreateConfig,
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
