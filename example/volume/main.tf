# query available zones in current region
data "vestack_zones" "foo" {
}

# create vpc
resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

# create subnet
resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block  = "172.16.0.0/24"
  zone_id     = data.vestack_zones.foo.zones[0].id
  vpc_id      = vestack_vpc.foo.id
}

# create security group
resource "vestack_security_group" "foo" {
  security_group_name = "acc-test-security-group"
  vpc_id              = vestack_vpc.foo.id
}

# query the image_id which match the specified instance_type
data "vestack_images" "foo" {
  os_type          = "Linux"
  visibility       = "public"
  instance_type_id = "ecs.g3il.large"
}

# create PrePaid ecs instance
resource "vestack_ecs_instance" "foo" {
  instance_name        = "acc-test-ecs"
  description          = "acc-test"
  host_name            = "tf-acc-test"
  image_id             = data.vestack_images.foo.images[0].image_id
  instance_type        = "ecs.g3il.large"
  password             = "93f0cb0614Aab12"
  instance_charge_type = "PrePaid"
  period               = 1
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

# create PrePaid data volume
resource "vestack_volume" "PreVolume" {
  volume_name          = "acc-test-volume"
  volume_type          = "ESSD_PL0"
  description          = "acc-test"
  kind                 = "data"
  size                 = 40
  zone_id              = data.vestack_zones.foo.zones[0].id
  volume_charge_type   = "PrePaid"
  instance_id          = vestack_ecs_instance.foo.id
  project_name         = "default"
  delete_with_instance = true
  tags {
    key   = "k1"
    value = "v1"
  }
}

# create PostPaid data volume
resource "vestack_volume" "PostVolume" {
  volume_name        = "acc-test-volume"
  volume_type        = "ESSD_PL0"
  description        = "acc-test"
  kind               = "data"
  size               = 40
  zone_id            = data.vestack_zones.foo.zones[0].id
  volume_charge_type = "PostPaid"
  project_name       = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

# attach PostPaid data volume to ecs instance
resource "vestack_volume_attach" "foo" {
  instance_id = vestack_ecs_instance.foo.id
  volume_id   = vestack_volume.PostVolume.id
}
