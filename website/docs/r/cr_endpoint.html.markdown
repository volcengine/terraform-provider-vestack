---
subcategory: "CR"
layout: "vestack"
page_title: "Vestack: vestack_cr_endpoint"
sidebar_current: "docs-vestack-resource-cr_endpoint"
description: |-
  Provides a resource to manage cr endpoint
---
# vestack_cr_endpoint
Provides a resource to manage cr endpoint
## Example Usage
```hcl
# endpoint cannot be created,please import by command `terraform import vestack_cr_endpoint.default endpoint:registryId`

resource "vestack_cr_endpoint" "default" {
  registry = "tf-1"
  enabled  = true
}
```
## Argument Reference
The following arguments are supported:
* `registry` - (Required, ForceNew) The CrRegistry name.
* `enabled` - (Optional) Whether enable public endpoint.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `status` - The status of public endpoint.


## Import
CR endpoints can be imported using the endpoint:registryName, e.g.
```
$ terraform import vestack_cr_endpoint.default endpoint:cr-basic
```

