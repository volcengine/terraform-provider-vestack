---
subcategory: "EIP"
layout: "vestack"
page_title: "Vestack: vestack_eip_associate"
sidebar_current: "docs-vestack-resource-eip_associate"
description: |-
  Provides a resource to manage eip associate
---
# vestack_eip_associate
Provides a resource to manage eip associate
## Example Usage
```hcl
data "vestack_zones" "foo" {
}

data "vestack_images" "foo" {
  os_type          = "Linux"
  visibility       = "public"
  instance_type_id = "ecs.g1.large"
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

resource "vestack_security_group" "foo" {
  vpc_id              = vestack_vpc.foo.id
  security_group_name = "acc-test-security-group"
}

resource "vestack_ecs_instance" "foo" {
  image_id             = data.vestack_images.foo.images[0].image_id
  instance_type        = "ecs.g1.large"
  instance_name        = "acc-test-ecs-name"
  password             = "93f0cb0614Aab12"
  instance_charge_type = "PostPaid"
  system_volume_type   = "ESSD_PL0"
  system_volume_size   = 40
  subnet_id            = vestack_subnet.foo.id
  security_group_ids   = [vestack_security_group.foo.id]
}

resource "vestack_eip_address" "foo" {
  billing_type = "PostPaidByTraffic"
}

resource "vestack_eip_associate" "foo" {
  allocation_id = vestack_eip_address.foo.id
  instance_id   = vestack_ecs_instance.foo.id
  instance_type = "EcsInstance"
}
```
## Argument Reference
The following arguments are supported:
* `allocation_id` - (Required, ForceNew) The allocation id of the EIP.
* `instance_id` - (Required, ForceNew) The instance id which be associated to the EIP.
* `instance_type` - (Required, ForceNew) The type of the associated instance,the value is `Nat` or `NetworkInterface` or `ClbInstance` or `EcsInstance` or `HaVip`.
* `private_ip_address` - (Optional, ForceNew) The private IP address of the instance will be associated to the EIP.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
Eip associate can be imported using the eip allocation_id:instance_id, e.g.
```
$ terraform import vestack_eip_associate.default eip-274oj9a8rs9a87fap8sf9515b:i-cm9t9ug9lggu79yr5tcw
```

