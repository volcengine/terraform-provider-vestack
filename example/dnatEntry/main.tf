data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block = "172.16.0.0/24"
  zone_id = data.vestack_zones.foo.zones[0].id
  vpc_id = vestack_vpc.foo.id
}

resource "vestack_nat_gateway" "foo" {
  vpc_id = vestack_vpc.foo.id
  subnet_id = vestack_subnet.foo.id
  spec = "Small"
  nat_gateway_name = "acc-test-ng"
  description = "acc-test"
  billing_type = "PostPaid"
  project_name = "default"
  tags {
    key = "k1"
    value = "v1"
  }
}

resource "vestack_eip_address" "foo" {
  name = "acc-test-eip"
  description = "acc-test"
  bandwidth = 1
  billing_type = "PostPaidByBandwidth"
  isp = "BGP"
}

resource "vestack_eip_associate" "foo" {
  allocation_id = vestack_eip_address.foo.id
  instance_id = vestack_nat_gateway.foo.id
  instance_type = "Nat"
}

resource "vestack_dnat_entry" "foo" {
  dnat_entry_name = "acc-test-dnat-entry"
  external_ip = vestack_eip_address.foo.eip_address
  external_port = 80
  internal_ip = "172.16.0.10"
  internal_port = 80
  nat_gateway_id = vestack_nat_gateway.foo.id
  protocol = "tcp"
  depends_on = [vestack_eip_associate.foo]
}
