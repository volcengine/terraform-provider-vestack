---
subcategory: "VPC"
layout: "vestack"
page_title: "Vestack: vestack_vpcs"
sidebar_current: "docs-vestack-datasource-vpcs"
description: |-
  Use this data source to query detailed information of vpcs
---
# vestack_vpcs
Use this data source to query detailed information of vpcs
## Example Usage
```hcl
data "vestack_vpcs" "default" {
  ids = ["vpc-mizl7m1kqccg5smt1bdpijuj"]
}
```
## Argument Reference
The following arguments are supported:
* `ids` - (Optional) A list of VPC IDs.
* `name_regex` - (Optional) A Name Regex of Vpc.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `total_count` - The total count of Vpc query.
* `vpcs` - The collection of Vpc query.
  * `account_id` - The account ID of VPC.
  * `associate_cens` - The associate cen list of VPC.
    * `cen_id` - The ID of CEN.
    * `cen_owner_id` - The owner ID of CEN.
    * `cen_status` - The status of CEN.
  * `auxiliary_cidr_blocks` - The auxiliary cidr block list of VPC.
  * `cidr_block` - The cidr block of VPC.
  * `creation_time` - The create time of VPC.
  * `description` - The description of VPC.
  * `dns_servers` - The dns server list of VPC.
  * `nat_gateway_ids` - The nat gateway ID list of VPC.
  * `route_table_ids` - The route table ID list of VPC.
  * `security_group_ids` - The security group ID list of VPC.
  * `status` - The status of VPC.
  * `subnet_ids` - The subnet ID list of VPC.
  * `update_time` - The update time of VPC.
  * `vpc_id` - The ID of VPC.
  * `vpc_name` - The name of VPC.


