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

// create internet nat gateway and snat entry and dnat entry
resource "vestack_nat_gateway" "internet_nat_gateway" {
  vpc_id           = vestack_vpc.foo.id
  subnet_id        = vestack_subnet.foo.id
  spec             = "Small"
  nat_gateway_name = "acc-test-internet_ng"
  description      = "acc-test"
  billing_type     = "PostPaid"
  project_name     = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "vestack_eip_address" "foo" {
  name         = "acc-test-eip"
  description  = "acc-test"
  bandwidth    = 1
  billing_type = "PostPaidByBandwidth"
  isp          = "BGP"
}

resource "vestack_eip_associate" "foo" {
  allocation_id = vestack_eip_address.foo.id
  instance_id   = vestack_nat_gateway.internet_nat_gateway.id
  instance_type = "Nat"
}

resource "vestack_snat_entry" "foo" {
  snat_entry_name = "acc-test-snat-entry"
  nat_gateway_id  = vestack_nat_gateway.internet_nat_gateway.id
  eip_id          = vestack_eip_address.foo.id
  source_cidr     = "172.16.0.0/24"
  depends_on      = [vestack_eip_associate.foo]
}

resource "vestack_dnat_entry" "foo" {
  dnat_entry_name = "acc-test-dnat-entry"
  external_ip     = vestack_eip_address.foo.eip_address
  external_port   = 80
  internal_ip     = "172.16.0.10"
  internal_port   = 80
  nat_gateway_id  = vestack_nat_gateway.internet_nat_gateway.id
  protocol        = "tcp"
  depends_on      = [vestack_eip_associate.foo]
}

// create intranet nat gateway and snat entry and dnat entry
resource "vestack_nat_gateway" "intranet_nat_gateway" {
  vpc_id           = vestack_vpc.foo.id
  subnet_id        = vestack_subnet.foo.id
  nat_gateway_name = "acc-test-intranet_ng"
  description      = "acc-test"
  network_type     = "intranet"
  billing_type     = "PostPaidByUsage"
  project_name     = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "vestack_nat_ip" "foo" {
  nat_gateway_id     = vestack_nat_gateway.intranet_nat_gateway.id
  nat_ip_name        = "acc-test-nat-ip"
  nat_ip_description = "acc-test"
  nat_ip             = "172.16.0.3"
}

resource "vestack_snat_entry" "foo-intranet" {
  snat_entry_name = "acc-test-snat-entry-intranet"
  nat_gateway_id  = vestack_nat_gateway.intranet_nat_gateway.id
  nat_ip_id       = vestack_nat_ip.foo.id
  source_cidr     = "172.16.0.0/24"
}

resource "vestack_dnat_entry" "foo-intranet" {
  nat_gateway_id  = vestack_nat_gateway.intranet_nat_gateway.id
  dnat_entry_name = "acc-test-dnat-entry-intranet"
  protocol        = "tcp"
  internal_ip     = "172.16.0.5"
  internal_port   = "82"
  external_ip     = vestack_nat_ip.foo.nat_ip
  external_port   = "87"
}
