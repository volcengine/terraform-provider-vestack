---
subcategory: "CR"
layout: "vestack"
page_title: "Vestack: vestack_cr_endpoints"
sidebar_current: "docs-vestack-datasource-cr_endpoints"
description: |-
  Use this data source to query detailed information of cr endpoints
---
# vestack_cr_endpoints
Use this data source to query detailed information of cr endpoints
## Example Usage
```hcl
data "vestack_cr_endpoints" "foo" {
  registry = "tf-1"
}
```
## Argument Reference
The following arguments are supported:
* `registry` - (Required) The CR instance name.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `endpoints` - The collection of endpoint query.
    * `enabled` - Whether public endpoint is enabled.
    * `registry` - The name of CR instance.
    * `status` - The status of public endpoint.
* `total_count` - The total count of tag query.


