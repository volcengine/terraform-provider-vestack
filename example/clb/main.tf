# query available zones in current region
data "vestack_zones" "foo" {
}

# create vpc
resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "192.168.0.0/16"
  project_name = "default"
}

# create subnet
resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block  = "192.168.0.0/24"
  zone_id     = data.vestack_zones.foo.zones[0].id
  vpc_id      = vestack_vpc.foo.id
}

# create ipv4 public clb
resource "vestack_clb" "public_clb" {
  type               = "public"
  subnet_id          = vestack_subnet.foo.id
  load_balancer_name = "acc-test-clb-public"
  load_balancer_spec = "small_1"
  description        = "acc-test-demo"
  project_name       = "default"
  eip_billing_config {
    isp              = "ChinaUnicom"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth        = 1
  }
  tags {
    key   = "k1"
    value = "v1"
  }
}

# create ipv4 private clb
resource "vestack_clb" "private_clb" {
  type               = "private"
  subnet_id          = vestack_subnet.foo.id
  load_balancer_name = "acc-test-clb-private"
  load_balancer_spec = "small_1"
  description        = "acc-test-demo"
  project_name       = "default"
}

# create eip
resource "vestack_eip_address" "eip" {
  billing_type = "PostPaidByBandwidth"
  bandwidth    = 1
  isp          = "ChinaUnicom"
  name         = "tf-eip"
  description  = "tf-test"
  project_name = "default"
}

# associate eip to clb
resource "vestack_eip_associate" "associate" {
  allocation_id = vestack_eip_address.eip.id
  instance_id   = vestack_clb.private_clb.id
  instance_type = "ClbInstance"
}

# create ipv6 vpc
resource "vestack_vpc" "vpc_ipv6" {
  vpc_name    = "acc-test-vpc-ipv6"
  cidr_block  = "192.168.0.0/16"
  enable_ipv6 = true
  project_name = "default"
  #ipv6_cidr_block = "fa00:230:0:de00::/56"
  ipv6_cidr_block_type = "ULA"

}

# create ipv6 subnet
resource "vestack_subnet" "subnet_ipv6" {
  subnet_name     = "acc-test-subnet-ipv6"
  cidr_block      = "192.168.0.0/24"
  zone_id         = data.vestack_zones.foo.zones[0].id
  vpc_id          = vestack_vpc.vpc_ipv6.id
  ipv6_cidr_block = 1
}

# create ipv6 private clb
resource "vestack_clb" "private_clb_ipv6" {
  type               = "private"
  subnet_id          = vestack_subnet.subnet_ipv6.id
  load_balancer_name = "acc-test-clb-ipv6"
  load_balancer_spec = "small_1"
  description        = "acc-test-demo"
  project_name       = "default"
  address_ip_version = "DualStack"
}

# create ipv6 gateway
resource "vestack_vpc_ipv6_gateway" "ipv6_gateway" {
  vpc_id = vestack_vpc.vpc_ipv6.id
  name   = "acc-test-ipv6-gateway"
}

resource "vestack_vpc_ipv6_address_bandwidth" "foo" {
  ipv6_address = vestack_clb.private_clb_ipv6.eni_ipv6_address
  billing_type = "PostPaidByBandwidth"
  bandwidth    = 5
  depends_on   = [vestack_vpc_ipv6_gateway.ipv6_gateway]
}
