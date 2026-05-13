package health_check_log_topic_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/health_check_log_topic"
)

const testAccVestackHealthCheckLogTopicCreateConfig = `
resource "vestack_health_check_log_topic" "foo" {
  log_topic_id     = "05df22e2-f561-4081-8cf3-b201af564407"
  load_balancer_id = "clb-mim12q0soe805smt1bebim25"
}
`

func TestAccVestackHealthCheckLogTopicResource_Basic(t *testing.T) {
	resourceName := "vestack_health_check_log_topic.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &health_check_log_topic.VestackHealthCheckLogTopicService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackHealthCheckLogTopicCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "log_topic_id", "05df22e2-f561-4081-8cf3-b201af564407"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_id", "clb-mim12q0soe805smt1bebim25"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
