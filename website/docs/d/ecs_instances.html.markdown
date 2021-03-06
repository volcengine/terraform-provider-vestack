---
subcategory: "ECS"
layout: "vestack"
page_title: "Vestack: vestack_ecs_instances"
sidebar_current: "docs-vestack-datasource-ecs_instances"
description: |-
  Use this data source to query detailed information of ecs instances
---
# vestack_ecs_instances
Use this data source to query detailed information of ecs instances
## Example Usage
```hcl
data "vestack_ecs_instances" "foo" {
  ids = ["i-ebgy6xmgjve0384ncgsc"]
}
```
## Argument Reference
The following arguments are supported:
* `hpc_cluster_id` - (Optional) The hpc cluster ID of ECS instance.
* `ids` - (Optional) A list of ECS instance IDs.
* `instance_charge_type` - (Optional) The charge type of ECS instance.
* `key_pair_name` - (Optional) The key pair name of ECS instance.
* `name_regex` - (Optional) A Name Regex of ECS instance.
* `output_file` - (Optional) File name where to save data source results.
* `primary_ip_address` - (Optional) The primary ip address of ECS instance.
* `status` - (Optional) The status of ECS instance.
* `vpc_id` - (Optional) The VPC ID of ECS instance.
* `zone_id` - (Optional) The available zone ID of ECS instance.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `instances` - The collection of ECS instance query.
  * `cpus` - The number of ECS instance CPU cores.
  * `created_at` - The create time of ECS instance.
  * `description` - The description of ECS instance.
  * `host_name` - The host name of ECS instance.
  * `image_id` - The image ID of ECS instance.
  * `instance_charge_type` - The charge type of ECS instance.
  * `instance_id` - The ID of ECS instance.
  * `instance_name` - The name of ECS instance.
  * `instance_type` - The spec type of ECS instance.
  * `key_pair_id` - The ssh key ID of ECS instance.
  * `key_pair_name` - The ssh key name of ECS instance.
  * `memory_size` - The memory size of ECS instance.
  * `network_interfaces` - The networkInterface detail collection of ECS instance.
    * `mac_address` - The mac address of networkInterface.
    * `network_interface_id` - The ID of networkInterface.
    * `primary_ip_address` - The private ip address of networkInterface.
    * `subnet_id` - The subnet ID of networkInterface.
    * `type` - The type of networkInterface.
    * `vpc_id` - The ID of networkInterface.
  * `os_name` - The os name of ECS instance.
  * `os_type` - The os type of ECS instance.
  * `status` - The status of ECS instance.
  * `stopped_mode` - The stop mode of ECS instance.
  * `updated_at` - The update time of ECS instance.
  * `volumes` - The volume detail collection of volume.
    * `delete_with_instance` - The delete with instance flag of volume.
    * `size` - The size of volume.
    * `volume_id` - The ID of volume.
    * `volume_name` - The Name of volume.
    * `volume_type` - The type of volume.
  * `vpc_id` - The VPC ID of ECS instance.
  * `zone_id` - The available zone ID of ECS instance.
* `total_count` - The total count of ECS instance query.


