resource "vestack_iam_policy" "foo" {
  policy_name     = "acc-test-policy"
  description     = "acc-test"
  policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}

data "vestack_iam_policies" "foo" {
  query = vestack_iam_policy.foo.description
  #  user_name = "user-test"
  #  role_name = "test-role"
}
