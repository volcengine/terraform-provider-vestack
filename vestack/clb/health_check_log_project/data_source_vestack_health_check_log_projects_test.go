package health_check_log_project_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/health_check_log_project"
)

const testAccVestackHealthCheckLogProjectsDatasourceConfig = `
data "vestack_health_check_log_projects" "foo" {
}
`

func TestAccVestackHealthCheckLogProjectsDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_health_check_log_projects.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &health_check_log_project.VestackHealthCheckLogProjectService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackHealthCheckLogProjectsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(acc.ResourceId, "health_check_log_projects.#"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "total_count"),
				),
			},
		},
	})
}
