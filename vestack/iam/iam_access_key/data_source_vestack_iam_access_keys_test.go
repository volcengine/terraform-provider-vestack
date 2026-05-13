package iam_access_key_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_access_key"
)

const testAccVestackIamAccessKeysDatasourceConfig = `
data "vestack_iam_access_keys" "foo"{
  user_name = "inner-user"
}
`

func TestAccVestackIamAccessKeysDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_iam_access_keys.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_access_key.VestackIamAccessKeyService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamAccessKeysDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "access_key_metadata.#", "1"),
				),
			},
		},
	})
}
