data "vestack_ecs_instances" "dataEcs" {
  ids = ["i-ycal1mtpucl8j0hjiihy"]
}

data "vestack_vpc_ipv6_addresses" "dataIpv6" {
  associated_instance_id = data.vestack_ecs_instances.dataEcs.instances.0.instance_id
}

resource "vestack_vpc_ipv6_address_bandwidth" "foo" {
  ipv6_address = data.vestack_vpc_ipv6_addresses.dataIpv6.ipv6_addresses.0.ipv6_address
  billing_type = "PostPaidByBandwidth"
  bandwidth    = 5
}