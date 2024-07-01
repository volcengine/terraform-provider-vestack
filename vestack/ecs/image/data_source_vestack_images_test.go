package image_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/image"
	"testing"
)

const testAccVestackImagesDatasourceConfig = `
data "vestack_images" "foo" {
	  os_type = "Linux"
	  visibility = "public"
	  instance_type_id = "ecs.g1.large"
}
`

func TestAccVestackImagesDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_images.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &image.VestackImageService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackImagesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "images.#", "26"),
				),
			},
		},
	})
}
