resource "vestack_ecs_deployment_set" "foo" {
  deployment_set_name = "acc-test-ecs-ds-${count.index}"
  description         = "acc-test"
  granularity         = "switch"
  strategy            = "Availability"
  count               = 3
}

data "vestack_ecs_deployment_sets" "foo" {
  granularity = "switch"
  ids         = vestack_ecs_deployment_set.foo[*].id
}
