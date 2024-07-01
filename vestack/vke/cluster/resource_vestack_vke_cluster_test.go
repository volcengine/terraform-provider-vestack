package cluster_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/cluster"
	"testing"
)

const testAccVestackVkeClusterCreateConfig = `
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

`

const testAccVestackVkeClusterUpdateConfig = `
resource "vestack_vpc" "foo" {
    vpc_name = "acc-test-project1"
    cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
    subnet_name = "acc-subnet-test-2"
    cidr_block = "172.16.0.0/24"
    zone_id = "cn-beijing-a"
    vpc_id = vestack_vpc.foo.id
}

resource "vestack_security_group" "foo" {
    vpc_id = vestack_vpc.foo.id
    security_group_name = "acc-test-security-group2"
}

resource "vestack_ecs_instance" "foo" {
    image_id = "image-ybqi99s7yq8rx7mnk44b"
    instance_type = "ecs.g1ie.large"
    instance_name = "acc-test-ecs-name2"
    password = "93f0cb0614Aab12"
    instance_charge_type = "PostPaid"
    system_volume_type = "ESSD_PL0"
    system_volume_size = 40
    subnet_id = vestack_subnet.foo.id
    security_group_ids = [vestack_security_group.foo.id]
    lifecycle {
        ignore_changes = [security_group_ids, instance_name]
    }
}

resource "vestack_vke_cluster" "foo" {
    name = "acc-test-2"
    description = "created by terraform update"
    delete_protection_enabled = false
    cluster_config {
        subnet_ids = [vestack_subnet.foo.id]
        api_server_public_access_enabled = false
        api_server_public_access_config {
            public_access_network_config {
                billing_type = "PostPaidByBandwidth"
                bandwidth = 2
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
    tags {
        key = "tf-k2"
        value = "tf-v2"
    }
}

`

func TestAccVestackVkeClusterResource_Basic(t *testing.T) {
	resourceName := "vestack_vke_cluster.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &cluster.VestackVkeClusterService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackVkeClusterCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.0.api_server_public_access_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_protection_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "created by terraform"),
					resource.TestCheckResourceAttr(acc.ResourceId, "logging_config.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "pods_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "pods_config.0.pod_network_mode", "VpcCniShared"),
					resource.TestCheckResourceAttr(acc.ResourceId, "services_config.#", "1"),
					vestack.TestCheckTypeSetElemAttr(acc.ResourceId, "services_config.0.service_cidrsv4.*", "172.30.0.0/18"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k1",
						"value": "tf-v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.0.api_server_public_access_config.0.public_access_network_config.0.bandwidth", "1"),
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

func TestAccVestackVkeClusterResource_Update(t *testing.T) {
	resourceName := "vestack_vke_cluster.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &cluster.VestackVkeClusterService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackVkeClusterCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.0.api_server_public_access_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_protection_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "created by terraform"),
					resource.TestCheckResourceAttr(acc.ResourceId, "logging_config.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "pods_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "pods_config.0.pod_network_mode", "VpcCniShared"),
					resource.TestCheckResourceAttr(acc.ResourceId, "services_config.#", "1"),
					vestack.TestCheckTypeSetElemAttr(acc.ResourceId, "services_config.0.service_cidrsv4.*", "172.30.0.0/18"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k1",
						"value": "tf-v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.0.api_server_public_access_config.0.public_access_network_config.0.bandwidth", "1"),
				),
			},
			{
				Config: testAccVestackVkeClusterUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.0.api_server_public_access_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_protection_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(acc.ResourceId, "logging_config.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "pods_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "pods_config.0.pod_network_mode", "VpcCniShared"),
					resource.TestCheckResourceAttr(acc.ResourceId, "services_config.#", "1"),
					vestack.TestCheckTypeSetElemAttr(acc.ResourceId, "services_config.0.service_cidrsv4.*", "172.30.0.0/18"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "2"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k1",
						"value": "tf-v1",
					}),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k2",
						"value": "tf-v2",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.0.api_server_public_access_config.0.public_access_network_config.0.bandwidth", "0"),
				),
			},
			{
				Config:             testAccVestackVkeClusterUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
