data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "192.168.0.0/16"
}

resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block = "192.168.71.0/24"
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
