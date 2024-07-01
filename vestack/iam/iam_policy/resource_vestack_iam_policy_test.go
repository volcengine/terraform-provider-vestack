package iam_policy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_policy"
)

const testAccVestackIamPolicyCreateConfig = `
resource "vestack_iam_policy" "foo" {
    policy_name = "acc-test-policy"
	description = "acc-test"
	policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}
`

func TestAccVestackIamPolicyResource_Basic(t *testing.T) {
	resourceName := "vestack_iam_policy.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_policy.VestackIamPolicyService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamPolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_name", "acc-test-policy"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_document", "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "policy_trn"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "policy_type"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_date"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "update_date"),
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

const testAccVestackIamPolicyUpdateConfig = `
resource "vestack_iam_policy" "foo" {
    policy_name = "acc-test-policy-new"
	description = "acc-test-new"
	policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingConfigurations\"],\"Resource\":[\"*\"]}]}"
}
`

func TestAccVestackIamPolicyResource_Update(t *testing.T) {
	resourceName := "vestack_iam_policy.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_policy.VestackIamPolicyService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamPolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_name", "acc-test-policy"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_document", "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "policy_trn"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "policy_type"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_date"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "update_date"),
				),
			},
			{
				Config: testAccVestackIamPolicyUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_name", "acc-test-policy-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_document", "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingConfigurations\"],\"Resource\":[\"*\"]}]}"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "policy_trn"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "policy_type"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_date"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "update_date"),
				),
			},
			{
				Config:             testAccVestackIamPolicyUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
