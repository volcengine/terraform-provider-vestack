---
subcategory: "IAM"
layout: "vestack"
page_title: "Vestack: vestack_iam_service_linked_role"
sidebar_current: "docs-vestack-resource-iam_service_linked_role"
description: |-
  Provides a resource to manage iam service linked role
---
# vestack_iam_service_linked_role
Provides a resource to manage iam service linked role
## Example Usage
```hcl
resource "vestack_iam_service_linked_role" "foo" {
  service_name = "transitrouter"
}
```
## Argument Reference
The following arguments are supported:
* `service_name` - (Required, ForceNew) The name of the service.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `description` - The description of the service linked Role.
* `display_name` - The display name of the service linked Role.
* `max_session_duration` - The max session duration of the service linked Role.
* `role_name` - The name of the service linked role.
* `trn` - The resource name of the service linked Role.
* `trust_policy_document` - The trust policy document of the service linked Role.


## Import
Iam service linked role can be imported using the servicx name and the service linked role name, e.g.
```
$ terraform import vestack_iam_service_linked_role.default ecs:ServiceRoleForEcs
```

