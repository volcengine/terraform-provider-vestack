---
subcategory: "EBS"
layout: "vestack"
page_title: "Vestack: vestack_volume_attach"
sidebar_current: "docs-vestack-resource-volume_attach"
description: |-
  Provides a resource to manage volume attach
---
# vestack_volume_attach
Provides a resource to manage volume attach
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
  description          = "acc-test"
  host_name            = "tf-acc-test"
  image_id             = data.vestack_images.foo.images[0].image_id
  instance_type        = "ecs.g1.large"
  password             = "93f0cb0614Aab12"
  instance_charge_type = "PostPaid"
  system_volume_type   = "ESSD_PL0"
  system_volume_size   = 40
  subnet_id            = vestack_subnet.foo.id
  security_group_ids   = [vestack_security_group.foo.id]
  project_name         = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "vestack_volume" "foo" {
  volume_name        = "acc-test-volume"
  volume_type        = "ESSD_PL0"
  description        = "acc-test"
  kind               = "data"
  size               = 40
  zone_id            = data.vestack_zones.foo.zones[0].id
  volume_charge_type = "PostPaid"
  project_name       = "default"
}

resource "vestack_volume_attach" "foo" {
  instance_id = vestack_ecs_instance.foo.id
  volume_id   = vestack_volume.foo.id
}
```
## Argument Reference
The following arguments are supported:
* `instance_id` - (Required, ForceNew) The Id of Instance.
* `volume_id` - (Required, ForceNew) The Id of Volume.
* `delete_with_instance` - (Optional) Delete Volume with Attached Instance.It is not recommended to use this field. If used, please ensure that the value of this field is consistent with the value of `delete_with_instance` in vestack_volume.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `created_at` - Creation time of Volume.
* `status` - Status of Volume.
* `updated_at` - Update time of Volume.


## Import
VolumeAttach can be imported using the id, e.g.
```
$ terraform import vestack_volume_attach.default vol-abc12345:i-abc12345
```

