package iam_policy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_policy"
)

const testAccVestackIamPoliciesDatasourceConfig = `
resource "vestack_iam_policy" "foo1" {
    policy_name = "acc-test-policy1"
	description = "acc-test"
	policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}

resource "vestack_iam_policy" "foo2" {
    policy_name = "acc-test-policy2"
	description = "acc-test"
	policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingConfigurations\"],\"Resource\":[\"*\"]}]}"
}

data "vestack_iam_policies" "foo"{
    query = "${vestack_iam_policy.foo1.description}"
}
`

func TestAccVestackIamPoliciesDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_iam_policies.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_policy.VestackIamPolicyService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamPoliciesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "policies.#", "2"),
				),
			},
		},
	})
}
