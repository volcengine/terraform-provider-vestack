package iam_access_key_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_access_key"
)

const testAccVestackIamAccessKeyCreateConfig = `
resource "vestack_iam_user" "foo" {
  	user_name = "acc-test-user"
  	description = "acc-test"
  	display_name = "name"
}

resource "vestack_iam_access_key" "foo" {
	user_name = "${vestack_iam_user.foo.user_name}"
    secret_file = "./sk"
    status = "active"
}
`

func TestAccVestackIamAccessKeyResource_Basic(t *testing.T) {
	resourceName := "vestack_iam_access_key.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_access_key.VestackIamAccessKeyService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamAccessKeyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secret_file", "./sk"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "active"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_date"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "secret"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "pgp_key"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "encrypted_secret"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "key_fingerprint"),
				),
			},
		},
	})
}

const testAccVestackIamAccessKeyUpdateConfig = `
resource "vestack_iam_user" "foo" {
  	user_name = "acc-test-user"
  	description = "acc-test"
  	display_name = "name"
}

resource "vestack_iam_access_key" "foo" {
	user_name = "${vestack_iam_user.foo.user_name}"
    secret_file = "./sk"
    status = "inactive"
}
`

func TestAccVestackIamAccessKeyResource_Update(t *testing.T) {
	resourceName := "vestack_iam_access_key.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_access_key.VestackIamAccessKeyService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamAccessKeyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secret_file", "./sk"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "active"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_date"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "secret"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "pgp_key"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "encrypted_secret"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "key_fingerprint"),
				),
			},
			{
				Config: testAccVestackIamAccessKeyUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secret_file", "./sk"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "inactive"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_date"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "secret"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "pgp_key"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "encrypted_secret"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "key_fingerprint"),
				),
			},
			{
				Config:             testAccVestackIamAccessKeyUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
