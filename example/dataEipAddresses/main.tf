resource "vestack_eip_address" "foo" {
  billing_type = "PostPaidByTraffic"
}
data "vestack_eip_addresses" "foo" {
  ids = [vestack_eip_address.foo.id]
}