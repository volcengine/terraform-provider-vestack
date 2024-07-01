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