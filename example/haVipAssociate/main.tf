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
  security_group_name = "acc-test-sg"
  vpc_id              = vestack_vpc.foo.id
}

resource "vestack_network_interface" "foo" {
  network_interface_name = "acc-test-eni"
  description            = "acc-test"
  subnet_id              = vestack_subnet.foo.id
  security_group_ids     = [vestack_security_group.foo.id]
  primary_ip_address     = "172.16.0.253"
  port_security_enabled  = false
  private_ip_address     = ["172.16.0.2"]
  project_name           = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "vestack_ha_vip" "foo" {
  ha_vip_name = "acc-test-ha-vip"
  description = "acc-test"
  subnet_id   = vestack_subnet.foo.id
  ip_address  = "172.16.0.5"
}

resource "vestack_ha_vip_associate" "foo" {
  ha_vip_id     = vestack_ha_vip.foo.id
  instance_type = "NetworkInterface"
  instance_id   = vestack_network_interface.foo.id
}
