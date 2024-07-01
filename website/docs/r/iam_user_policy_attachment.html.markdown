---
subcategory: "IAM"
layout: "vestack"
page_title: "Vestack: vestack_iam_user_policy_attachment"
sidebar_current: "docs-vestack-resource-iam_user_policy_attachment"
description: |-
  Provides a resource to manage iam user policy attachment
---
# vestack_iam_user_policy_attachment
Provides a resource to manage iam user policy attachment
## Example Usage
```hcl
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
```
## Argument Reference
The following arguments are supported:
* `policy_name` - (Required, ForceNew) The name of the Policy.
* `policy_type` - (Required, ForceNew) The type of the Policy.
* `user_name` - (Required, ForceNew) The name of the user.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
Iam user policy attachment can be imported using the UserName:PolicyName:PolicyType, e.g.
```
$ terraform import vestack_iam_user_policy_attachment.default TerraformTestUser:TerraformTestPolicy:Custom
```

