---
subcategory: "VPC"
layout: "vestack"
page_title: "Vestack: vestack_network_interface"
sidebar_current: "docs-vestack-resource-network_interface"
description: |-
  Provides a resource to manage network interface
---
# vestack_network_interface
Provides a resource to manage network interface
## Example Usage
```hcl
resource "vestack_network_interface" "foo" {
  subnet_id              = "subnet-2fe79j7c8o5c059gp68ksxr93"
  security_group_ids     = ["sg-2fepz3c793g1s59gp67y21r34"]
  primary_ip_address     = "192.168.5.253"
  network_interface_name = "tf-test-up"
  description            = "tf-test-up"
  port_security_enabled  = false
  project_name           = "default"
  private_ip_address     = ["192.168.5.2"]
  //secondary_private_ip_address_count = 0
}
```
## Argument Reference
The following arguments are supported:
* `security_group_ids` - (Required) The list of the security group id to which the secondary ENI belongs.
* `subnet_id` - (Required, ForceNew) The id of the subnet to which the ENI is connected.
* `description` - (Optional) The description of the ENI.
* `network_interface_name` - (Optional) The name of the ENI.
* `port_security_enabled` - (Optional) Set port security enable or disable.
* `primary_ip_address` - (Optional, ForceNew) The primary IP address of the ENI.
* `private_ip_address` - (Optional) The list of private ip address. This field conflicts with `secondary_private_ip_address_count`.
* `project_name` - (Optional) The ProjectName of the ENI.
* `secondary_private_ip_address_count` - (Optional) The count of secondary private ip address. This field conflicts with `private_ip_address`.
* `tags` - (Optional) Tags.

The `tags` object supports the following:

* `key` - (Required) The Key of Tags.
* `value` - (Required) The Value of Tags.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `status` - The status of the ENI.


## Import
Network interface can be imported using the id, e.g.
```
$ terraform import vestack_network_interface.default eni-bp1fgnh68xyz9****
```

