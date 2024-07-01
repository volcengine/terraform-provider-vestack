---
subcategory: "TOS(BETA)"
layout: "vestack"
page_title: "Vestack: vestack_tos_objects"
sidebar_current: "docs-vestack-datasource-tos_objects"
description: |-
  Use this data source to query detailed information of tos objects
---
# vestack_tos_objects
Use this data source to query detailed information of tos objects
## Example Usage
```hcl
data "vestack_tos_objects" "default" {
  bucket_name = "test"
}
```
## Argument Reference
The following arguments are supported:
* `bucket_name` - (Required) The name the TOS bucket.
* `name_regex` - (Optional) A Name Regex of TOS Object.
* `object_name` - (Optional) The name the TOS Object.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `objects` - The collection of TOS Object query.
    * `content` - The content the TOS Object when content type is json or text and xml.
    * `name` - The name the TOS Object.
    * `size` - The name the TOS Object size.
    * `storage_class` - The name the TOS Object storage class.
* `total_count` - The total count of TOS Object query.


