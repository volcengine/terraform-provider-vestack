package default_node_pool_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/default_node_pool"
	"testing"
)

const testAccVestackVkeDefaultNodePoolCreateConfig = `
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

resource "vestack_vke_default_node_pool" "foo" {
    cluster_id = vestack_vke_cluster.foo.id
    node_config {
        security {
            login {
                password = "amw4WTdVcTRJVVFsUXpVTw=="
            }
            security_group_ids = [vestack_security_group.foo.id]
            security_strategies = ["Hids"]
        }
        initialize_script = "ISMvYmluL2Jhc2gKZWNobyAx"

    }
    kubernetes_config {
        labels {
            key   = "tf-key1"
            value = "tf-value1"
        }
        labels {
            key   = "tf-key2"
            value = "tf-value2"
        }
        taints {
            key = "tf-key3"
            value = "tf-value3"
            effect = "NoSchedule"
        }
        taints {
            key = "tf-key4"
            value = "tf-value4"
            effect = "NoSchedule"
        }
        cordon = true
    }
    tags {
        key = "tf-k1"
        value = "tf-v1"
    }
}

`

const testAccVestackVkeDefaultNodePoolUpdateConfig = `
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

resource "vestack_vke_default_node_pool" "foo" {
    cluster_id = vestack_vke_cluster.foo.id
    node_config {
        security {
            login {
                password = "UHdkMTIzNDU2"
            }
            security_group_ids = [vestack_security_group.foo.id]
            security_strategies = ["Hids"]
        }
        initialize_script = "ISMvYmluL2Jhc2gKZWNobyAx"

    }
    kubernetes_config {
        labels {
            key   = "tf-key1"
            value = "tf-value1"
        }
        labels {
            key   = "tf-key2"
            value = "tf-value2"
        }
		labels {
            key   = "tf-key3"
            value = "tf-value3"
        }
        taints {
            key = "tf-key3"
            value = "tf-value3"
            effect = "NoSchedule"
        }
        taints {
            key = "tf-key4"
            value = "tf-value4"
            effect = "NoSchedule"
        }
		taints {
            key = "tf-key5"
            value = "tf-value5"
            effect = "NoSchedule"
        }
        cordon = true
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

func TestAccVestackVkeDefaultNodePoolResource_Basic(t *testing.T) {
	resourceName := "vestack_vke_default_node_pool.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return default_node_pool.NewDefaultNodePoolService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackVkeDefaultNodePoolCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "2"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "tf-key1",
						"value": "tf-value1",
					}),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "tf-key2",
						"value": "tf-value2",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "tf-key3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "tf-value3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.key", "tf-key4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.value", "tf-value4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ISMvYmluL2Jhc2gKZWNobyAx"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "amw4WTdVcTRJVVFsUXpVTw=="),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k1",
						"value": "tf-v1",
					}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_config.0.security.0.login.0.password", "is_import"},
			},
		},
	})
}

func TestAccVestackVkeDefaultNodePoolResource_Update(t *testing.T) {
	resourceName := "vestack_vke_default_node_pool.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return default_node_pool.NewDefaultNodePoolService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackVkeDefaultNodePoolCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "2"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "tf-key1",
						"value": "tf-value1",
					}),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "tf-key2",
						"value": "tf-value2",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "tf-key3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "tf-value3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.key", "tf-key4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.value", "tf-value4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ISMvYmluL2Jhc2gKZWNobyAx"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "amw4WTdVcTRJVVFsUXpVTw=="),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k1",
						"value": "tf-v1",
					}),
				),
			},
			{
				Config: testAccVestackVkeDefaultNodePoolUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "3"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "tf-key1",
						"value": "tf-value1",
					}),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "tf-key2",
						"value": "tf-value2",
					}),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "tf-key3",
						"value": "tf-value3",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "tf-key3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "tf-value3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.key", "tf-key4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.value", "tf-value4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.2.key", "tf-key5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.2.value", "tf-value5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.2.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ISMvYmluL2Jhc2gKZWNobyAx"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "2"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k1",
						"value": "tf-v1",
					}),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k2",
						"value": "tf-v2",
					}),
				),
			},
			{
				Config:             testAccVestackVkeDefaultNodePoolUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
