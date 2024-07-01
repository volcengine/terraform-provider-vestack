---
subcategory: "VPC"
layout: "vestack"
page_title: "Vestack: vestack_route_table"
sidebar_current: "docs-vestack-resource-route_table"
description: |-
  Provides a resource to manage route table
---
# vestack_route_table
Provides a resource to manage route table
## Example Usage
```hcl
resource "vestack_route_table" "foo" {
  vpc_id           = "vpc-2feppmy1ugt1c59gp688n1fld"
  route_table_name = "tf-project-1"
  description      = "tf-test1"
  project_name     = "yuwao"
}
```
## Argument Reference
The following arguments are supported:
* `vpc_id` - (Required, ForceNew) The id of the VPC.
* `description` - (Optional) The description of the route table.
* `project_name` - (Optional) The ProjectName of the route table.
* `route_table_name` - (Optional) The name of the route table.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
Route table can be imported using the id, e.g.
```
$ terraform import vestack_route_table.default vtb-274e0syt9av407fap8tle16kb
```

