---
subcategory: "CLB"
layout: "vestack"
page_title: "Vestack: vestack_acls"
sidebar_current: "docs-vestack-datasource-acls"
description: |-
  Use this data source to query detailed information of acls
---
# vestack_acls
Use this data source to query detailed information of acls
## Example Usage
```hcl

```
## Argument Reference
The following arguments are supported:
* `acl_name` - (Optional) The name of acl.
* `ids` - (Optional) A list of Acl IDs.
* `name_regex` - (Optional) A Name Regex of Acl.
* `output_file` - (Optional) File name where to save data source results.
* `project_name` - (Optional) The ProjectName of Acl.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `acls` - The collection of Acl query.
    * `acl_entry_count` - The count of acl entry.
    * `acl_id` - The ID of Acl.
    * `acl_name` - The Name of Acl.
    * `create_time` - Creation time of Acl.
    * `description` - The description of Acl.
    * `id` - The ID of Acl.
    * `listeners` - The listeners of Acl.
    * `project_name` - The ProjectName of Acl.
    * `update_time` - Update time of Acl.
* `total_count` - The total count of Acl query.


