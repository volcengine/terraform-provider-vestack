---
subcategory: "CLB"
layout: "vestack"
page_title: "Vestack: vestack_clb_rule"
sidebar_current: "docs-vestack-resource-clb_rule"
description: |-
  Provides a resource to manage clb rule
---
# vestack_clb_rule
Provides a resource to manage clb rule
## Example Usage
```hcl
data "vestack_zones" "foo" {
}
resource "vestack_vpc" "foo" {
  vpc_name     = "acc-test-vpc"
  cidr_block   = "192.168.0.0/16"
  project_name = "default"
}

resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block  = "192.168.0.0/24"
  zone_id     = data.vestack_zones.foo.zones[0].id
  vpc_id      = vestack_vpc.foo.id
}

resource "vestack_clb" "foo" {
  type               = "public"
  subnet_id          = vestack_subnet.foo.id
  load_balancer_spec = "small_1"
  description        = "acc0Demo"
  load_balancer_name = "acc-test-create"
  eip_billing_config {
    isp              = "ChinaUnicom"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth        = 1
  }
}

resource "vestack_server_group" "foo" {
  load_balancer_id  = vestack_clb.foo.id
  server_group_name = "acc-test-create"
  description       = "hello demo11"
}

resource "vestack_listener" "foo" {
  load_balancer_id = vestack_clb.foo.id
  listener_name    = "acc-test-listener"
  protocol         = "HTTP"
  port             = 90
  server_group_id  = vestack_server_group.foo.id
  health_check {
    enabled              = "on"
    interval             = 10
    timeout              = 3
    healthy_threshold    = 5
    un_healthy_threshold = 2
    domain               = "volcengine.com"
    http_code            = "http_2xx"
    method               = "GET"
    uri                  = "/"
  }
  enabled = "on"
}
resource "vestack_clb_rule" "foo" {
  listener_id     = vestack_listener.foo.id
  server_group_id = vestack_server_group.foo.id
  domain          = "test-volc123.com"
  url             = "/tftest"
}
```
## Argument Reference
The following arguments are supported:
* `listener_id` - (Required, ForceNew) The ID of listener.
* `server_group_id` - (Required) Server Group Id.
* `description` - (Optional) The description of the Rule.
* `domain` - (Optional, ForceNew) The domain of Rule.
* `url` - (Optional, ForceNew) The Url of Rule.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
Rule can be imported using the id, e.g.
Notice: resourceId is ruleId, due to the lack of describeRuleAttributes in openapi, for import resources, please use ruleId:listenerId to import.
we will fix this problem later.
```
$ terraform import vestack_clb_rule.foo rule-273zb9hzi1gqo7fap8u1k3utb:lsn-273ywvnmiu70g7fap8u2xzg9d
```

