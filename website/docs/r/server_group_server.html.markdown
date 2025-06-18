---
subcategory: "CLB"
layout: "vestack"
page_title: "Vestack: vestack_server_group_server"
sidebar_current: "docs-vestack-resource-server_group_server"
description: |-
  Provides a resource to manage server group server
---
# vestack_server_group_server
Provides a resource to manage server group server
## Example Usage
```hcl

```
## Argument Reference
The following arguments are supported:
* `instance_id` - (Required, ForceNew) The ID of ecs instance or the network card bound to ecs instance.
* `port` - (Required) The port receiving request.
* `server_group_id` - (Required, ForceNew) The ID of the ServerGroup.
* `type` - (Required, ForceNew) The type of instance. Optional choice contains `ecs`, `eni`.
* `description` - (Optional) The description of the instance.
* `ip` - (Optional, ForceNew) The private ip of the instance.
* `weight` - (Optional) The weight of the instance, range in 0~100.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `server_id` - The server id of instance in ServerGroup.


## Import
ServerGroupServer can be imported using the id, e.g.
```
$ terraform import vestack_server_group_server.default rsp-274xltv2*****8tlv3q3s:rs-3ciynux6i1x4w****rszh49sj
```

