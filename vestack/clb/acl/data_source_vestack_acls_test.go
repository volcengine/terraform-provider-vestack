package acl_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/acl"
	"testing"
)

const testAccVestackAclsDatasourceConfig = `
resource "vestack_acl" "foo" {
	acl_name = "acc-test-acl-${count.index}"
	description = "acc-test-demo"
	project_name = "default"
	acl_entries {
    	entry = "172.20.1.0/24"
    	description = "e1"
  	}
	count = 3
}

data "vestack_acls" "foo"{
    ids = vestack_acl.foo[*].id
}
`

func TestAccVestackAclsDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_acls.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &acl.VestackAclService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackAclsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "acls.#", "3"),
				),
			},
		},
	})
}
