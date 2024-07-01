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

resource "vestack_ecs_deployment_set" "foo" {
  deployment_set_name = "acc-test-ecs-ds"
  description         = "acc-test"
  granularity         = "switch"
  strategy            = "Availability"
}

resource "vestack_ecs_deployment_set_associate" "foo" {
  deployment_set_id = vestack_ecs_deployment_set.foo.id
  instance_id       = vestack_ecs_instance.foo.id
  depends_on        = [vestack_ecs_instance_state.foo]
}