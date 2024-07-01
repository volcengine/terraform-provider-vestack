package support_addon_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/support_addon"
	"testing"
)

const testAccVestackVkeSupportAddonsDatasourceConfig = `
data "vestack_vke_support_addons" "foo"{
}
`

func TestAccVestackVkeSupportAddonsDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_vke_support_addons.foo"

	_ = &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &support_addon.VestackVkeSupportAddonService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackVkeSupportAddonsDatasourceConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}
