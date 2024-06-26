resource "vestack_iam_user" "foo" {
  user_name    = "acc-test-user"
  description  = "acc-test"
  display_name = "name"
}

resource "vestack_iam_access_key" "foo" {
  user_name   = vestack_iam_user.foo.user_name
  secret_file = "./sk"
  status      = "active"
  #  pgp_key = "keybase:some_person_that_exists"
}
