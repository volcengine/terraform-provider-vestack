# Tag cannot be created,please import by command `terraform import vestack_cr_tag.default registry:namespace:repository:tag`
resource "vestack_cr_tag" "default" {
  registry   = "enterprise-1"
  namespace  = "langyu"
  repository = "repo"
  name       = "v2"
}