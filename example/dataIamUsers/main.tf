resource "vestack_iam_user" "foo" {
  user_name    = "acc-test-user"
  description  = "acc test"
  display_name = "name"
}
data "vestack_iam_users" "foo" {
  user_names = [vestack_iam_user.foo.user_name]
}