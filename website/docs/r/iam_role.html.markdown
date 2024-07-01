---
subcategory: "IAM"
layout: "vestack"
page_title: "Vestack: vestack_iam_role"
sidebar_current: "docs-vestack-resource-iam_role"
description: |-
  Provides a resource to manage iam role
---
# vestack_iam_role
Provides a resource to manage iam role
## Example Usage
```hcl
resource "vestack_iam_role" "foo" {
  role_name             = "acc-test-role"
  display_name          = "acc-test"
  description           = "acc-test"
  trust_policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"auto_scaling\"]}}]}"
  max_session_duration  = 3600
}
```
## Argument Reference
The following arguments are supported:
* `display_name` - (Required) The display name of the Role.
* `role_name` - (Required, ForceNew) The name of the Role.
* `trust_policy_document` - (Required) The trust policy document of the Role.
* `description` - (Optional) The description of the Role.
* `max_session_duration` - (Optional) The max session duration of the Role.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `trn` - The resource name of the Role.


## Import
Iam role can be imported using the id, e.g.
```
$ terraform import vestack_iam_role.default TerraformTestRole
```

