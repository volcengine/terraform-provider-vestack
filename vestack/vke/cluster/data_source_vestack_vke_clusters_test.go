package cluster_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/cluster"
	"testing"
)

const testAccVestackVkeClustersDatasourceConfig = `
data "vestack_zones" "foo"{
}

resource "vestack_vpc" "foo" {
    vpc_name = "acc-test-project1"
    cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
    subnet_name = "acc-subnet-test-2"
    cidr_block = "172.16.0.0/24"
    zone_id = data.vestack_zones.foo.zones[0].id
    vpc_id = vestack_vpc.foo.id
}

resource "vestack_security_group" "foo" {
    vpc_id = vestack_vpc.foo.id
    security_group_name = "acc-test-security-group2"
}

resource "vestack_vke_cluster" "foo" {
    name = "acc-test-1"
    description = "created by terraform"
    delete_protection_enabled = false
    cluster_config {
        subnet_ids = [vestack_subnet.foo.id]
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
            subnet_ids = [vestack_subnet.foo.id]
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

data "vestack_vke_clusters" "foo"{
    ids = [vestack_vke_cluster.foo.id]
}
`

func TestAccVestackVkeClustersDatasource_Basic(t *testing.T) {
	resourceName := "data.vestack_vke_clusters.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &cluster.VestackVkeClusterService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers: vestack.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackVkeClustersDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "clusters.#", "1"),
				),
			},
		},
	})
}
