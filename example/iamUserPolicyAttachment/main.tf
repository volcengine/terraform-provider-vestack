resource "vestack_iam_user" "user" {
  user_name   = "TfTest"
  description = "test"
}

resource "vestack_iam_policy" "policy" {
  policy_name     = "TerraformResourceTest1"
  description     = "created by terraform 1"
  policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}

resource "vestack_iam_user_policy_attachment" "foo" {
  user_name   = vestack_iam_user.user.user_name
  policy_name = vestack_iam_policy.policy.policy_name
  policy_type = vestack_iam_policy.policy.policy_type
}