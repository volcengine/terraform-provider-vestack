---
subcategory: "CR"
layout: "vestack"
page_title: "Vestack: vestack_cr_namespace"
sidebar_current: "docs-vestack-resource-cr_namespace"
description: |-
  Provides a resource to manage cr namespace
---
# vestack_cr_namespace
Provides a resource to manage cr namespace
## Example Usage
```hcl
resource "vestack_cr_namespace" "foo" {
  registry = "tf-2"
  name     = "namespace-1"
}

resource "vestack_cr_namespace" "foo1" {
  registry = "tf-1"
  name     = "namespace-2"
}
```
## Argument Reference
The following arguments are supported:
* `name` - (Required, ForceNew) The name of CrNamespace.
* `registry` - (Required, ForceNew) The registry name.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `create_time` - The time when namespace created.


## Import
CR namespace can be imported using the registry:name, e.g.
```
$ terraform import vestack_cr_namespace.default cr-basic:namespace-1
```

