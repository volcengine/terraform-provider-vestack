---
subcategory: "CLB"
layout: "vestack"
page_title: "Vestack: vestack_server_groups"
sidebar_current: "docs-vestack-datasource-server_groups"
description: |-
  Use this data source to query detailed information of server groups
---
# vestack_server_groups
Use this data source to query detailed information of server groups
## Example Usage
```hcl

```
## Argument Reference
The following arguments are supported:
* `ids` - (Optional) A list of ServerGroup IDs.
* `load_balancer_id` - (Optional) The id of the Clb.
* `name_regex` - (Optional) A Name Regex of ServerGroup.
* `output_file` - (Optional) File name where to save data source results.
* `server_group_name` - (Optional) The name of the ServerGroup.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `groups` - The collection of ServerGroup query.
    * `address_ip_version` - The address ip version of the ServerGroup.
    * `create_time` - The create time of the ServerGroup.
    * `description` - The description of the ServerGroup.
    * `id` - The ID of the ServerGroup.
    * `server_group_id` - The ID of the ServerGroup.
    * `server_group_name` - The name of the ServerGroup.
    * `update_time` - The update time of the ServerGroup.
* `total_count` - The total count of ServerGroup query.


