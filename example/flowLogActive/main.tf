data "vestack_zones" "foo" {
}

resource "vestack_vpc" "foo" {
  vpc_name     = "acc-test-vpc"
  cidr_block   = "172.16.0.0/16"
  project_name = "default"
}

resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block  = "172.16.0.0/24"
  zone_id     = data.vestack_zones.foo.zones[0].id
  vpc_id      = vestack_vpc.foo.id
}

resource "vestack_flow_log" "foo" {
  flow_log_name        = "acc-test-flow-log"
  description          = "acc-test"
  resource_type        = "subnet"
  resource_id          = vestack_subnet.foo.id
  traffic_type         = "All"
  log_project_name     = "acc-test-project"
  log_topic_name       = "acc-test-topic"
  aggregation_interval = 10
  project_name         = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "vestack_flow_log_active" "foo" {
  flow_log_id = vestack_flow_log.foo.id
}