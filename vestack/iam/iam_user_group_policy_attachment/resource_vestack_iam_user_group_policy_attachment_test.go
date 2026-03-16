package iam_user_group_policy_attachment_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_user_group_policy_attachment"
)

const testAccVestackIamUserGroupPolicyAttachmentCreateConfig = `
resource "vestack_iam_policy" "foo" {
    policy_name = "acc-test-policy"
	description = "acc-test"
	policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}

resource "vestack_iam_user_group" "foo" {
  user_group_name = "acc-test-group"
  description = "acc-test"
  display_name = "acc-test"
}

resource "vestack_iam_user_group_policy_attachment" "foo" {
    policy_name = vestack_iam_policy.foo.policy_name
    policy_type = "Custom"
    user_group_name = vestack_iam_user_group.foo.user_group_name
}
`

func TestAccVestackIamUserGroupPolicyAttachmentResource_Basic(t *testing.T) {
	resourceName := "vestack_iam_user_group_policy_attachment.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return iam_user_group_policy_attachment.NewIamUserGroupPolicyAttachmentService(client)
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
				Config: testAccVestackIamUserGroupPolicyAttachmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_name", "acc-test-policy"),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_type", "Custom"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_group_name", "acc-test-group"),
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
