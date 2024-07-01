resource "vestack_ecs_key_pair" "foo" {
  key_pair_name = "acc-test-key-name"
  description   = "acc-test"
}

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
  password             = "your password"
  instance_charge_type = "PostPaid"
  system_volume_type   = "ESSD_PL0"
  system_volume_size   = 40
  subnet_id            = vestack_subnet.foo.id
  security_group_ids   = [vestack_security_group.foo.id]
}

resource "vestack_ecs_key_pair_associate" "foo" {
  instance_id = vestack_ecs_instance.foo.id
  key_pair_id = vestack_ecs_key_pair.foo.id
}