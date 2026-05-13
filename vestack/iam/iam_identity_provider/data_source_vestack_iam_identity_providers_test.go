package iam_identity_provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_identity_provider"
)

func TestAccVestackIamIdentityProvidersDataSource_Basic(t *testing.T) {
	resourceName := "data.vestack_iam_identity_providers.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_identity_provider.VestackIamIdentityProviderService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackIamIdentityProvidersDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(acc.ResourceId, "total_count"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "providers.#"),
				),
			},
		},
	})
}

func testAccVestackIamIdentityProvidersDataSourceConfig() string {
	return `
data "vestack_iam_identity_providers" "foo" {
}
`
}
