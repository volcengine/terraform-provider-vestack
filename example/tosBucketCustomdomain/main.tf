# Create a custom domain for TOS bucket
resource "vestack_tos_bucket_customdomain" "default" {
  bucket_name = "tflyb7"
  custom_domain_rule {
    domain   = "www.163.com"
    protocol = "tos"
  }
}

resource "vestack_tos_bucket_customdomain" "default1" {
  bucket_name = "tflyb7"
  custom_domain_rule {
    domain   = "www.2345.com"
    protocol = "tos"
  }
}