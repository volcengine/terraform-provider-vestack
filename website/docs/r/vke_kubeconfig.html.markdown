---
subcategory: "VKE"
layout: "vestack"
page_title: "Vestack: vestack_vke_kubeconfig"
sidebar_current: "docs-vestack-resource-vke_kubeconfig"
description: |-
  Provides a resource to manage vke kubeconfig
---
# vestack_vke_kubeconfig
Provides a resource to manage vke kubeconfig
## Example Usage
```hcl
data "vestack_zones" "foo" {
}

resource "vestack_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block  = "172.16.0.0/24"
  zone_id     = data.vestack_zones.foo.zones[0].id
  vpc_id      = vestack_vpc.foo.id
}

resource "vestack_security_group" "foo" {
  security_group_name = "acc-test-security-group"
  vpc_id              = vestack_vpc.foo.id
}

resource "vestack_vke_cluster" "foo" {
  name                      = "acc-test-cluster"
  description               = "created by terraform"
  delete_protection_enabled = false
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
}

resource "vestack_vke_kubeconfig" "foo" {
  cluster_id     = vestack_vke_cluster.foo.id
  type           = "Private"
  valid_duration = 2
}
```
## Argument Reference
The following arguments are supported:
* `cluster_id` - (Required, ForceNew) The cluster id of the Kubeconfig.
* `type` - (Required, ForceNew) The type of the Kubeconfig, the value of type should be Public or Private.
* `valid_duration` - (Optional, ForceNew) The ValidDuration of the Kubeconfig, the range of the ValidDuration is 1 hour to 43800 hour.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
VkeKubeconfig can be imported using the id, e.g.
```
$ terraform import vestack_vke_kubeconfig.default kce8simvqtofl0l6u4qd0
```

