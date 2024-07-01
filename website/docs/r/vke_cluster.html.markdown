---
subcategory: "VKE"
layout: "vestack"
page_title: "Vestack: vestack_vke_cluster"
sidebar_current: "docs-vestack-resource-vke_cluster"
description: |-
  Provides a resource to manage vke cluster
---
# vestack_vke_cluster
Provides a resource to manage vke cluster
## Example Usage
```hcl
resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-project1"
  cidr_block = "192.168.0.0/16"
}

resource "vestack_subnet" "foo" {
  subnet_name = "acc-subnet-test-2"
  cidr_block  = "192.168.1.0/24"
  zone_id     = "cn-e10-jicheng-a"
  vpc_id      = vestack_vpc.foo.id
}
resource "vestack_subnet" "controller_subnet" {
  subnet_name = "acc-subnet-controller-1"
  cidr_block  = "192.168.2.0/24"
  zone_id     = "cn-e10-jicheng-a"
  vpc_id      = vestack_vpc.foo.id
}

resource "vestack_security_group" "foo" {
  vpc_id              = vestack_vpc.foo.id
  security_group_name = "acc-test-security-group2"
}

#resource "vestack_ecs_instance" "foo" {
#  image_id = "image-ycmpf7lm129tifwxg27g"
#  instance_type = "ecs.g1i.xlarge"
#  instance_name = "acc-test-ecs-name2"
#  password = "93f0cb0614Aab12"
#  instance_charge_type = "PostPaid"
#  system_volume_type = "ESSD_PL0"
#  system_volume_size = 40
#  subnet_id = vestack_subnet.foo.id
#  security_group_ids = [vestack_security_group.foo.id]
#  lifecycle {
#    ignore_changes = [security_group_ids, instance_name]
#  }
#}

resource "vestack_vke_cluster" "foo" {
  name                      = "acc-test-1"
  type                      = "Standard"
  description               = "created by terraform"
  delete_protection_enabled = false
  control_plane_nodes_config {
    provider = "VeStack"
    ve_stack {
      new_node_configs {
        count            = 3
        subnet_ids       = [vestack_subnet.controller_subnet.id]
        instance_type_id = "ecs.g1i.xlarge"
        system_volume {
          size = 40
          type = "ESSD_PL0"
        }
        security {
          security_strategies = []
          login {
            password = "Um9vdEAxMjM="
          }
        }
      }
    }
  }
  cluster_config {
    subnet_ids                       = [vestack_subnet.foo.id]
    api_server_public_access_enabled = true
    api_server_public_access_config {
      public_access_network_config {
        billing_type = "PostPaidByBandwidth"
        bandwidth    = 1
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
    key   = "tf-k1"
    value = "tf-v1"
  }
  logging_config {
    log_setups {
      enabled  = false
      log_ttl  = 30
      log_type = "Audit"
    }
  }
}
```
## Argument Reference
The following arguments are supported:
* `cluster_config` - (Required) The config of the cluster.
* `control_plane_nodes_config` - (Required) The control plane node information for the VKE cluster instance.
* `name` - (Required) The name of the cluster.
* `pods_config` - (Required) The config of the pods.
* `services_config` - (Required, ForceNew) The config of the services.
* `client_token` - (Optional) ClientToken is a case-sensitive string of no more than 64 ASCII characters passed in by the caller.
* `delete_protection_enabled` - (Optional) The delete protection of the cluster, the value is `true` or `false`.
* `description` - (Optional) The description of the cluster.
* `kubernetes_version` - (Optional, ForceNew) The version of Kubernetes specified when creating a VKE cluster (specified to patch version), if not specified, the latest Kubernetes version supported by VKE is used by default, which is a 3-segment version format starting with a lowercase v, that is, KubernetesVersion with IsLatestVersion=True in the return value of ListSupportedVersions.
* `logging_config` - (Optional) Cluster log configuration information.
* `tags` - (Optional) Tags.
* `type` - (Optional) Type of the Cluster.

The `api_server_public_access_config` object supports the following:

* `public_access_network_config` - (Optional, ForceNew) Public network access network configuration.

The `bgp_config` object supports the following:

* `as_number` - (Required, ForceNew) Value range [64512, 65534].
* `mode` - (Required, ForceNew) BGP mode, optional values are: FullMesh | RouteReflectors: - FullMesh: all nodes serve as RR peers. - RouteReflectors: The Master node as the RR contains the routing information of all nodes, and the node only has routes pointing to the RR node.
* `route_reflector_peer_points` - (Required, ForceNew) Pod CIDR for the Flannel container network.
* `external_route_reflector_enabled` - (Optional, ForceNew) Whether to enable external RoutReflector, optional values true | false: - true: Enable external RoutReflector instead of using Master node as RoutReflector. - false: Use Master node as RoutReflector. Default value: false.

The `calico_config` object supports the following:

* `pod_cidrs` - (Required, ForceNew) Pod CIDR for the Flannel container network.
* `bgp_config` - (Optional, ForceNew) Configuration information of BGP mode under Calico. Only supported in Onpremise cluster & CalicoBgp mode.
* `max_pods_per_node` - (Optional, ForceNew) The maximum number of single-node Pod instances for a Flannel container network, the value can be `16` or `32` or `64` or `128` or `256`, default value is `64`.

The `cluster_config` object supports the following:

* `subnet_ids` - (Required, ForceNew) The subnet ID for the cluster control plane to communicate within the private network.
* `api_server_public_access_config` - (Optional) Cluster API Server public network access configuration.
* `api_server_public_access_enabled` - (Optional) Cluster API Server public network access configuration, the value is `true` or `false`.
* `resource_public_access_default_enabled` - (Optional, ForceNew) Node public network access configuration, the value is `true` or `false`.

The `control_plane_nodes_config` object supports the following:

* `provider` - (Required) Node resource provider name, available values: VeStack: Resources built on veStack full-stack version.
* `ve_stack` - (Optional) The resources in veStack are used for the master node in the VKE cluster.

The `data_volumes` object supports the following:

* `mount_point` - (Optional) The target mounting directory after disk formatting.
* `size` - (Optional, ForceNew) Disk size, unit GB, value range is 20~32768, default value is 20.
* `type` - (Optional, ForceNew) The Type of DataVolumes, the value can be `ESSD_PL0` or `ESSD_FlexPL`.

The `existed_node_config` object supports the following:

* `additional_container_storage_enabled` - (Optional) Select the data disk of the configuration node and format it and mount it as the storage directory for container images and logs. The value is:false: (default) off.true: enable.
* `container_storage_path` - (Optional) Use this data disk device to mount the container and image storage directory /var/lib/containerd. It is only valid when AdditionalContainerStorageEnabled=true and cannot be empty.The following conditions must be met, otherwise the initialization will fail:Only cloud server instances with mounted data disks are supported.When specifying the data disk device name, please ensure that the data disk device exists, and the problem will be automatically initialized.When specifying a data disk partition or logical volume name, make sure that the partition or logical volume exists and is an exct4 file system.NoticeWhen specifying a data disk device, it will be automatically formatted and mounted directly. Please be sure to back up the data in advance.When specifying a data disk partition or logical volume name, no formatting is required.
* `initialize_script` - (Optional) cript that is executed after ECS nodes are created and Kubernetes components are deployed. Supports Shell format, the length after Base64 encoding does not exceed 16 KB.
* `instances` - (Optional) ECS node information list.
* `keep_instance_name` - (Optional) Keep the node name to join the cluster, the priority is higher than NamePrefix, the value is:false: (Default) Do not maintain node names.true: Keep the node name as the original host instance name.
* `name_prefix` - (Optional) The node naming prefix has a lower priority than KeepInstanceName. When the value is empty, it means that the node naming prefix is not enabled. Among them, the prefix verification rules:Supports English letters, numbers and dashes -, dashes - cannot be used continuously.It can only start with an English letter and end with an English letter or number.The length is 2 to 51 characters.
* `security` - (Optional) Node security configuration.

The `flannel_config` object supports the following:

* `max_pods_per_node` - (Optional, ForceNew) The maximum number of single-node Pod instances for a Flannel container network, the value can be `16` or `32` or `64` or `128` or `256`.
* `pod_cidrs` - (Optional, ForceNew) Pod CIDR for the Flannel container network.

The `instances` object supports the following:


The `log_setups` object supports the following:

* `log_type` - (Required) The currently enabled log type.
* `enabled` - (Optional) Whether to enable the log option, true means enable, false means not enable, the default is false. When Enabled is changed from false to true, a new Topic will be created.
* `log_ttl` - (Optional) The storage time of logs in Log Service. After the specified log storage time is exceeded, the expired logs in this log topic will be automatically cleared. The unit is days, and the default is 30 days. The value range is 1 to 3650, specifying 3650 days means permanent storage.

The `logging_config` object supports the following:

* `log_project_id` - (Optional) The TLS log item ID of the collection target.
* `log_setups` - (Optional) Cluster logging options. This structure can only be modified and added, and cannot be deleted. When encountering a `cannot be deleted` error, please query the log setups of the current cluster and fill in the current `tf` file.

The `login` object supports the following:

* `password` - (Optional) 

The `new_node_configs` object supports the following:

* `instance_type_id` - (Required) 
* `security` - (Required) 
* `subnet_ids` - (Required, ForceNew) The subnet ID for the master node.
* `system_volume` - (Required) The SystemVolume of NodeConfig.
* `count` - (Optional) numbers of master, must be 1 3 5 7.
* `data_volumes` - (Optional, ForceNew) The DataVolumes of NodeConfig.
* `initialize_script` - (Optional) 

The `pods_config` object supports the following:

* `pod_network_mode` - (Required, ForceNew) The container network model of the cluster, the value is `Flannel` or `VpcCniShared` or `VpcCniHybrid` or `CalicoVxlan` or `CalicoBgp`. Flannel: Flannel network model, an independent Underlay container network solution, combined with the global routing capability of VPC, to achieve a high-performance network experience for the cluster. VpcCniShared: VPC-CNI network model, an Underlay container network solution based on the ENI of the private network elastic network card, with high network communication performance. CalicoVxlan: Calico network Vxlan mode, an overlay container network solution independent of the control plane. CalicoBgp: Calico network BGP mode, configure BGP between nodes or peer network infrastructure to distribute routing information (OnPremise cluster supported only).
* `calico_config` - (Optional) Calico network configuration.
* `flannel_config` - (Optional, ForceNew) Flannel network configuration.
* `vpc_cni_config` - (Optional) VPC-CNI network configuration.

The `public_access_network_config` object supports the following:

* `bandwidth` - (Optional) The peak bandwidth of the public IP, unit: Mbps.
* `billing_type` - (Optional) Billing type of public IP, the value is `PostPaidByBandwidth` or `PostPaidByTraffic`.
* `isp` - (Optional) Line type of public network IP, value:BGP: (Default) BGP circuit.ChinaMobile: China Mobile.ChinaUnicom: China Unicom.ChinaTelecom: China Telecom.

The `route_reflector_peer_points` object supports the following:

* `ip_address` - (Required, ForceNew) IP address of RouteReflector.
* `port` - (Required, ForceNew) Port of RouteReflector.
* `as_number` - (Optional, ForceNew) If ExternalRouteReflectorEnabled=true, this parameter is optional, otherwise this parameter cannot be empty [64512, 65534].

The `security` object supports the following:

* `login` - (Required) Node access mode configuration.Support password mode or key pair mode. When they are passed in at the same time, the key pair will be used first.
* `security_group_ids` - (Optional) List of security group IDs in which the node network is located.Call the DescribeSecurityGroups interface of the private network to obtain the security group ID.NoticeMust be in the same private network as the cluster.When the value is empty, the default security group of the cluster node is used by default (the naming format is <cluster ID>-common).A single node pool supports up to 5 security groups (including the default security group of cluster nodes).

The `services_config` object supports the following:

* `service_cidrsv4` - (Required, ForceNew) The IPv4 private network address exposed by the service.

The `system_volume` object supports the following:

* `size` - (Optional) Disk size, unit GB, value range is 40~2048, default value is 40.
* `type` - (Optional) The Type of SystemVolume.

The `tags` object supports the following:

* `key` - (Required) The Key of Tags.
* `value` - (Required) The Value of Tags.

The `ve_stack` object supports the following:

* `new_node_configs` - (Required) Configuration for auto create new nodes.
* `deployment_set_id` - (Optional) Deployment set ID. If specified, the master node will be added to the deployment set group. Currently, only the new node method is supported.
* `existed_node_config` - (Optional) Use an existing node as the cluster master node configuration.

The `vpc_cni_config` object supports the following:

* `subnet_ids` - (Optional) A list of Pod subnet IDs for the VPC-CNI container network.
* `vpc_id` - (Optional, ForceNew) The private network where the cluster control plane network resides.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `eip_allocation_id` - Eip allocation Id.
* `kubeconfig_private` - Kubeconfig data with private network access, returned in BASE64 encoding, it is suggested to use vke_kubeconfig instead.
* `kubeconfig_public` - Kubeconfig data with public network access, returned in BASE64 encoding, it is suggested to use vke_kubeconfig instead.


## Import
VkeCluster can be imported using the id, e.g.
```
$ terraform import vestack_vke_cluster.default cc9l74mvqtofjnoj5****
```

