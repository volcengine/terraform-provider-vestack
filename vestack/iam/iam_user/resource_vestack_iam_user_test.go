package iam_user_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_user"
)

const testAccVestackIamUserCreateConfig = `
resource "vestack_iam_user" "foo" {
  user_name = "acc-test-user"
  description = "acc test"
  display_name = "name"
  tags {
    key = "k1"
    value = "v1"
  }
}
`

const testAccVestackIamUserUpdateConfig = `
resource "vestack_iam_user" "foo" {
    description = "acc test update"
    display_name = "name2"
    email = "xxx@163.com"
    user_name = "acc-test-user2"
	tags {
    key = "k2"
    value = "v2"
  }
}
`

func TestAccVestackIamUserResource_Basic(t *testing.T) {
	resourceName := "vestack_iam_user.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_user.VestackIamUserService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamUserCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "display_name", "name"),
					resource.TestCheckResourceAttr(acc.ResourceId, "email", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "mobile_phone", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.0.key", "k1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.0.value", "v1"),
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

func TestAccVestackIamUserResource_Update(t *testing.T) {
	resourceName := "vestack_iam_user.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_user.VestackIamUserService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamUserCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "display_name", "name"),
					resource.TestCheckResourceAttr(acc.ResourceId, "email", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "mobile_phone", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.0.key", "k1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.0.value", "v1"),
				),
			},
			{
				Config: testAccVestackIamUserUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc test update"),
					resource.TestCheckResourceAttr(acc.ResourceId, "display_name", "name2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "email", "xxx@163.com"),
					resource.TestCheckResourceAttr(acc.ResourceId, "mobile_phone", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.0.key", "k2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.0.value", "v2"),
				),
			},
			{
				Config:             testAccVestackIamUserUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
