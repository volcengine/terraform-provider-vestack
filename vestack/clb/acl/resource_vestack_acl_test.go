package acl_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/clb/acl"
	"testing"
)

const testAccVestackAclCreateConfig = `
resource "vestack_acl" "foo" {
	acl_name = "acc-test-acl"
	description = "acc-test-demo"
	project_name = "default"
	acl_entries {
    	entry = "172.20.1.0/24"
    	description = "e1"
  	}
}
`

func TestAccVestackAclResource_Basic(t *testing.T) {
	resourceName := "vestack_acl.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &acl.VestackAclService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackAclCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "acl_name", "acc-test-acl"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-demo"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_time"),
					resource.TestCheckResourceAttr(acc.ResourceId, "acl_entries.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "acl_entries.*", map[string]string{
						"entry":       "172.20.1.0/24",
						"description": "e1",
					}),
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

const testAccVestackAclUpdateConfig = `
resource "vestack_acl" "foo" {
    acl_name = "acc-test-acl-new"
    description = "acc-test-demo-new"
    project_name = "default"
	acl_entries {
    	entry = "172.20.2.0/24"
    	description = "e2"
  	}
	acl_entries {
    	entry = "172.20.3.0/24"
    	description = "e3"
  	}
}
`

func TestAccVestackAclResource_Update(t *testing.T) {
	resourceName := "vestack_acl.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &acl.VestackAclService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackAclCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "acl_name", "acc-test-acl"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-demo"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_time"),
					resource.TestCheckResourceAttr(acc.ResourceId, "acl_entries.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "acl_entries.*", map[string]string{
						"entry":       "172.20.1.0/24",
						"description": "e1",
					}),
				),
			},
			{
				Config: testAccVestackAclUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "acl_name", "acc-test-acl-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-demo-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_time"),
					resource.TestCheckResourceAttr(acc.ResourceId, "acl_entries.#", "2"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "acl_entries.*", map[string]string{
						"entry":       "172.20.2.0/24",
						"description": "e2",
					}),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "acl_entries.*", map[string]string{
						"entry":       "172.20.3.0/24",
						"description": "e3",
					}),
				),
			},
			{
				Config:             testAccVestackAclUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
