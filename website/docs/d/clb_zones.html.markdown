---
subcategory: "CLB"
layout: "vestack"
page_title: "Vestack: vestack_clb_zones"
sidebar_current: "docs-vestack-datasource-clb_zones"
description: |-
  Use this data source to query detailed information of clb zones
---
# vestack_clb_zones
Use this data source to query detailed information of clb zones
## Example Usage
```hcl

```
## Argument Reference
The following arguments are supported:
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `master_zones` - The master zones list.
    * `slave_zones` - The slave zones list.
        * `zone_id` - The slave zone id.
    * `zone_id` - The master zone id.
* `total_count` - The total count of zone query.


