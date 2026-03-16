package health_check_log_project_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/health_check_log_project"
)

const testAccVestackHealthCheckLogProjectCreateConfig = `
resource "vestack_health_check_log_project" "foo" {
}
`

func TestAccVestackHealthCheckLogProjectResource_Basic(t *testing.T) {
	resourceName := "vestack_health_check_log_project.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &health_check_log_project.VestackHealthCheckLogProjectService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackHealthCheckLogProjectCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "log_project_id"),
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
