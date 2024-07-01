package default_node_pool_batch_attach_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/default_node_pool_batch_attach"
	"testing"
)

const testAccVestackVkeDefaultNodePoolBatchAttachCreateConfig = `
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

resource "vestack_ecs_instance" "foo" {
	//默认节点池中的节点只能使用指定镜像，请参考 https://www.vestack.com/docs/6460/115194#%E9%95%9C%E5%83%8F-id-%E9%85%8D%E7%BD%AE%E8%AF%B4%E6%98%8E
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

resource "vestack_vke_default_node_pool_batch_attach" "foo" {
    cluster_id = vestack_vke_cluster.foo.id
    default_node_pool_id = vestack_vke_default_node_pool.foo.id
    instances {
        instance_id = vestack_ecs_instance.foo.id
        keep_instance_name = true
        additional_container_storage_enabled = false
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
}
`

const testAccVestackVkeDefaultNodePoolBatchAttachUpdateConfig = `
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

resource "vestack_ecs_instance" "foo" {
	//默认节点池中的节点只能使用指定镜像，请参考 https://www.vestack.com/docs/6460/115194#%E9%95%9C%E5%83%8F-id-%E9%85%8D%E7%BD%AE%E8%AF%B4%E6%98%8E
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

resource "vestack_ecs_instance" "foo2" {
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

resource "vestack_vke_default_node_pool_batch_attach" "foo" {
    cluster_id = vestack_vke_cluster.foo.id
    default_node_pool_id = vestack_vke_default_node_pool.foo.id
    instances {
        instance_id = vestack_ecs_instance.foo.id
        keep_instance_name = true
        additional_container_storage_enabled = false
    }
	instances {
        instance_id = vestack_ecs_instance.foo2.id
        keep_instance_name = true
        additional_container_storage_enabled = false
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
}
`

func TestAccVestackVkeDefaultNodePoolBatchAttachResource_Basic(t *testing.T) {
	resourceName := "vestack_vke_default_node_pool_batch_attach.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return default_node_pool_batch_attach.NewVestackVkeDefaultNodePoolBatchAttachService(client)
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
				Config: testAccVestackVkeDefaultNodePoolBatchAttachCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
				),
			},
		},
	})
}

func TestAccVestackVkeDefaultNodePoolBatchAttachResource_Update(t *testing.T) {
	resourceName := "vestack_vke_default_node_pool_batch_attach.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return default_node_pool_batch_attach.NewVestackVkeDefaultNodePoolBatchAttachService(client)
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
				Config: testAccVestackVkeDefaultNodePoolBatchAttachCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
				),
			},
			{
				Config: testAccVestackVkeDefaultNodePoolBatchAttachUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
				),
			},
			{
				Config:             testAccVestackVkeDefaultNodePoolBatchAttachUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
