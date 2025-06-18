---
subcategory: "CLB"
layout: "vestack"
page_title: "Vestack: vestack_server_group"
sidebar_current: "docs-vestack-resource-server_group"
description: |-
  Provides a resource to manage server group
---
# vestack_server_group
Provides a resource to manage server group
## Example Usage
```hcl

```
## Argument Reference
The following arguments are supported:
* `load_balancer_id` - (Required, ForceNew) The ID of the Clb.
* `address_ip_version` - (Optional, ForceNew) The address ip version of the ServerGroup. Valid values: `ipv4`, `ipv6`. Default is `ipv4`.
* `description` - (Optional) The description of ServerGroup.
* `server_group_id` - (Optional) The ID of the ServerGroup.
* `server_group_name` - (Optional) The name of the ServerGroup.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
ServerGroup can be imported using the id, e.g.
```
$ terraform import vestack_server_group.default rsp-273yv0kir1vk07fap8tt9jtwg
```

