package eip_address_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/eip/eip_address"
	"testing"
)

const testAccVestackEipAddressesDatasourceConfig = `
resource "vestack_eip_address" "foo" {
    billing_type = "PostPaidByTraffic"
}
data "vestack_eip_addresses" "foo"{
    ids = ["${vestack_eip_address.foo.id}"]
}
`

func TestAccVestackEipAddressesDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_eip_addresses.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &eip_address.VestackEipAddressService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEipAddressesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "addresses.#", "1"),
				),
			},
		},
	})
}
