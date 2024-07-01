package iam_login_profile_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_login_profile"
)

const testAccVestackIamLoginProfileCreateConfig = `
resource "vestack_iam_user" "foo" {
  	user_name = "acc-test-user"
  	description = "acc-test"
  	display_name = "name"
}

resource "vestack_iam_login_profile" "foo" {
    user_name = "${vestack_iam_user.foo.user_name}"
  	password = "93f0cb0614Aab12"
  	login_allowed = true
	password_reset_required = false
}
`

func TestAccVestackIamLoginProfileResource_Basic(t *testing.T) {
	resourceName := "vestack_iam_login_profile.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_login_profile.VestackIamLoginProfileService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamLoginProfileCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "login_allowed", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "password_reset_required", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

const testAccVestackIamLoginProfileUpdateConfig = `
resource "vestack_iam_user" "foo" {
  	user_name = "acc-test-user"
  	description = "acc-test"
  	display_name = "name"
}

resource "vestack_iam_login_profile" "foo" {
    user_name = "${vestack_iam_user.foo.user_name}"
  	password = "93f0cb0614Aab12177"
  	login_allowed = false
	password_reset_required = true
}
`

func TestAccVestackIamLoginProfileResource_Update(t *testing.T) {
	resourceName := "vestack_iam_login_profile.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_login_profile.VestackIamLoginProfileService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamLoginProfileCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "login_allowed", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "password_reset_required", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
				),
			},
			{
				Config: testAccVestackIamLoginProfileUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "login_allowed", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12177"),
					resource.TestCheckResourceAttr(acc.ResourceId, "password_reset_required", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
				),
			},
			{
				Config:             testAccVestackIamLoginProfileUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
