package rule_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/rule"
	"testing"
)

const testAccVestackClbRuleCreateConfig = `
data "vestack_zones" "foo"{
}
resource "vestack_vpc" "foo" {
	  vpc_name   = "acc-test-vpc"
	  cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
	  subnet_name = "acc-test-subnet"
	  cidr_block = "172.16.0.0/24"
	  zone_id = "${data.vestack_zones.foo.zones[0].id}"
	  vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_clb" "foo" {
  type = "public"
  subnet_id = "${vestack_subnet.foo.id}"
  load_balancer_spec = "small_1"
  description = "acc0Demo"
  load_balancer_name = "acc-test-create"
  eip_billing_config {
    isp = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth = 1
  }
}

resource "vestack_server_group" "foo" {
  load_balancer_id = "${vestack_clb.foo.id}"
  server_group_name = "acc-test-create"
  description = "hello demo11"
}

resource "vestack_listener" "foo" {
  load_balancer_id = "${vestack_clb.foo.id}"
  listener_name = "acc-test-listener"
  protocol = "HTTP"
  port = 90
  server_group_id = "${vestack_server_group.foo.id}"
  health_check {
    enabled = "on"
    interval = 10
    timeout = 3
    healthy_threshold = 5
    un_healthy_threshold = 2
    domain = "vestack.com"
    http_code = "http_2xx"
    method = "GET"
    uri = "/"
  }
  enabled = "on"
}
resource "vestack_clb_rule" "foo" {
  listener_id = "${vestack_listener.foo.id}"
  server_group_id = "${vestack_server_group.foo.id}"
  domain = "test-volc123.com"
  url = "/yyyy"
}
`

const testAccVestackClbRuleUpdateConfig = `
data "vestack_zones" "foo"{
}
resource "vestack_vpc" "foo" {
	  vpc_name   = "acc-test-vpc"
	  cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
	  subnet_name = "acc-test-subnet"
	  cidr_block = "172.16.0.0/24"
	  zone_id = "${data.vestack_zones.foo.zones[0].id}"
	  vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_clb" "foo" {
  type = "public"
  subnet_id = "${vestack_subnet.foo.id}"
  load_balancer_spec = "small_1"
  description = "acc0Demo"
  load_balancer_name = "acc-test-create"
  eip_billing_config {
    isp = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth = 1
  }
}

resource "vestack_server_group" "foo" {
  load_balancer_id = "${vestack_clb.foo.id}"
  server_group_name = "acc-test-create"
  description = "hello demo11"
}

resource "vestack_listener" "foo" {
  load_balancer_id = "${vestack_clb.foo.id}"
  listener_name = "acc-test-listener"
  protocol = "HTTP"
  port = 90
  server_group_id = "${vestack_server_group.foo.id}"
  health_check {
    enabled = "on"
    interval = 10
    timeout = 3
    healthy_threshold = 5
    un_healthy_threshold = 2
    domain = "vestack.com"
    http_code = "http_2xx"
    method = "GET"
    uri = "/"
  }
  enabled = "on"
}
resource "vestack_clb_rule" "foo" {
  listener_id = "${vestack_listener.foo.id}"
  server_group_id = "${vestack_server_group.foo.id}"
  domain = "acc-test-volc123.com"
  url = "/accyyyy"
}
`

func TestAccVestackClbRuleResource_Basic(t *testing.T) {
	resourceName := "vestack_clb_rule.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &rule.VestackRuleService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackClbRuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "domain", "test-volc123.com"),
					resource.TestCheckResourceAttr(acc.ResourceId, "url", "/yyyy"),
				),
			},
		},
	})
}

func TestAccVestackClbRuleResource_Update(t *testing.T) {
	resourceName := "vestack_clb_rule.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &rule.VestackRuleService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackClbRuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "domain", "test-volc123.com"),
					resource.TestCheckResourceAttr(acc.ResourceId, "url", "/yyyy"),
				),
			},
			{
				Config: testAccVestackClbRuleUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "domain", "acc-test-volc123.com"),
					resource.TestCheckResourceAttr(acc.ResourceId, "url", "/accyyyy"),
				),
			},
			{
				Config:             testAccVestackClbRuleUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
