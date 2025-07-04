---
subcategory: "CLB"
layout: "vestack"
page_title: "Vestack: vestack_clb"
sidebar_current: "docs-vestack-resource-clb"
description: |-
  Provides a resource to manage clb
---
# vestack_clb
Provides a resource to manage clb
## Notice
When Destroy this resource,If the resource charge type is PrePaid,Please unsubscribe the resource 
in  [Vestack Console],when complete console operation,yon can
use 'terraform state rm ${resourceId}' to remove.
## Example Usage
```hcl
# query available zones in current region
data "vestack_zones" "foo" {
}

# create vpc
resource "vestack_vpc" "foo" {
  vpc_name     = "acc-test-vpc"
  cidr_block   = "192.168.0.0/16"
  project_name = "default"
}

# create subnet
resource "vestack_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block  = "192.168.0.0/24"
  zone_id     = data.vestack_zones.foo.zones[0].id
  vpc_id      = vestack_vpc.foo.id
}

# create ipv4 public clb
resource "vestack_clb" "public_clb" {
  type               = "public"
  subnet_id          = vestack_subnet.foo.id
  load_balancer_name = "acc-test-clb-public"
  load_balancer_spec = "small_1"
  description        = "acc-test-demo"
  project_name       = "default"
  eip_billing_config {
    isp              = "ChinaUnicom"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth        = 1
  }
  tags {
    key   = "k1"
    value = "v1"
  }
}

# create ipv4 private clb
resource "vestack_clb" "private_clb" {
  type               = "private"
  subnet_id          = vestack_subnet.foo.id
  load_balancer_name = "acc-test-clb-private"
  load_balancer_spec = "small_1"
  description        = "acc-test-demo"
  project_name       = "default"
}

# create eip
resource "vestack_eip_address" "eip" {
  billing_type = "PostPaidByBandwidth"
  bandwidth    = 1
  isp          = "ChinaUnicom"
  name         = "tf-eip"
  description  = "tf-test"
  project_name = "default"
}

# associate eip to clb
resource "vestack_eip_associate" "associate" {
  allocation_id = vestack_eip_address.eip.id
  instance_id   = vestack_clb.private_clb.id
  instance_type = "ClbInstance"
}

# create ipv6 vpc
resource "vestack_vpc" "vpc_ipv6" {
  vpc_name     = "acc-test-vpc-ipv6"
  cidr_block   = "192.168.0.0/16"
  enable_ipv6  = true
  project_name = "default"
  #ipv6_cidr_block = "fa00:230:0:de00::/56"
  ipv6_cidr_block_type = "ULA"

}

# create ipv6 subnet
resource "vestack_subnet" "subnet_ipv6" {
  subnet_name     = "acc-test-subnet-ipv6"
  cidr_block      = "192.168.0.0/24"
  zone_id         = data.vestack_zones.foo.zones[0].id
  vpc_id          = vestack_vpc.vpc_ipv6.id
  ipv6_cidr_block = 1
}

# create ipv6 private clb
resource "vestack_clb" "private_clb_ipv6" {
  type               = "private"
  subnet_id          = vestack_subnet.subnet_ipv6.id
  load_balancer_name = "acc-test-clb-ipv6"
  load_balancer_spec = "small_1"
  description        = "acc-test-demo"
  project_name       = "default"
  address_ip_version = "DualStack"
}

# create ipv6 gateway
resource "vestack_vpc_ipv6_gateway" "ipv6_gateway" {
  vpc_id = vestack_vpc.vpc_ipv6.id
  name   = "acc-test-ipv6-gateway"
}

resource "vestack_vpc_ipv6_address_bandwidth" "foo" {
  ipv6_address = vestack_clb.private_clb_ipv6.eni_ipv6_address
  billing_type = "PostPaidByBandwidth"
  bandwidth    = 5
  depends_on   = [vestack_vpc_ipv6_gateway.ipv6_gateway]
}
```
## Argument Reference
The following arguments are supported:
* `subnet_id` - (Required, ForceNew) The id of the Subnet.
* `type` - (Required, ForceNew) The type of the CLB. And optional choice contains `public` or `private`.
* `address_ip_version` - (Optional, ForceNew) The address ip version of the Clb. Valid values: `ipv4`, `DualStack`. Default is `ipv4`.
When the value of this field is `DualStack`, the type of the CLB must be `private`, and suggest using a combination of resource `vestack_vpc_ipv6_gateway` and `vestack_vpc_ipv6_address_bandwidth` to achieve ipv6 public network access function.
* `description` - (Optional) The description of the CLB.
* `eip_billing_config` - (Optional, ForceNew) The billing configuration of the EIP which automatically associated to CLB. This field is valid when the type of CLB is `public`.When the type of the CLB is `private`, suggest using a combination of resource `vestack_eip_address` and `vestack_eip_associate` to achieve public network access function.
* `eni_address` - (Optional, ForceNew) The eni address of the CLB.
* `eni_ipv6_address` - (Optional, ForceNew) The eni ipv6 address of the Clb.
* `load_balancer_billing_type` - (Optional) The billing type of the CLB, valid values: `PostPaid`, `PrePaid`, `PostPaidByLCU`. Default is `PostPaid`.
* `load_balancer_name` - (Optional) The name of the CLB.
* `load_balancer_spec` - (Optional) The specification of the CLB, the value can be `small_1`, `small_2`, `medium_1`, `medium_2`, `large_1`, `large_2`. When the value of the `load_balancer_billing_type` is `PostPaidByLCU`, this field does not need to be specified.
* `master_zone_id` - (Optional) The master zone ID of the CLB.
* `modification_protection_reason` - (Optional) The reason of the console modification protection.
* `modification_protection_status` - (Optional) The status of the console modification protection, the value can be `NonProtection` or `ConsoleProtection`.
* `period` - (Optional) The period of the NatGateway, the valid value range in 1~9 or 12 or 24 or 36. Default value is 12. The period unit defaults to `Month`.This field is only effective when creating a PrePaid NatGateway. When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.
* `project_name` - (Optional) The ProjectName of the CLB.
* `region_id` - (Optional, ForceNew) The region of the request.
* `slave_zone_id` - (Optional) The slave zone ID of the CLB.
* `tags` - (Optional) Tags.
* `vpc_id` - (Optional, ForceNew) The id of the VPC.

The `eip_billing_config` object supports the following:

* `eip_billing_type` - (Required, ForceNew) The billing type of the EIP which automatically assigned to CLB. And optional choice contains `PostPaidByBandwidth` or `PostPaidByTraffic` or `PrePaid`.When creating a `PrePaid` public CLB, this field must be specified as `PrePaid` simultaneously.When the LoadBalancerBillingType changes from `PostPaid` to `PrePaid`, please manually modify the value of this field to `PrePaid` simultaneously.
* `isp` - (Required, ForceNew) The ISP of the EIP which automatically associated to CLB, the value can be `BGP` or `ChinaMobile` or `ChinaUnicom` or `ChinaTelecom` or `SingleLine_BGP` or `Static_BGP` or `Fusion_BGP`.
* `bandwidth` - (Optional) The peek bandwidth of the EIP which automatically assigned to CLB.

The `tags` object supports the following:

* `key` - (Required) The Key of Tags.
* `value` - (Required) The Value of Tags.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `eip_address` - The Eip address of the Clb.
* `eip_id` - The Eip ID of the Clb.
* `ipv6_eip_id` - The Ipv6 Eip ID of the Clb.
* `renew_type` - The renew type of the CLB. When the value of the load_balancer_billing_type is `PrePaid`, the query returns this field.


## Import
CLB can be imported using the id, e.g.
```
$ terraform import vestack_clb.default clb-273y2ok6ets007fap8txvf6us
```

