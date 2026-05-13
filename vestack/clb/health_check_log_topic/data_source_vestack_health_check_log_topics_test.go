package health_check_log_topic_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/health_check_log_topic"
)

const testAccVestackHealthCheckLogTopicsDatasourceConfig = `
data "vestack_health_check_log_topics" "foo"{
    log_topic_id = "05df22e2-f561-4081-8cf3-b201af564407"
}
`

func TestAccVestackHealthCheckLogTopicsDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_health_check_log_topics.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &health_check_log_topic.VestackHealthCheckLogTopicService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackHealthCheckLogTopicsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(acc.ResourceId, "health_check_log_topics.#"),
				),
			},
		},
	})
}
