---
subcategory: "NAT"
layout: "vestack"
page_title: "Vestack: vestack_nat_gateway"
sidebar_current: "docs-vestack-resource-nat_gateway"
description: |-
  Provides a resource to manage nat gateway
---
# vestack_nat_gateway
Provides a resource to manage nat gateway
## Notice
When Destroy this resource,If the resource charge type is PrePaid,Please unsubscribe the resource 
in  [Vestack Console],when complete console operation,yon can
use 'terraform state rm ${resourceId}' to remove.
## Example Usage
```hcl
data "vestack_zones" "foo" {
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "192.168.0.0/16"
}

resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block  = "192.168.71.0/24"
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
```
## Argument Reference
The following arguments are supported:
* `subnet_id` - (Required, ForceNew) The ID of the Subnet.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `billing_type` - (Optional, ForceNew) The billing type of the NatGateway, the value is `PostPaid` or `PrePaid`.
* `description` - (Optional) The description of the NatGateway.
* `nat_gateway_name` - (Optional) The name of the NatGateway.
* `period` - (Optional, ForceNew) The period of the NatGateway, the valid value range in 1~9 or 12 or 24 or 36. Default value is 12. The period unit defaults to `Month`.This field is only effective when creating a PrePaid NatGateway. When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.
* `project_name` - (Optional) The ProjectName of the NatGateway.
* `spec` - (Optional) The specification of the NatGateway. Optional choice contains `Small`(default), `Medium`, `Large` or leave blank.
* `tags` - (Optional) Tags.

The `tags` object supports the following:

* `key` - (Required) The Key of Tags.
* `value` - (Required) The Value of Tags.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
NatGateway can be imported using the id, e.g.
```
$ terraform import vestack_nat_gateway.default ngw-vv3t043k05sm****
```

