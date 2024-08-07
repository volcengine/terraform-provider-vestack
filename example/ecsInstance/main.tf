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
  data_volumes {
    volume_type          = "ESSD_PL0"
    size                 = 50
    delete_with_instance = true
  }
  subnet_id          = vestack_subnet.foo.id
  security_group_ids = [vestack_security_group.foo.id]

  #  deployment_set_id = ""
  #  ipv6_address_count = 1
  #  secondary_network_interfaces {
  #    subnet_id = vestack_subnet.foo.id
  #    security_group_ids = [vestack_security_group.foo.id]
  #  }

  project_name = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
  ha_strategy = "offsite_rebuild"
}
