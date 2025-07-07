---
subcategory: "NAT"
layout: "vestack"
page_title: "Vestack: vestack_snat_entry"
sidebar_current: "docs-vestack-resource-snat_entry"
description: |-
  Provides a resource to manage snat entry
---
# vestack_snat_entry
Provides a resource to manage snat entry
## Example Usage
```hcl
data "vestack_zones" "foo" {
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block  = "172.16.0.0/24"
  zone_id     = data.vestack_zones.foo.zones[0].id
  vpc_id      = vestack_vpc.foo.id
}

resource "vestack_nat_gateway" "foo" {
  vpc_id           = vestack_vpc.foo.id
  subnet_id        = vestack_subnet.foo.id
  spec             = "Small"
  nat_gateway_name = "acc-test-ng"
  description      = "acc-test"
  billing_type     = "PostPaid"
  project_name     = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "vestack_eip_address" "foo" {
  name         = "acc-test-eip"
  description  = "acc-test"
  bandwidth    = 1
  billing_type = "PostPaidByBandwidth"
  isp          = "BGP"
}

resource "vestack_eip_associate" "foo" {
  allocation_id = vestack_eip_address.foo.id
  instance_id   = vestack_nat_gateway.foo.id
  instance_type = "Nat"
}

resource "vestack_snat_entry" "foo" {
  snat_entry_name = "acc-test-snat-entry"
  nat_gateway_id  = vestack_nat_gateway.foo.id
  eip_id          = vestack_eip_address.foo.id
  source_cidr     = "172.16.0.0/24"
  depends_on      = [vestack_eip_associate.foo]
}
```
## Argument Reference
The following arguments are supported:
* `eip_id` - (Required) The id of the public ip address used by the SNAT entry.
* `nat_gateway_id` - (Required, ForceNew) The id of the nat gateway to which the entry belongs.
* `snat_entry_name` - (Optional) The name of the SNAT entry.
* `source_cidr` - (Optional, ForceNew) The SourceCidr of the SNAT entry. Only one of `subnet_id,source_cidr` can be specified.
* `subnet_id` - (Optional, ForceNew) The id of the subnet that is required to access the internet. Only one of `subnet_id,source_cidr` can be specified.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `status` - The status of the SNAT entry.


## Import
Snat entry can be imported using the id, e.g.
```
$ terraform import vestack_snat_entry.default snat-3fvhk47kf56****
```

