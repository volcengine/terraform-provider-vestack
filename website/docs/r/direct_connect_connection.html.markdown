---
subcategory: "DIRECT_CONNECT"
layout: "vestack"
page_title: "Vestack: vestack_direct_connect_connection"
sidebar_current: "docs-vestack-resource-direct_connect_connection"
description: |-
  Provides a resource to manage direct connect connection
---
# vestack_direct_connect_connection
Provides a resource to manage direct connect connection
## Example Usage
```hcl
# 混合云场景下物理专线是通过OPSAPI创建的，通过Terraform创建物理专线时，provider需要使用运维端账号的AKSK
resource "vestack_direct_connect_connection" "foo" {
  direct_connect_connection_name = "tf-test-connection"
  description                    = "tf-test"
  port_id                        = "dcp-xxxxx"
  # 运维端不持有资源，owner_account_id填入租户的account_id
  owner_account_id       = "1000000xxx"
  owner_project_name     = "default"
  line_operator          = "ChinaOther"
  port_type              = "10GBase"
  port_spec              = "10G"
  bandwidth              = 1000
  peer_location          = "XX路XX号XX楼XX机房"
  customer_name          = "tf-a"
  customer_contact_phone = "12345678911"
  customer_contact_email = "email@aaa.com"
}
```
## Argument Reference
The following arguments are supported:
* `bandwidth` - (Required, ForceNew) The line band width,unit:Mbps.
* `customer_contact_email` - (Required) The dedicated line contact email.
* `customer_contact_phone` - (Required) The dedicated line contact phone.
* `customer_name` - (Required) The dedicated line contact name.
* `line_operator` - (Required, ForceNew) The physical leased line operator.valid value contains `ChinaTelecom`,`ChinaMobile`,`ChinaUnicom`,`ChinaOther`.
* `owner_account_id` - (Required, ForceNew) The direct connect connection owner account id.
* `owner_project_name` - (Required, ForceNew) The direct connect connection owner project name.
* `peer_location` - (Required, ForceNew) The local IDC address.
* `port_id` - (Required, ForceNew) The direct connect access point port id.
* `port_spec` - (Required, ForceNew) The physical leased line port spec.valid value contains `1G`,`10G`.
* `port_type` - (Required, ForceNew) The physical leased line port type and spec.valid value contains `1000Base-T`,`10GBase-T`,`1000Base`,`10GBase`,`40GBase`,`100GBase`.
* `description` - (Optional) The description of direct connect.
* `direct_connect_connection_name` - (Optional) The name of direct connect.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
DirectConnectConnection can be imported using the id, e.g.
```
$ terraform import vestack_direct_connect_connection.default dcc-7qthudw0ll6jmc****
```

