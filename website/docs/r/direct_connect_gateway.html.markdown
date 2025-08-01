---
subcategory: "DIRECT_CONNECT"
layout: "vestack"
page_title: "Vestack: vestack_direct_connect_gateway"
sidebar_current: "docs-vestack-resource-direct_connect_gateway"
description: |-
  Provides a resource to manage direct connect gateway
---
# vestack_direct_connect_gateway
Provides a resource to manage direct connect gateway
## Example Usage
```hcl
resource "vestack_direct_connect_gateway" "foo" {
  direct_connect_gateway_name = "tf-test-gateway"
  description                 = "tf-test"
  tags {
    key   = "k1"
    value = "v1"
  }
}
```
## Argument Reference
The following arguments are supported:
* `description` - (Optional) The description of direct connect gateway.
* `direct_connect_gateway_name` - (Optional) The name of direct connect gateway.
* `tags` - (Optional) The direct connect gateway tags.

The `tags` object supports the following:

* `key` - (Required) The tag key.
* `value` - (Optional) The tag value.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
DirectConnectGateway can be imported using the id, e.g.
```
$ terraform import vestack_direct_connect_gateway.default resource_id
```

