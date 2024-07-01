package ecs_instance_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/volcengine/terraform-provider-vestack/vestack"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_instance"
	"testing"
)

const testAccVestackEcsInstanceCreateConfig = `
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

data "vestack_images" "foo" {
  	os_type = "Linux"
  	visibility = "public"
  	instance_type_id = "ecs.g1.large"
}

resource "vestack_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
	description = "acc-test"
	host_name = "tf-acc-test"
  	image_id = "${data.vestack_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1.large"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	data_volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
	subnet_id = "${vestack_subnet.foo.id}"
	security_group_ids = ["${vestack_security_group.foo.id}"]
	project_name = "default"
	tags {
    	key = "k1"
    	value = "v1"
  	}
}
`

func TestAccVestackEcsInstanceResource_Basic(t *testing.T) {
	resourceName := "vestack_ecs_instance.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance.VestackEcsService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEcsInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "primary_ip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "hpc_cluster_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "security_enhancement_strategy"},
			},
		},
	})
}

const testAccVestackEcsInstanceUpdateBasicAttributeConfig = `
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

data "vestack_images" "foo" {
  	os_type = "Linux"
  	visibility = "public"
  	instance_type_id = "ecs.g1.large"
}

resource "vestack_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs-new"
	description = "acc-test-new"
	host_name = "tf-acc-test"
	user_data = "ZWNobyBoZWxsbyBlY3Mh"
  	image_id = "${data.vestack_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1.large"
  	password = "93f0cb0614Aab12new"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	data_volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
	subnet_id = "${vestack_subnet.foo.id}"
	security_group_ids = ["${vestack_security_group.foo.id}"]
	project_name = "default"
	tags {
    	key = "k1"
    	value = "v1"
  	}
}
`

func TestAccVestackEcsInstanceResource_Update_BasicAttribute(t *testing.T) {
	resourceName := "vestack_ecs_instance.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance.VestackEcsService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEcsInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "primary_ip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "hpc_cluster_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config: testAccVestackEcsInstanceUpdateBasicAttributeConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", "echo hello ecs!"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "primary_ip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "hpc_cluster_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config:             testAccVestackEcsInstanceUpdateBasicAttributeConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccVestackEcsInstanceUpdateSecurityGroupConfig = `
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
  	security_group_name = "acc-test-security-group-${count.index}"
  	vpc_id = "${vestack_vpc.foo.id}"
	count = 3
}

data "vestack_images" "foo" {
  	os_type = "Linux"
  	visibility = "public"
  	instance_type_id = "ecs.g1.large"
}

resource "vestack_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
	description = "acc-test"
	host_name = "tf-acc-test"
  	image_id = "${data.vestack_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1.large"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	data_volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
	subnet_id = "${vestack_subnet.foo.id}"
	security_group_ids = vestack_security_group.foo[*].id
	project_name = "default"
	tags {
    	key = "k1"
    	value = "v1"
  	}
}
`

func TestAccVestackEcsInstanceResource_Update_SecurityGroup(t *testing.T) {
	resourceName := "vestack_ecs_instance.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance.VestackEcsService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEcsInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "primary_ip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "hpc_cluster_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config: testAccVestackEcsInstanceUpdateSecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "primary_ip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "hpc_cluster_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config:             testAccVestackEcsInstanceUpdateSecurityGroupConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccVestackEcsInstanceUpdateSystemVolumeConfig = `
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

data "vestack_images" "foo" {
  	os_type = "Linux"
  	visibility = "public"
  	instance_type_id = "ecs.g1.large"
}

resource "vestack_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
	description = "acc-test"
	host_name = "tf-acc-test"
  	image_id = "${data.vestack_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1.large"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 50
	data_volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
	subnet_id = "${vestack_subnet.foo.id}"
	security_group_ids = ["${vestack_security_group.foo.id}"]
	project_name = "default"
	tags {
    	key = "k1"
    	value = "v1"
  	}
}
`

func TestAccVestackEcsInstanceResource_Update_SystemVolume(t *testing.T) {
	resourceName := "vestack_ecs_instance.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance.VestackEcsService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEcsInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "primary_ip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "hpc_cluster_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config: testAccVestackEcsInstanceUpdateSystemVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "primary_ip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "hpc_cluster_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config:             testAccVestackEcsInstanceUpdateSystemVolumeConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccVestackEcsInstanceUpdateInstanceTypeConfig = `
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

data "vestack_images" "foo" {
  	os_type = "Linux"
  	visibility = "public"
  	instance_type_id = "ecs.g1.xlarge"
}

resource "vestack_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
	description = "acc-test"
	host_name = "tf-acc-test"
  	image_id = "${data.vestack_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1.xlarge"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	data_volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
	subnet_id = "${vestack_subnet.foo.id}"
	security_group_ids = ["${vestack_security_group.foo.id}"]
	project_name = "default"
	tags {
    	key = "k1"
    	value = "v1"
  	}
}
`

func TestAccVestackEcsInstanceResource_Update_InstanceType(t *testing.T) {
	resourceName := "vestack_ecs_instance.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance.VestackEcsService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEcsInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "primary_ip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "hpc_cluster_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config: testAccVestackEcsInstanceUpdateInstanceTypeConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "primary_ip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "hpc_cluster_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config:             testAccVestackEcsInstanceUpdateInstanceTypeConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccVestackEcsInstanceUpdateImageConfig = `
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

data "vestack_images" "foo" {
  	os_type = "Linux"
  	visibility = "public"
  	instance_type_id = "ecs.g1.large"
}

resource "vestack_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
	description = "acc-test"
	host_name = "tf-acc-test"
  	image_id = "${data.vestack_images.foo.images[1].image_id}"
  	instance_type = "ecs.g1.large"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	data_volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
	subnet_id = "${vestack_subnet.foo.id}"
	security_group_ids = ["${vestack_security_group.foo.id}"]
	project_name = "default"
	tags {
    	key = "k1"
    	value = "v1"
  	}
}
`

func TestAccVestackEcsInstanceResource_Update_Image(t *testing.T) {
	resourceName := "vestack_ecs_instance.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance.VestackEcsService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEcsInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "primary_ip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "hpc_cluster_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config: testAccVestackEcsInstanceUpdateImageConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "primary_ip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "hpc_cluster_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config:             testAccVestackEcsInstanceUpdateImageConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccVestackEcsInstanceUpdateTagsConfig = `
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

data "vestack_images" "foo" {
  	os_type = "Linux"
  	visibility = "public"
  	instance_type_id = "ecs.g1.large"
}

resource "vestack_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
	description = "acc-test"
	host_name = "tf-acc-test"
  	image_id = "${data.vestack_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1.large"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	data_volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
	subnet_id = "${vestack_subnet.foo.id}"
	security_group_ids = ["${vestack_security_group.foo.id}"]
	project_name = "default"
	tags {
    	key = "k2"
    	value = "v2"
  	}
	tags {
    	key = "k3"
    	value = "v3"
  	}
}
`

func TestAccVestackEcsInstanceResource_Update_Tags(t *testing.T) {
	resourceName := "vestack_ecs_instance.foo"

	acc := &vestack.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance.VestackEcsService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			vestack.AccTestPreCheck(t)
		},
		Providers:    vestack.GetTestAccProviders(),
		CheckDestroy: vestack.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVestackEcsInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "primary_ip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "hpc_cluster_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config: testAccVestackEcsInstanceUpdateTagsConfig,
				Check: resource.ComposeTestCheckFunc(
					vestack.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "2"),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k2",
						"value": "v2",
					}),
					vestack.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k3",
						"value": "v3",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "primary_ip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "hpc_cluster_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config:             testAccVestackEcsInstanceUpdateTagsConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
