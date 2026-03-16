package listener_health_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/listener_health"
)

const testAccVestackListenerHealthsDatasourceConfig = `
data "vestack_listener_healths" "foo" {
  listener_id = "lsn-13f1749gp981s3n6nu5c4209a"
}
`

func TestAccVestackListenerHealthsDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_listener_healths.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &listener_health.VestackListenerHealthService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackListenerHealthsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(acc.ResourceId, "health_info.#"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "total_count"),
				),
			},
		},
	})
}
