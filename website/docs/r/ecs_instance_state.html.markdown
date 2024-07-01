---
subcategory: "ECS"
layout: "vestack"
page_title: "Vestack: vestack_ecs_instance_state"
sidebar_current: "docs-vestack-resource-ecs_instance_state"
description: |-
  Provides a resource to manage ecs instance state
---
# vestack_ecs_instance_state
Provides a resource to manage ecs instance state
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

resource "vestack_security_group" "foo" {
  security_group_name = "acc-test-security-group"
  vpc_id              = vestack_vpc.foo.id
}

data "vestack_images" "foo" {
  os_type          = "Linux"
  visibility       = "public"
  instance_type_id = "ecs.g1.large"
}

resource "vestack_ecs_instance" "foo" {
  instance_name        = "acc-test-ecs"
  image_id             = data.vestack_images.foo.images[0].image_id
  instance_type        = "ecs.g1.large"
  password             = "93f0cb0614Aab12"
  instance_charge_type = "PostPaid"
  system_volume_type   = "ESSD_PL0"
  system_volume_size   = 40
  subnet_id            = vestack_subnet.foo.id
  security_group_ids   = [vestack_security_group.foo.id]
}

resource "vestack_ecs_instance_state" "foo" {
  instance_id  = vestack_ecs_instance.foo.id
  action       = "Stop"
  stopped_mode = "KeepCharging"
}
```
## Argument Reference
The following arguments are supported:
* `action` - (Required) Start or Stop of Instance Action, the value can be `Start`, `Stop` or `ForceStop`.
* `instance_id` - (Required, ForceNew) Id of Instance.
* `stopped_mode` - (Optional) Stop Mode of Instance, the value can be `KeepCharging` or `StopCharging`, default `KeepCharging`.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `status` - Status of Instance.


## Import
State Instance can be imported using the id, e.g.
```
$ terraform import vestack_ecs_instance_state.default state:i-mizl7m1kqccg5smt1bdpijuj
```

