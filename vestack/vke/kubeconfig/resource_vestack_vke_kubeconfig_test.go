package kubeconfig_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/kubeconfig"
)

const testAccVestackVkeKubeconfigCreateConfig = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
	vpc_name   = "acc-test-vpc"
  	cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
  	subnet_name = "acc-test-subnet"
  	cidr_block = "172.16.0.0/24"
  	zone_id = "${data.vestack_zones.foo.zones[0].id}"
	vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_security_group" "foo" {
  	security_group_name = "acc-test-security-group"
  	vpc_id = "${vestack_vpc.foo.id}"
}

resource "vestack_vke_cluster" "foo" {
    name = "acc-test-cluster"
    description = "created by terraform"
    delete_protection_enabled = false
    cluster_config {
        subnet_ids = ["${vestack_subnet.foo.id}"]
        api_server_public_access_enabled = true
        api_server_public_access_config {
            public_access_network_config {
                billing_type = "PostPaidByBandwidth"
                bandwidth = 1
            }
        }
        resource_public_access_default_enabled = true
    }
    pods_config {
        pod_network_mode = "VpcCniShared"
        vpc_cni_config {
            subnet_ids = ["${vestack_subnet.foo.id}"]
        }
    }
    services_config {
        service_cidrsv4 = ["172.30.0.0/18"]
    }
    tags {
        key = "tf-k1"
        value = "tf-v1"
    }
}

resource "vestack_vke_kubeconfig" "foo" {
    cluster_id = "${vestack_vke_cluster.foo.id}"
    type = "Private"
	valid_duration = 2
}
`

func TestAccVestackVkeKubeconfigResource_Basic(t *testing.T) {
	resourceName := "vestack_vke_kubeconfig.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &kubeconfig.VestackVkeKubeconfigService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackVkeKubeconfigCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "Private"),
					resource.TestCheckResourceAttr(acc.ResourceId, "valid_duration", "2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"valid_duration"},
			},
		},
	})
}
