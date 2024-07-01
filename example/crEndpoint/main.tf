# endpoint cannot be created,please import by command `terraform import vestack_cr_endpoint.default endpoint:registryId`

resource "vestack_cr_endpoint" "default" {
  registry = "tf-1"
  enabled  = true
}