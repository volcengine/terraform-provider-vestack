---
subcategory: "IAM"
layout: "vestack"
page_title: "Vestack: vestack_iam_login_profile"
sidebar_current: "docs-vestack-resource-iam_login_profile"
description: |-
  Provides a resource to manage iam login profile
---
# vestack_iam_login_profile
Provides a resource to manage iam login profile
## Example Usage
```hcl
resource "vestack_iam_user" "foo" {
  user_name    = "acc-test-user"
  description  = "acc-test"
  display_name = "name"
}

resource "vestack_iam_login_profile" "foo" {
  user_name               = vestack_iam_user.foo.user_name
  password                = "93f0cb0614Aab12"
  login_allowed           = true
  password_reset_required = false
}
```
## Argument Reference
The following arguments are supported:
* `password` - (Required) The password.
* `user_name` - (Required, ForceNew) The user name.
* `login_allowed` - (Optional) The flag of login allowed.
* `password_reset_required` - (Optional) Is required reset password when next time login in.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
Login profile can be imported using the UserName, e.g.
```
$ terraform import vestack_iam_login_profile.default user_name
```

