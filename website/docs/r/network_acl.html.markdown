---
subcategory: "VPC"
layout: "vestack"
page_title: "Vestack: vestack_network_acl"
sidebar_current: "docs-vestack-resource-network_acl"
description: |-
  Provides a resource to manage network acl
---
# vestack_network_acl
Provides a resource to manage network acl
## Example Usage
```hcl
resource "vestack_network_acl" "foo" {
  vpc_id           = "vpc-2d6jskar243k058ozfdae13ne"
  network_acl_name = "tf-test-acl"

  ingress_acl_entries {
    network_acl_entry_name = "ingress1"
    policy                 = "accept"
    protocol               = "all"
    source_cidr_ip         = "192.168.0.0/24"
  }

  egress_acl_entries {
    network_acl_entry_name = "egress2"
    policy                 = "accept"
    protocol               = "all"
    destination_cidr_ip    = "192.168.0.0/16"
  }

  ingress_acl_entries {
    network_acl_entry_name = "ingress3"
    policy                 = "accept"
    protocol               = "tcp"
    port                   = "80/80"
    source_cidr_ip         = "192.168.0.0/24"
  }

  project_name = "default"
}
```
## Argument Reference
The following arguments are supported:
* `vpc_id` - (Required, ForceNew) The vpc id of Network Acl.
* `description` - (Optional) The description of the Network Acl.
* `egress_acl_entries` - (Optional) The egress entries of Network Acl.
* `ingress_acl_entries` - (Optional) The ingress entries of Network Acl.
* `network_acl_name` - (Optional) The name of Network Acl.
* `project_name` - (Optional) The project name of the network acl.

The `egress_acl_entries` object supports the following:

* `description` - (Optional) The description of entry.
* `destination_cidr_ip` - (Optional) The DestinationCidrIp of entry.
* `network_acl_entry_name` - (Optional) The name of entry.
* `policy` - (Optional) The policy of entry. Default is `accept`. The value can be `accept` or `drop`.
* `port` - (Optional) The port of entry. Default is `-1/-1`. When Protocol is `all`, `icmp` or `gre`, the port range is `-1/-1`, which means no port restriction.When the Protocol is `tcp` or `udp`, the port range is `1~65535`, and the format is `1/200`, `80/80`,which means port 1 to port 200, port 80.
* `protocol` - (Optional) The protocol of entry. The value can be `icmp` or `gre` or `tcp` or `udp` or `all`. Default is `all`.

The `ingress_acl_entries` object supports the following:

* `description` - (Optional) The description of entry.
* `network_acl_entry_name` - (Optional) The name of entry.
* `policy` - (Optional) The policy of entry, default is `accept`. The value can be `accept` or `drop`.
* `port` - (Optional) The port of entry. Default is `-1/-1`. When Protocol is `all`, `icmp` or `gre`, the port range is `-1/-1`, which means no port restriction. When the Protocol is `tcp` or `udp`, the port range is `1~65535`, and the format is `1/200`, `80/80`, which means port 1 to port 200, port 80.
* `protocol` - (Optional) The protocol of entry, default is `all`. The value can be `icmp` or `gre` or `tcp` or `udp` or `all`.
* `source_cidr_ip` - (Optional) The SourceCidrIp of entry.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
Network Acl can be imported using the id, e.g.
```
$ terraform import vestack_network_acl.default nacl-172leak37mi9s4d1w33pswqkh
```

