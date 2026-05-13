package acl_entry_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/acl_entry"
	"testing"
)

const testAccVestackAclEntryCreateConfig = `
resource "vestack_acl" "foo" {
	acl_name = "acc-test-acl"
	description = "acc-test-demo"
	project_name = "default"
}

resource "vestack_acl_entry" "foo" {
    acl_id = "${vestack_acl.foo.id}"
    entry = "172.20.1.0/24"
	description = "entry"
}
`

func TestAccVestackAclEntryResource_Basic(t *testing.T) {
	resourceName := "vestack_acl_entry.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &acl_entry.VestackAclEntryService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackAclEntryCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "entry", "172.20.1.0/24"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "entry"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "acl_id"),
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
