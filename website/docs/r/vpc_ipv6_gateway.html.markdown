---
subcategory: "VPC"
layout: "vestack"
page_title: "Vestack: vestack_vpc_ipv6_gateway"
sidebar_current: "docs-vestack-resource-vpc_ipv6_gateway"
description: |-
  Provides a resource to manage vpc ipv6 gateway
---
# vestack_vpc_ipv6_gateway
Provides a resource to manage vpc ipv6 gateway
## Example Usage
```hcl
resource "vestack_vpc_ipv6_gateway" "foo" {
  vpc_id      = "vpc-12afxho4sxyio17q7y2kkp8ej"
  name        = "tf-test-1"
  description = "test"
}
```
## Argument Reference
The following arguments are supported:
* `vpc_id` - (Required, ForceNew) The ID of the VPC which the Ipv6Gateway belongs to.
* `description` - (Optional) The description of the Ipv6Gateway.
* `name` - (Optional) The name of the Ipv6Gateway.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `creation_time` - Creation time of the Ipv6Gateway.
* `ipv6_gateway_id` - The ID of the Ipv6Gateway.
* `status` - The Status of the Ipv6Gateway.
* `update_time` - Update time of the Ipv6Gateway.


## Import
Ipv6Gateway can be imported using the id, e.g.
```
$ terraform import vestack_vpc_ipv6_gateway.default ipv6gw-12bcapllb5ukg17q7y2sd3thx
```

