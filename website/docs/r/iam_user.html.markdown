---
subcategory: "IAM"
layout: "vestack"
page_title: "Vestack: vestack_iam_user"
sidebar_current: "docs-vestack-resource-iam_user"
description: |-
  Provides a resource to manage iam user
---
# vestack_iam_user
Provides a resource to manage iam user
## Example Usage
```hcl
resource "vestack_iam_user" "foo" {
  user_name    = "tf-test"
  description  = "test"
  display_name = "name"
}
```
## Argument Reference
The following arguments are supported:
* `user_name` - (Required) The name of the user.
* `description` - (Optional) The description of the user.
* `display_name` - (Optional) The display name of the user.
* `email` - (Optional) The email of the user.
* `mobile_phone` - (Optional) The mobile phone of the user.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `account_id` - The account id of the user.
* `create_date` - The create date of the user.
* `trn` - The trn of the user.
* `update_date` - The update date of the user.


## Import
Iam user can be imported using the UserName, e.g.
```
$ terraform import vestack_iam_user.default user_name
```

