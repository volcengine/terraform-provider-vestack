package ecs_key_pair_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_key_pair"
	"testing"
)

const testAccVestackEcsKeyPairsDatasourceConfig = `
resource "vestack_ecs_key_pair" "foo" {
  key_pair_name = "acc-test-key-name"
  description ="acc-test"
}
data "vestack_ecs_key_pairs" "foo"{
    key_pair_name = "${vestack_ecs_key_pair.foo.key_pair_name}"
}
`

func TestAccVestackEcsKeyPairsDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_ecs_key_pairs.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_key_pair.VestackEcsKeyPairService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEcsKeyPairsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pairs.#", "1"),
				),
			},
		},
	})
}
