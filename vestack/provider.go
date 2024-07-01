package vestack

import (
	"github.com/volcengine/terraform-provider-vestack/logger"
	"github.com/volcengine/terraform-provider-vestack/vestack/tos/bucket_policy"
	//"github.com/volcengine/terraform-provider-vestack/vestack/cen/cen_service_route_entry"
	//"github.com/volcengine/terraform-provider-vestack/vestack/cloudfs/cloudfs_access"
	//"github.com/volcengine/terraform-provider-vestack/vestack/cloudfs/cloudfs_file_system"
	//"github.com/volcengine/terraform-provider-vestack/vestack/cloudfs/cloudfs_namespace"
	//"github.com/volcengine/terraform-provider-vestack/vestack/cloudfs/cloudfs_ns_quota"
	//"github.com/volcengine/terraform-provider-vestack/vestack/cloudfs/cloudfs_quota"
	//"github.com/volcengine/terraform-provider-vestack/vestack/nas/nas_file_system"
	//"github.com/volcengine/terraform-provider-vestack/vestack/nas/nas_region"
	//"github.com/volcengine/terraform-provider-vestack/vestack/nas/nas_snapshot"
	//"github.com/volcengine/terraform-provider-vestack/vestack/nas/nas_zone"
	"strings"

	//"github.com/volcengine/terraform-provider-vestack/vestack/fast_track/fast_track_associate"
	//"github.com/volcengine/terraform-provider-vestack/vestack/internet_tunnel/internet_tunnel_gateway"
	//"github.com/volcengine/terraform-provider-vestack/vestack/internet_tunnel/internet_tunnel_gateway_route"
	//"github.com/volcengine/terraform-provider-vestack/vestack/internet_tunnel/internet_tunnel_virtual_interface"
	//"github.com/volcengine/terraform-provider-vestack/vestack/kafka/kafka_region"
	//"github.com/volcengine/terraform-provider-vestack/vestack/kafka/kafka_zone"
	//"github.com/volcengine/terraform-provider-vestack/vestack/rds/rds_parameter_template"
	//"github.com/volcengine/terraform-provider-vestack/vestack/redis/instance_state"
	//"github.com/volcengine/terraform-provider-vestack/vestack/redis/pitr_time_period"
	//"github.com/volcengine/terraform-provider-vestack/vestack/vmp/vmp_alert"
	//"github.com/volcengine/terraform-provider-vestack/vestack/vmp/vmp_alerting_rule"
	//
	//"github.com/volcengine/terraform-provider-vestack/vestack/vmp/vmp_notify_group_policy"
	//"github.com/volcengine/terraform-provider-vestack/vestack/vmp/vmp_notify_policy"
	//
	//"github.com/volcengine/terraform-provider-vestack/vestack/vmp/vmp_contact"
	//"github.com/volcengine/terraform-provider-vestack/vestack/vmp/vmp_contact_group"
	//"github.com/volcengine/terraform-provider-vestack/vestack/vmp/vmp_instance_type"
	//"github.com/volcengine/terraform-provider-vestack/vestack/vmp/vmp_rule"
	//"github.com/volcengine/terraform-provider-vestack/vestack/vmp/vmp_rule_file"
	//"github.com/volcengine/terraform-provider-vestack/vestack/vmp/vmp_workspace"
	//
	//"github.com/volcengine/terraform-provider-vestack/vestack/mongodb/ssl_state"
	//"github.com/volcengine/terraform-provider-vestack/vestack/tls/alarm"
	//"github.com/volcengine/terraform-provider-vestack/vestack/tls/alarm_notify_group"
	//"github.com/volcengine/terraform-provider-vestack/vestack/tls/host"
	//"github.com/volcengine/terraform-provider-vestack/vestack/tls/host_group"
	//"github.com/volcengine/terraform-provider-vestack/vestack/tls/kafka_consumer"
	//tlsRule "github.com/volcengine/terraform-provider-vestack/vestack/tls/rule"
	//"github.com/volcengine/terraform-provider-vestack/vestack/tls/rule_applier"
	//"github.com/volcengine/terraform-provider-vestack/vestack/tls/shard"
	//"github.com/volcengine/terraform-provider-vestack/vestack/tos/bucket_policy"
	//
	//plSecurityGroup "github.com/volcengine/terraform-provider-vestack/vestack/privatelink/security_group"
	//"github.com/volcengine/terraform-provider-vestack/vestack/privatelink/vpc_endpoint"
	//"github.com/volcengine/terraform-provider-vestack/vestack/privatelink/vpc_endpoint_connection"
	//"github.com/volcengine/terraform-provider-vestack/vestack/privatelink/vpc_endpoint_service"
	//"github.com/volcengine/terraform-provider-vestack/vestack/privatelink/vpc_endpoint_service_permission"
	//"github.com/volcengine/terraform-provider-vestack/vestack/privatelink/vpc_endpoint_service_resource"
	//"github.com/volcengine/terraform-provider-vestack/vestack/privatelink/vpc_endpoint_zone"
	//
	//"github.com/volcengine/terraform-provider-vestack/vestack/mongodb/spec"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ve "github.com/volcengine/terraform-provider-vestack/common"
	//"github.com/volcengine/terraform-provider-vestack/vestack/anycast_eip/anycast_eip_address"
	//"github.com/volcengine/terraform-provider-vestack/vestack/anycast_eip/anycast_eip_associate"
	//"github.com/volcengine/terraform-provider-vestack/vestack/anycast_eip/anycast_pop_location"
	//"github.com/volcengine/terraform-provider-vestack/vestack/anycast_eip/anycast_server_region"
	//"github.com/volcengine/terraform-provider-vestack/vestack/autoscaling/scaling_activity"
	//"github.com/volcengine/terraform-provider-vestack/vestack/autoscaling/scaling_configuration"
	//"github.com/volcengine/terraform-provider-vestack/vestack/autoscaling/scaling_configuration_attachment"
	//"github.com/volcengine/terraform-provider-vestack/vestack/autoscaling/scaling_group"
	//"github.com/volcengine/terraform-provider-vestack/vestack/autoscaling/scaling_group_enabler"
	//"github.com/volcengine/terraform-provider-vestack/vestack/autoscaling/scaling_instance"
	//"github.com/volcengine/terraform-provider-vestack/vestack/autoscaling/scaling_instance_attachment"
	//"github.com/volcengine/terraform-provider-vestack/vestack/autoscaling/scaling_lifecycle_hook"
	//"github.com/volcengine/terraform-provider-vestack/vestack/autoscaling/scaling_policy"
	//bioosCluster "github.com/volcengine/terraform-provider-vestack/vestack/bioos/cluster"
	//"github.com/volcengine/terraform-provider-vestack/vestack/bioos/cluster_bind"
	//"github.com/volcengine/terraform-provider-vestack/vestack/bioos/workspace"
	//"github.com/volcengine/terraform-provider-vestack/vestack/cen/cen"
	//"github.com/volcengine/terraform-provider-vestack/vestack/cen/cen_attach_instance"
	//"github.com/volcengine/terraform-provider-vestack/vestack/cen/cen_bandwidth_package"
	//"github.com/volcengine/terraform-provider-vestack/vestack/cen/cen_bandwidth_package_associate"
	//"github.com/volcengine/terraform-provider-vestack/vestack/cen/cen_grant_instance"
	//"github.com/volcengine/terraform-provider-vestack/vestack/cen/cen_inter_region_bandwidth"
	//"github.com/volcengine/terraform-provider-vestack/vestack/cen/cen_route_entry"
	//"github.com/volcengine/terraform-provider-vestack/vestack/clb/acl"
	//"github.com/volcengine/terraform-provider-vestack/vestack/clb/acl_entry"
	//"github.com/volcengine/terraform-provider-vestack/vestack/clb/certificate"
	//"github.com/volcengine/terraform-provider-vestack/vestack/clb/clb"
	//"github.com/volcengine/terraform-provider-vestack/vestack/clb/listener"
	//"github.com/volcengine/terraform-provider-vestack/vestack/clb/rule"
	//"github.com/volcengine/terraform-provider-vestack/vestack/clb/server_group"
	//"github.com/volcengine/terraform-provider-vestack/vestack/clb/server_group_server"
	//clbZone "github.com/volcengine/terraform-provider-vestack/vestack/clb/zone"
	"github.com/volcengine/terraform-provider-vestack/vestack/cr/cr_authorization_token"
	"github.com/volcengine/terraform-provider-vestack/vestack/cr/cr_endpoint"
	"github.com/volcengine/terraform-provider-vestack/vestack/cr/cr_namespace"
	"github.com/volcengine/terraform-provider-vestack/vestack/cr/cr_registry"
	"github.com/volcengine/terraform-provider-vestack/vestack/cr/cr_registry_state"
	"github.com/volcengine/terraform-provider-vestack/vestack/cr/cr_repository"
	"github.com/volcengine/terraform-provider-vestack/vestack/cr/cr_tag"
	"github.com/volcengine/terraform-provider-vestack/vestack/cr/cr_vpc_endpoint"
	//"github.com/volcengine/terraform-provider-vestack/vestack/direct_connect/direct_connect_bgp_peer"
	//"github.com/volcengine/terraform-provider-vestack/vestack/direct_connect/direct_connect_connection"
	//"github.com/volcengine/terraform-provider-vestack/vestack/direct_connect/direct_connect_gateway"
	//"github.com/volcengine/terraform-provider-vestack/vestack/direct_connect/direct_connect_gateway_route"
	//"github.com/volcengine/terraform-provider-vestack/vestack/direct_connect/direct_connect_virtual_interface"
	"github.com/volcengine/terraform-provider-vestack/vestack/ebs/volume"
	"github.com/volcengine/terraform-provider-vestack/vestack/ebs/volume_attach"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_command"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_deployment_set"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_deployment_set_associate"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_instance"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_instance_state"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_invocation"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_invocation_result"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_key_pair"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_key_pair_associate"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/ecs_launch_template"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/image"
	"github.com/volcengine/terraform-provider-vestack/vestack/ecs/zone"
	"github.com/volcengine/terraform-provider-vestack/vestack/eip/eip_address"
	"github.com/volcengine/terraform-provider-vestack/vestack/eip/eip_associate"
	//"github.com/volcengine/terraform-provider-vestack/vestack/escloud/instance"
	//"github.com/volcengine/terraform-provider-vestack/vestack/escloud/region"
	//esZone "github.com/volcengine/terraform-provider-vestack/vestack/escloud/zone"
	//"github.com/volcengine/terraform-provider-vestack/vestack/fast_track/fast_track"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_access_key"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_login_profile"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_policy"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_role"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_role_policy_attachment"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_service_linked_role"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_user"
	"github.com/volcengine/terraform-provider-vestack/vestack/iam/iam_user_policy_attachment"
	//plbAcl "github.com/volcengine/terraform-provider-vestack/vestack/inner/plb/acl"
	//plbAclEntry "github.com/volcengine/terraform-provider-vestack/vestack/inner/plb/acl_entry"
	//plbListener "github.com/volcengine/terraform-provider-vestack/vestack/inner/plb/listener"
	//"github.com/volcengine/terraform-provider-vestack/vestack/inner/plb/plb"
	//"github.com/volcengine/terraform-provider-vestack/vestack/inner/plb/plb_vpc_associate"
	//plbServerGroup "github.com/volcengine/terraform-provider-vestack/vestack/inner/plb/server_group"
	//plbServerGroupServer "github.com/volcengine/terraform-provider-vestack/vestack/inner/plb/server_group_server"
	//"github.com/volcengine/terraform-provider-vestack/vestack/inner/shuttle/shuttle"
	//"github.com/volcengine/terraform-provider-vestack/vestack/inner/shuttle/shuttle_client"
	//"github.com/volcengine/terraform-provider-vestack/vestack/inner/shuttle/shuttle_server"
	//"github.com/volcengine/terraform-provider-vestack/vestack/inner/shuttle_v1/shuttle_association_v1"
	//"github.com/volcengine/terraform-provider-vestack/vestack/inner/shuttle_v1/shuttle_client_v1"
	//"github.com/volcengine/terraform-provider-vestack/vestack/inner/shuttle_v1/shuttle_server_v1"
	//trEntry "github.com/volcengine/terraform-provider-vestack/vestack/inner/transit_router/route_entry"
	//trTable "github.com/volcengine/terraform-provider-vestack/vestack/inner/transit_router/route_table"
	//"github.com/volcengine/terraform-provider-vestack/vestack/inner/transit_router/route_table_association"
	//"github.com/volcengine/terraform-provider-vestack/vestack/inner/transit_router/route_table_propagation"
	//"github.com/volcengine/terraform-provider-vestack/vestack/inner/transit_router/transit_router"
	//"github.com/volcengine/terraform-provider-vestack/vestack/inner/transit_router/transit_router_vpc_attachment"
	//"github.com/volcengine/terraform-provider-vestack/vestack/inner/transit_router/transit_router_vpn_attachment"
	//"github.com/volcengine/terraform-provider-vestack/vestack/internet_tunnel/internet_tunnel_address"
	//"github.com/volcengine/terraform-provider-vestack/vestack/internet_tunnel/internet_tunnel_address_associate"
	//"github.com/volcengine/terraform-provider-vestack/vestack/internet_tunnel/internet_tunnel_bandwidth"
	//"github.com/volcengine/terraform-provider-vestack/vestack/internet_tunnel/internet_tunnel_bgp_peer"
	//"github.com/volcengine/terraform-provider-vestack/vestack/kafka/kafka_consumed_partition"
	//"github.com/volcengine/terraform-provider-vestack/vestack/kafka/kafka_consumed_topic"
	//"github.com/volcengine/terraform-provider-vestack/vestack/kafka/kafka_group"
	//"github.com/volcengine/terraform-provider-vestack/vestack/kafka/kafka_instance"
	//"github.com/volcengine/terraform-provider-vestack/vestack/kafka/kafka_public_address"
	//"github.com/volcengine/terraform-provider-vestack/vestack/kafka/kafka_sasl_user"
	//"github.com/volcengine/terraform-provider-vestack/vestack/kafka/kafka_topic"
	//"github.com/volcengine/terraform-provider-vestack/vestack/kafka/kafka_topic_partition"
	//"github.com/volcengine/terraform-provider-vestack/vestack/mongodb/account"
	//"github.com/volcengine/terraform-provider-vestack/vestack/mongodb/allow_list"
	//"github.com/volcengine/terraform-provider-vestack/vestack/mongodb/allow_list_associate"
	//"github.com/volcengine/terraform-provider-vestack/vestack/mongodb/endpoint"
	//mongodbInstance "github.com/volcengine/terraform-provider-vestack/vestack/mongodb/instance"
	//"github.com/volcengine/terraform-provider-vestack/vestack/mongodb/instance_parameter"
	//"github.com/volcengine/terraform-provider-vestack/vestack/mongodb/instance_parameter_log"
	//mongodbRegion "github.com/volcengine/terraform-provider-vestack/vestack/mongodb/region"
	//mongodbZone "github.com/volcengine/terraform-provider-vestack/vestack/mongodb/zone"
	//"github.com/volcengine/terraform-provider-vestack/vestack/nat/dnat_entry"
	//"github.com/volcengine/terraform-provider-vestack/vestack/nat/nat_gateway"
	//"github.com/volcengine/terraform-provider-vestack/vestack/nat/snat_entry"
	//"github.com/volcengine/terraform-provider-vestack/vestack/rds/rds_account"
	//"github.com/volcengine/terraform-provider-vestack/vestack/rds/rds_account_privilege"
	//"github.com/volcengine/terraform-provider-vestack/vestack/rds/rds_database"
	//"github.com/volcengine/terraform-provider-vestack/vestack/rds/rds_instance"
	//"github.com/volcengine/terraform-provider-vestack/vestack/rds/rds_ip_list"
	//"github.com/volcengine/terraform-provider-vestack/vestack/rds_mysql/allowlist"
	//"github.com/volcengine/terraform-provider-vestack/vestack/rds_mysql/allowlist_associate"
	//"github.com/volcengine/terraform-provider-vestack/vestack/rds_mysql/rds_mysql_account"
	//"github.com/volcengine/terraform-provider-vestack/vestack/rds_mysql/rds_mysql_database"
	//"github.com/volcengine/terraform-provider-vestack/vestack/rds_mysql/rds_mysql_instance"
	//"github.com/volcengine/terraform-provider-vestack/vestack/rds_mysql/rds_mysql_instance_readonly_node"
	//"github.com/volcengine/terraform-provider-vestack/vestack/rds_v2/rds_instance_v2"
	//
	//tlsIndex "github.com/volcengine/terraform-provider-vestack/vestack/tls/index"
	//tlsProject "github.com/volcengine/terraform-provider-vestack/vestack/tls/project"
	//tlsTopic "github.com/volcengine/terraform-provider-vestack/vestack/tls/topic"
	//
	//redisAccount "github.com/volcengine/terraform-provider-vestack/vestack/redis/account"
	//redis_allow_list "github.com/volcengine/terraform-provider-vestack/vestack/redis/allow_list"
	//redis_allow_list_associate "github.com/volcengine/terraform-provider-vestack/vestack/redis/allow_list_associate"
	//redis_backup "github.com/volcengine/terraform-provider-vestack/vestack/redis/backup"
	//redis_backup_restore "github.com/volcengine/terraform-provider-vestack/vestack/redis/backup_restore"
	//redisContinuousBackup "github.com/volcengine/terraform-provider-vestack/vestack/redis/continuous_backup"
	//redis_endpoint "github.com/volcengine/terraform-provider-vestack/vestack/redis/endpoint"
	//redisInstance "github.com/volcengine/terraform-provider-vestack/vestack/redis/instance"
	//redisRegion "github.com/volcengine/terraform-provider-vestack/vestack/redis/region"
	//redisZone "github.com/volcengine/terraform-provider-vestack/vestack/redis/zone"

	"github.com/volcengine/terraform-provider-vestack/vestack/tos/bucket"
	"github.com/volcengine/terraform-provider-vestack/vestack/tos/object"
	//"github.com/volcengine/terraform-provider-vestack/vestack/veenedge/available_resource"
	//"github.com/volcengine/terraform-provider-vestack/vestack/veenedge/cloud_server"
	//veInstance "github.com/volcengine/terraform-provider-vestack/vestack/veenedge/instance"
	//"github.com/volcengine/terraform-provider-vestack/vestack/veenedge/instance_types"
	//veVpc "github.com/volcengine/terraform-provider-vestack/vestack/veenedge/vpc"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/addon"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/cluster"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/default_node_pool"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/default_node_pool_batch_attach"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/kubeconfig"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/node"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/node_pool"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/support_addon"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/ipv6_address"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/ipv6_address_bandwidth"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/ipv6_gateway"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/network_acl"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/network_acl_associate"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/network_interface"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/network_interface_attach"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/route_entry"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/route_table"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/route_table_associate"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/security_group"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/security_group_rule"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/subnet"
	"github.com/volcengine/terraform-provider-vestack/vestack/vpc/vpc"
	//"github.com/volcengine/terraform-provider-vestack/vestack/vpn/customer_gateway"
	//"github.com/volcengine/terraform-provider-vestack/vestack/vpn/vpn_connection"
	//"github.com/volcengine/terraform-provider-vestack/vestack/vpn/vpn_gateway"
	//"github.com/volcengine/terraform-provider-vestack/vestack/vpn/vpn_gateway_route"
	//
	//"github.com/volcengine/terraform-provider-vestack/vestack/nas/nas_mount_point"
	//"github.com/volcengine/terraform-provider-vestack/vestack/nas/nas_permission_group"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VESTACK_ACCESS_KEY", nil),
				Description: "The Access Key for Vestack Provider",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VESTACK_SECRET_KEY", nil),
				Description: "The Secret Key for Vestack Provider",
			},
			"session_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VESTACK_SESSION_TOKEN", nil),
				Description: "The Session Token for Vestack Provider",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VESTACK_REGION", nil),
				Description: "The Region for Vestack Provider",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VESTACK_ENDPOINT", nil),
				Description: "The Customer Endpoint for Vestack Provider",
			},
			"disable_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VESTACK_DISABLE_SSL", nil),
				Description: "Disable SSL for Vestack Provider",
			},
			"customer_headers": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VESTACK_CUSTOMER_HEADERS", nil),
				Description: "CUSTOMER HEADERS for Vestack Provider",
			},
			"customer_endpoints": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VESTACK_CUSTOMER_ENDPOINTS", nil),
				Description: "CUSTOMER ENDPOINTS for Vestack Provider",
			},
			"proxy_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VESTACK_PROXY_URL", nil),
				Description: "PROXY URL for Vestack Provider",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"vestack_vpcs":                        vpc.DataSourceVestackVpcs(),
			"vestack_subnets":                     subnet.DataSourceVestackSubnets(),
			"vestack_route_tables":                route_table.DataSourceVestackRouteTables(),
			"vestack_route_entries":               route_entry.DataSourceVestackRouteEntries(),
			"vestack_security_groups":             security_group.DataSourceVestackSecurityGroups(),
			"vestack_security_group_rules":        security_group_rule.DataSourceVestackSecurityGroupRules(),
			"vestack_network_interfaces":          network_interface.DataSourceVestackNetworkInterfaces(),
			"vestack_network_acls":                network_acl.DataSourceVestackNetworkAcls(),
			"vestack_vpc_ipv6_gateways":           ipv6_gateway.DataSourceVestackIpv6Gateways(),
			"vestack_vpc_ipv6_address_bandwidths": ipv6_address_bandwidth.DataSourceVestackIpv6AddressBandwidths(),
			"vestack_vpc_ipv6_addresses":          ipv6_address.DataSourceVestackIpv6Addresses(),

			// ================ EIP ================
			"vestack_eip_addresses": eip_address.DataSourceVestackEipAddresses(),

			// ================ AnyCast Eip ================
			//"vestack_anycast_eip_addresses":  anycast_eip_address.DataSourceVestackAnyCastEipAddresses(),
			//"vestack_anycast_pop_locations":  anycast_pop_location.DataSourceVestackAnyCastPopLocations(),
			//"vestack_anycast_server_regions": anycast_server_region.DataSourceVestackAnyCastServerRegions(),

			// ================ CLB ================
			//"vestack_acls":                 acl.DataSourceVestackAcls(),
			//"vestack_clbs":                 clb.DataSourceVestackClbs(),
			//"vestack_listeners":            listener.DataSourceVestackListeners(),
			//"vestack_server_groups":        server_group.DataSourceVestackServerGroups(),
			//"vestack_certificates":         certificate.DataSourceVestackCertificates(),
			//"vestack_clb_rules":            rule.DataSourceVestackRules(),
			//"vestack_server_group_servers": server_group_server.DataSourceVestackServerGroupServers(),
			//"vestack_clb_zones":            clbZone.DataSourceVestackClbZones(),

			// ================ EBS ================
			"vestack_volumes": volume.DataSourceVestackVolumes(),

			// ================ ECS ================
			"vestack_ecs_instances":          ecs_instance.DataSourceVestackEcsInstances(),
			"vestack_images":                 image.DataSourceVestackImages(),
			"vestack_zones":                  zone.DataSourceVestackZones(),
			"vestack_ecs_deployment_sets":    ecs_deployment_set.DataSourceVestackEcsDeploymentSets(),
			"vestack_ecs_key_pairs":          ecs_key_pair.DataSourceVestackEcsKeyPairs(),
			"vestack_ecs_launch_templates":   ecs_launch_template.DataSourceVestackEcsLaunchTemplates(),
			"vestack_ecs_commands":           ecs_command.DataSourceVestackEcsCommands(),
			"vestack_ecs_invocations":        ecs_invocation.DataSourceVestackEcsInvocations(),
			"vestack_ecs_invocation_results": ecs_invocation_result.DataSourceVestackEcsInvocationResults(),

			// ================ NAT ================
			//"vestack_snat_entries": snat_entry.DataSourceVestackSnatEntries(),
			//"vestack_nat_gateways": nat_gateway.DataSourceVestackNatGateways(),
			//"vestack_dnat_entries": dnat_entry.DataSourceVestackDnatEntries(),

			// ================ AutoScaling ================
			//"vestack_scaling_groups":          scaling_group.DataSourceVestackScalingGroups(),
			//"vestack_scaling_configurations":  scaling_configuration.DataSourceVestackScalingConfigurations(),
			//"vestack_scaling_policies":        scaling_policy.DataSourceVestackScalingPolicies(),
			//"vestack_scaling_activities":      scaling_activity.DataSourceVestackScalingActivities(),
			//"vestack_scaling_lifecycle_hooks": scaling_lifecycle_hook.DataSourceVestackScalingLifecycleHooks(),
			//"vestack_scaling_instances":       scaling_instance.DataSourceVestackScalingInstances(),

			// ================ Cen ================
			//"vestack_cens":                        cen.DataSourceVestackCens(),
			//"vestack_cen_attach_instances":        cen_attach_instance.DataSourceVestackCenAttachInstances(),
			//"vestack_cen_bandwidth_packages":      cen_bandwidth_package.DataSourceVestackCenBandwidthPackages(),
			//"vestack_cen_inter_region_bandwidths": cen_inter_region_bandwidth.DataSourceVestackCenInterRegionBandwidths(),
			//"vestack_cen_service_route_entries":   cen_service_route_entry.DataSourceVestackCenServiceRouteEntries(),
			//"vestack_cen_route_entries":           cen_route_entry.DataSourceVestackCenRouteEntries(),

			// ================ VPN ================
			//"vestack_vpn_gateways":       vpn_gateway.DataSourceVestackVpnGateways(),
			//"vestack_customer_gateways":  customer_gateway.DataSourceVestackCustomerGateways(),
			//"vestack_vpn_connections":    vpn_connection.DataSourceVestackVpnConnections(),
			//"vestack_vpn_gateway_routes": vpn_gateway_route.DataSourceVestackVpnGatewayRoutes(),

			// ================ VKE ================
			"vestack_vke_nodes":          node.DataSourceVestackVkeNodes(),
			"vestack_vke_clusters":       cluster.DataSourceVestackVkeVkeClusters(),
			"vestack_vke_node_pools":     node_pool.DataSourceVestackNodePools(),
			"vestack_vke_addons":         addon.DataSourceVestackVkeAddons(),
			"vestack_vke_support_addons": support_addon.DataSourceVestackVkeVkeSupportedAddons(),
			"vestack_vke_kubeconfigs":    kubeconfig.DataSourceVestackVkeKubeconfigs(),

			// ================ IAM ================
			"vestack_iam_policies": iam_policy.DataSourceVestackIamPolicies(),
			"vestack_iam_roles":    iam_role.DataSourceVestackIamRoles(),
			"vestack_iam_users":    iam_user.DataSourceVestackIamUsers(),

			// ================ RDS V1 ==============
			//"vestack_rds_instances":           rds_instance.DataSourceVestackRdsInstances(),
			//"vestack_rds_databases":           rds_database.DataSourceVestackRdsDatabases(),
			//"vestack_rds_accounts":            rds_account.DataSourceVestackRdsAccounts(),
			//"vestack_rds_ip_lists":            rds_ip_list.DataSourceVestackRdsIpLists(),
			//"vestack_rds_parameter_templates": rds_parameter_template.DataSourceVestackRdsParameterTemplates(),

			// ================ ESCloud =============
			//"vestack_escloud_instances": instance.DataSourceVestackESCloudInstances(),
			//"vestack_escloud_regions":   region.DataSourceVestackESCloudRegions(),
			//"vestack_escloud_zones":     esZone.DataSourceVestackESCloudZones(),

			// ================ TOS ================
			"vestack_tos_buckets": bucket.DataSourceVestackTosBuckets(),
			"vestack_tos_objects": object.DataSourceVestackTosObjects(),

			// ================ PLB ================
			//"vestack_plbs":                     plb.DataSourceVestackPlbs(),
			//"vestack_plb_listeners":            plbListener.DataSourceVestackListeners(),
			//"vestack_plb_server_groups":        plbServerGroup.DataSourceVestackServerGroups(),
			//"vestack_plb_server_group_servers": plbServerGroupServer.DataSourceVestackServerGroupServers(),
			//"vestack_plb_acls":                 plbAcl.DataSourceVestackAcls(),

			// ================ Redis =============
			//"vestack_redis_allow_lists":       redis_allow_list.DataSourceVestackRedisAllowLists(),
			//"vestack_redis_backups":           redis_backup.DataSourceVestackRedisBackups(),
			//"vestack_redis_regions":           redisRegion.DataSourceVestackRedisRegions(),
			//"vestack_redis_zones":             redisZone.DataSourceVestackRedisZones(),
			//"vestack_redis_accounts":          redisAccount.DataSourceVestackRedisAccounts(),
			//"vestack_redis_instances":         redisInstance.DataSourceVestackRedisDbInstances(),
			//"vestack_redis_pitr_time_windows": pitr_time_period.DataSourceVestackRedisPitrTimeWindows(),

			// ================ CR ================
			"vestack_cr_registries":           cr_registry.DataSourceVestackCrRegistries(),
			"vestack_cr_namespaces":           cr_namespace.DataSourceVestackCrNamespaces(),
			"vestack_cr_repositories":         cr_repository.DataSourceVestackCrRepositories(),
			"vestack_cr_tags":                 cr_tag.DataSourceVestackCrTags(),
			"vestack_cr_authorization_tokens": cr_authorization_token.DataSourceVestackCrAuthorizationTokens(),
			"vestack_cr_endpoints":            cr_endpoint.DataSourceVestackCrEndpoints(),
			"vestack_cr_vpc_endpoints":        cr_vpc_endpoint.DataSourceVestackCrVpcEndpoints(),

			// ================ Shuttle =============================
			//"vestack_shuttles":                shuttle.DataSourceVestackShuttles(),
			//"vestack_shuttle_clients":         shuttle_client.DataSourceVestackShuttleClients(),
			//"vestack_shuttle_servers":         shuttle_server.DataSourceVestackShuttleServers(),
			//"vestack_shuttle_associations_v1": shuttle_association_v1.DataSourceVestackShuttleAssociationsV1(),
			//"vestack_shuttle_clients_v1":      shuttle_client_v1.DataSourceVestackShuttleClientsV1(),
			//"vestack_shuttle_servers_v1":      shuttle_server_v1.DataSourceVestackShuttleServersV1(),
			//
			//// ================ Veenedge ================
			//"vestack_veenedge_cloud_servers":       cloud_server.DataSourceVestackVeenedgeCloudServers(),
			//"vestack_veenedge_instances":           veInstance.DataSourceVestackInstances(),
			//"vestack_veenedge_instance_types":      instance_types.DataSourceVestackInstanceTypes(),
			//"vestack_veenedge_available_resources": available_resource.DataSourceVestackAvailableResources(),
			//"vestack_veenedge_vpcs":                veVpc.DataSourceVestackVpcs(),
			//
			//// ================ Kafka ================
			//"vestack_kafka_sasl_users":          kafka_sasl_user.DataSourceVestackKafkaSaslUsers(),
			//"vestack_kafka_topic_partitions":    kafka_topic_partition.DataSourceVestackKafkaTopicPartitions(),
			//"vestack_kafka_groups":              kafka_group.DataSourceVestackKafkaGroups(),
			//"vestack_kafka_topics":              kafka_topic.DataSourceVestackKafkaTopics(),
			//"vestack_kafka_instances":           kafka_instance.DataSourceVestackKafkaInstances(),
			//"vestack_kafka_regions":             kafka_region.DataSourceVestackRegions(),
			//"vestack_kafka_zones":               kafka_zone.DataSourceVestackZones(),
			//"vestack_kafka_consumed_topics":     kafka_consumed_topic.DataSourceVestackKafkaConsumedTopics(),
			//"vestack_kafka_consumed_partitions": kafka_consumed_partition.DataSourceVestackKafkaConsumedPartitions(),
			//
			//// ================ MongoDB =============
			//"vestack_mongodb_instances":               mongodbInstance.DataSourceVestackMongoDBInstances(),
			//"vestack_mongodb_endpoints":               endpoint.DataSourceVestackMongoDBEndpoints(),
			//"vestack_mongodb_allow_lists":             allow_list.DataSourceVestackMongoDBAllowLists(),
			//"vestack_mongodb_instance_parameters":     instance_parameter.DataSourceVestackMongoDBInstanceParameters(),
			//"vestack_mongodb_instance_parameter_logs": instance_parameter_log.DataSourceVestackMongoDBInstanceParameterLogs(),
			//"vestack_mongodb_regions":                 mongodbRegion.DataSourceVestackMongoDBRegions(),
			//"vestack_mongodb_zones":                   mongodbZone.DataSourceVestackMongoDBZones(),
			//"vestack_mongodb_accounts":                account.DataSourceVestackMongoDBAccounts(),
			//"vestack_mongodb_specs":                   spec.DataSourceVestackMongoDBSpecs(),
			//"vestack_mongodb_ssl_states":              ssl_state.DataSourceVestackMongoDBSSLStates(),
			//
			//// ================ Bioos ==================
			//"vestack_bioos_clusters":   bioosCluster.DataSourceVestackBioosClusters(),
			//"vestack_bioos_workspaces": workspace.DataSourceVestackBioosWorkspaces(),
			//
			//// ================ DirectConnect ================
			//"vestack_direct_connect_connections":        direct_connect_connection.DataSourceVestackDirectConnectConnections(),
			//"vestack_direct_connect_gateways":           direct_connect_gateway.DataSourceVestackDirectConnectGateways(),
			//"vestack_direct_connect_virtual_interfaces": direct_connect_virtual_interface.DataSourceVestackDirectConnectVirtualInterfaces(),
			//"vestack_direct_connect_bgp_peers":          direct_connect_bgp_peer.DataSourceVestackDirectConnectBgpPeers(),
			//"vestack_direct_connect_gateway_routes":     direct_connect_gateway_route.DataSourceVestackDirectConnectGatewayRoutes(),
			//
			//// ================ TransitRouter =============
			//"vestack_transit_routers":                         transit_router.DataSourceVestackTransitRouters(),
			//"vestack_transit_router_vpc_attachments":          transit_router_vpc_attachment.DataSourceVestackTransitRouterVpcAttachments(),
			//"vestack_transit_router_vpn_attachments":          transit_router_vpn_attachment.DataSourceVestackTransitRouterVpnAttachments(),
			//"vestack_transit_router_route_tables":             trTable.DataSourceVestackTransitRouterRouteTables(),
			//"vestack_transit_router_route_entries":            trEntry.DataSourceVestackTransitRouterRouteEntries(),
			//"vestack_transit_router_route_table_associations": route_table_association.DataSourceVestackTransitRouterRouteTableAssociations(),
			//"vestack_transit_router_route_table_propagations": route_table_propagation.DataSourceVestackTransitRouterRouteTablePropagations(),
			//
			//// ================ PrivateLink ==================
			//"vestack_privatelink_vpc_endpoints":                    vpc_endpoint.DataSourceVestackPrivatelinkVpcEndpoints(),
			//"vestack_privatelink_vpc_endpoint_services":            vpc_endpoint_service.DataSourceVestackPrivatelinkVpcEndpointServices(),
			//"vestack_privatelink_vpc_endpoint_service_permissions": vpc_endpoint_service_permission.DataSourceVestackPrivatelinkVpcEndpointServicePermissions(),
			//"vestack_privatelink_vpc_endpoint_connections":         vpc_endpoint_connection.DataSourceVestackPrivatelinkVpcEndpointConnections(),
			//"vestack_privatelink_vpc_endpoint_zones":               vpc_endpoint_zone.DataSourceVestackPrivatelinkVpcEndpointZones(),
			//
			//// ================ RDS V2 ==============
			//"vestack_rds_instances_v2": rds_instance_v2.DataSourceVestackRdsInstances(),
			//
			//// ================ RdsMysql ================
			//"vestack_rds_mysql_instances":  rds_mysql_instance.DataSourceVestackRdsMysqlInstances(),
			//"vestack_rds_mysql_accounts":   rds_mysql_account.DataSourceVestackRdsMysqlAccounts(),
			//"vestack_rds_mysql_databases":  rds_mysql_database.DataSourceVestackRdsMysqlDatabases(),
			//"vestack_rds_mysql_allowlists": allowlist.DataSourceVestackRdsMysqlAllowLists(),
			//
			//// ================ TLS ================
			//"vestack_tls_rules":               tlsRule.DataSourceVestackTlsRules(),
			//"vestack_tls_alarms":              alarm.DataSourceVestackTlsAlarms(),
			//"vestack_tls_alarm_notify_groups": alarm_notify_group.DataSourceVestackTlsAlarmNotifyGroups(),
			//"vestack_tls_rule_appliers":       rule_applier.DataSourceVestackTlsRuleAppliers(),
			//"vestack_tls_shards":              shard.DataSourceVestackTlsShards(),
			//"vestack_tls_kafka_consumers":     kafka_consumer.DataSourceVestackTlsKafkaConsumers(),
			//"vestack_tls_host_groups":         host_group.DataSourceVestackTlsHostGroups(),
			//"vestack_tls_hosts":               host.DataSourceVestackTlsHosts(),
			//"vestack_tls_projects":            tlsProject.DataSourceVestackTlsProjects(),
			//"vestack_tls_topics":              tlsTopic.DataSourceVestackTlsTopics(),
			//"vestack_tls_indexes":             tlsIndex.DataSourceVestackTlsIndexes(),
			//
			//// ================ VMP ================
			//"vestack_vmp_workspaces":            vmp_workspace.DataSourceVestackVmpWorkspaces(),
			//"vestack_vmp_instance_types":        vmp_instance_type.DataSourceVestackVmpInstanceTypes(),
			//"vestack_vmp_rule_files":            vmp_rule_file.DataSourceVestackVmpRuleFiles(),
			//"vestack_vmp_rules":                 vmp_rule.DataSourceVestackVmpRules(),
			//"vestack_vmp_contact_groups":        vmp_contact_group.DataSourceVestackVmpContactGroups(),
			//"vestack_vmp_contacts":              vmp_contact.DataSourceVestackVmpContacts(),
			//"vestack_vmp_alerting_rules":        vmp_alerting_rule.DataSourceVestackVmpAlertingRules(),
			//"vestack_vmp_alerts":                vmp_alert.DataSourceVestackVmpAlerts(),
			//"vestack_vmp_notify_group_policies": vmp_notify_group_policy.DataSourceVestackVmpNotifyGroupPolicies(),
			//"vestack_vmp_notify_policies":       vmp_notify_policy.DataSourceVestackVmpNotifyPolicies(),
			//
			//// ================ InternetTunnel ================
			//"vestack_internet_tunnel_gateways":           internet_tunnel_gateway.DataSourceVestackInternetTunnelGateways(),
			//"vestack_internet_tunnel_gateway_routes":     internet_tunnel_gateway_route.DataSourceVestackInternetTunnelGatewayRoutes(),
			//"vestack_internet_tunnel_virtual_interfaces": internet_tunnel_virtual_interface.DataSourceVestackInternetTunnelVirtualInterfaces(),
			//"vestack_internet_tunnel_bgp_peers":          internet_tunnel_bgp_peer.DataSourceVestackInternetTunnelBgpPeers(),
			//"vestack_internet_tunnel_bandwidths":         internet_tunnel_bandwidth.DataSourceVestackInternetTunnelBandwidths(),
			//"vestack_internet_tunnel_addresses":          internet_tunnel_address.DataSourceVestackInternetTunnelAddresses(),
			//
			//// ================ FastTrack ================
			//"vestack_fast_tracks": fast_track.DataSourceVestackFastTracks(),
			//
			//// ================ NAS ================
			//"vestack_nas_file_systems":      nas_file_system.DataSourceVestackNasFileSystems(),
			//"vestack_nas_snapshots":         nas_snapshot.DataSourceVestackNasSnapshots(),
			//"vestack_nas_zones":             nas_zone.DataSourceVestackNasZones(),
			//"vestack_nas_regions":           nas_region.DataSourceVestackNasRegions(),
			//"vestack_nas_mount_points":      nas_mount_point.DataSourceVestackNasMountPoints(),
			//"vestack_nas_permission_groups": nas_permission_group.DataSourceVestackNasPermissionGroups(),
			//
			//// ================ Cloudfs ================
			//"vestack_cloudfs_quotas":       cloudfs_quota.DataSourceVestackCloudfsQuotas(),
			//"vestack_cloudfs_file_systems": cloudfs_file_system.DataSourceVestackCloudfsFileSystems(),
			//"vestack_cloudfs_accesses":     cloudfs_access.DataSourceVestackCloudfsAccesses(),
			//"vestack_cloudfs_ns_quotas":    cloudfs_ns_quota.DataSourceVestackCloudfsNsQuotas(),
			//"vestack_cloudfs_namespaces":   cloudfs_namespace.DataSourceVestackCloudfsNamespaces(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"vestack_vpc":                        vpc.ResourceVestackVpc(),
			"vestack_subnet":                     subnet.ResourceVestackSubnet(),
			"vestack_route_table":                route_table.ResourceVestackRouteTable(),
			"vestack_route_entry":                route_entry.ResourceVestackRouteEntry(),
			"vestack_route_table_associate":      route_table_associate.ResourceVestackRouteTableAssociate(),
			"vestack_security_group":             security_group.ResourceVestackSecurityGroup(),
			"vestack_network_interface":          network_interface.ResourceVestackNetworkInterface(),
			"vestack_network_interface_attach":   network_interface_attach.ResourceVestackNetworkInterfaceAttach(),
			"vestack_security_group_rule":        security_group_rule.ResourceVestackSecurityGroupRule(),
			"vestack_network_acl":                network_acl.ResourceVestackNetworkAcl(),
			"vestack_network_acl_associate":      network_acl_associate.ResourceVestackNetworkAclAssociate(),
			"vestack_vpc_ipv6_gateway":           ipv6_gateway.ResourceVestackIpv6Gateway(),
			"vestack_vpc_ipv6_address_bandwidth": ipv6_address_bandwidth.ResourceVestackIpv6AddressBandwidth(),

			// ================ EIP ================
			"vestack_eip_address":   eip_address.ResourceVestackEipAddress(),
			"vestack_eip_associate": eip_associate.ResourceVestackEipAssociate(),

			// ================ AnyCast Eip ================
			//"vestack_anycast_eip_address":   anycast_eip_address.ResourceVestackAnyCastEipAddress(),
			//"vestack_anycast_eip_associate": anycast_eip_associate.ResourceVestackAnycastEipAssociate(),
			//
			//// ================ CLB ================
			//"vestack_acl":                 acl.ResourceVestackAcl(),
			//"vestack_clb":                 clb.ResourceVestackClb(),
			//"vestack_listener":            listener.ResourceVestackListener(),
			//"vestack_server_group":        server_group.ResourceVestackServerGroup(),
			//"vestack_certificate":         certificate.ResourceVestackCertificate(),
			//"vestack_clb_rule":            rule.ResourceVestackRule(),
			//"vestack_server_group_server": server_group_server.ResourceVestackServerGroupServer(),
			//"vestack_acl_entry":           acl_entry.ResourceVestackAclEntry(),

			// ================ EBS ================
			"vestack_volume":        volume.ResourceVestackVolume(),
			"vestack_volume_attach": volume_attach.ResourceVestackVolumeAttach(),

			// ================ ECS ================
			"vestack_ecs_instance":                 ecs_instance.ResourceVestackEcsInstance(),
			"vestack_ecs_instance_state":           ecs_instance_state.ResourceVestackEcsInstanceState(),
			"vestack_ecs_deployment_set":           ecs_deployment_set.ResourceVestackEcsDeploymentSet(),
			"vestack_ecs_deployment_set_associate": ecs_deployment_set_associate.ResourceVestackEcsDeploymentSetAssociate(),
			"vestack_ecs_key_pair":                 ecs_key_pair.ResourceVestackEcsKeyPair(),
			"vestack_ecs_key_pair_associate":       ecs_key_pair_associate.ResourceVestackEcsKeyPairAssociate(),
			"vestack_ecs_launch_template":          ecs_launch_template.ResourceVestackEcsLaunchTemplate(),
			"vestack_ecs_command":                  ecs_command.ResourceVestackEcsCommand(),
			"vestack_ecs_invocation":               ecs_invocation.ResourceVestackEcsInvocation(),

			// ================ NAT ================
			//"vestack_snat_entry":  snat_entry.ResourceVestackSnatEntry(),
			//"vestack_nat_gateway": nat_gateway.ResourceVestackNatGateway(),
			//"vestack_dnat_entry":  dnat_entry.ResourceVestackDnatEntry(),
			//
			//// ================ AutoScaling ================
			//"vestack_scaling_group":                    scaling_group.ResourceVestackScalingGroup(),
			//"vestack_scaling_configuration":            scaling_configuration.ResourceVestackScalingConfiguration(),
			//"vestack_scaling_policy":                   scaling_policy.ResourceVestackScalingPolicy(),
			//"vestack_scaling_instance_attachment":      scaling_instance_attachment.ResourceVestackScalingInstanceAttachment(),
			//"vestack_scaling_lifecycle_hook":           scaling_lifecycle_hook.ResourceVestackScalingLifecycleHook(),
			//"vestack_scaling_group_enabler":            scaling_group_enabler.ResourceVestackScalingGroupEnabler(),
			//"vestack_scaling_configuration_attachment": scaling_configuration_attachment.ResourceVestackScalingConfigurationAttachment(),
			//
			//// ================ Cen ================
			//"vestack_cen":                             cen.ResourceVestackCen(),
			//"vestack_cen_attach_instance":             cen_attach_instance.ResourceVestackCenAttachInstance(),
			//"vestack_cen_grant_instance":              cen_grant_instance.ResourceVestackCenGrantInstance(),
			//"vestack_cen_bandwidth_package":           cen_bandwidth_package.ResourceVestackCenBandwidthPackage(),
			//"vestack_cen_bandwidth_package_associate": cen_bandwidth_package_associate.ResourceVestackCenBandwidthPackageAssociate(),
			//"vestack_cen_inter_region_bandwidth":      cen_inter_region_bandwidth.ResourceVestackCenInterRegionBandwidth(),
			//"vestack_cen_service_route_entry":         cen_service_route_entry.ResourceVestackCenServiceRouteEntry(),
			//"vestack_cen_route_entry":                 cen_route_entry.ResourceVestackCenRouteEntry(),
			//
			//// ================ VPN ================
			//"vestack_vpn_gateway":       vpn_gateway.ResourceVestackVpnGateway(),
			//"vestack_customer_gateway":  customer_gateway.ResourceVestackCustomerGateway(),
			//"vestack_vpn_connection":    vpn_connection.ResourceVestackVpnConnection(),
			//"vestack_vpn_gateway_route": vpn_gateway_route.ResourceVestackVpnGatewayRoute(),

			// ================ VKE ================
			"vestack_vke_node":                           node.ResourceVestackVkeNode(),
			"vestack_vke_cluster":                        cluster.ResourceVestackVkeCluster(),
			"vestack_vke_node_pool":                      node_pool.ResourceVestackNodePool(),
			"vestack_vke_addon":                          addon.ResourceVestackVkeAddon(),
			"vestack_vke_default_node_pool":              default_node_pool.ResourceVestackDefaultNodePool(),
			"vestack_vke_default_node_pool_batch_attach": default_node_pool_batch_attach.ResourceVestackDefaultNodePoolBatchAttach(),
			"vestack_vke_kubeconfig":                     kubeconfig.ResourceVestackVkeKubeconfig(),

			// ================ IAM ================
			"vestack_iam_policy":                 iam_policy.ResourceVestackIamPolicy(),
			"vestack_iam_role":                   iam_role.ResourceVestackIamRole(),
			"vestack_iam_role_policy_attachment": iam_role_policy_attachment.ResourceVestackIamRolePolicyAttachment(),
			"vestack_iam_access_key":             iam_access_key.ResourceVestackIamAccessKey(),
			"vestack_iam_user":                   iam_user.ResourceVestackIamUser(),
			"vestack_iam_login_profile":          iam_login_profile.ResourceVestackIamLoginProfile(),
			"vestack_iam_user_policy_attachment": iam_user_policy_attachment.ResourceVestackIamUserPolicyAttachment(),
			"vestack_iam_service_linked_role":    iam_service_linked_role.ResourceVestackIamServiceLinkedRole(),

			// ================ RDS V1 ==============
			//"vestack_rds_instance":           rds_instance.ResourceVestackRdsInstance(),
			//"vestack_rds_database":           rds_database.ResourceVestackRdsDatabase(),
			//"vestack_rds_account":            rds_account.ResourceVestackRdsAccount(),
			//"vestack_rds_ip_list":            rds_ip_list.ResourceVestackRdsIpList(),
			//"vestack_rds_account_privilege":  rds_account_privilege.ResourceVestackRdsAccountPrivilege(),
			//"vestack_rds_parameter_template": rds_parameter_template.ResourceVestackRdsParameterTemplate(),
			//
			//// ================ ESCloud ================
			//"vestack_escloud_instance": instance.ResourceVestackESCloudInstance(),

			//================= TOS =================
			"vestack_tos_bucket":        bucket.ResourceVestackTosBucket(),
			"vestack_tos_object":        object.ResourceVestackTosObject(),
			"vestack_tos_bucket_policy": bucket_policy.ResourceVestackTosBucketPolicy(),

			// ================ PLB ================
			//"vestack_plb":                     plb.ResourceVestackPlb(),
			//"vestack_plb_vpc_associate":       plb_vpc_associate.ResourceVestackPlbVpcAssociate(),
			//"vestack_plb_listener":            plbListener.ResourceVestackListener(),
			//"vestack_plb_server_group":        plbServerGroup.ResourceVestackServerGroup(),
			//"vestack_plb_server_group_server": plbServerGroupServer.ResourceVestackServerGroupServer(),
			//"vestack_plb_acl":                 plbAcl.ResourceVestackAcl(),
			//"vestack_plb_acl_entry":           plbAclEntry.ResourceVestackAclEntry(),
			//
			//// ================ Redis ==============
			//"vestack_redis_allow_list":           redis_allow_list.ResourceVestackRedisAllowList(),
			//"vestack_redis_endpoint":             redis_endpoint.ResourceVestackRedisEndpoint(),
			//"vestack_redis_allow_list_associate": redis_allow_list_associate.ResourceVestackRedisAllowListAssociate(),
			//"vestack_redis_backup":               redis_backup.ResourceVestackRedisBackup(),
			//"vestack_redis_backup_restore":       redis_backup_restore.ResourceVestackRedisBackupRestore(),
			//"vestack_redis_account":              redisAccount.ResourceVestackRedisAccount(),
			//"vestack_redis_instance":             redisInstance.ResourceVestackRedisDbInstance(),
			//"vestack_redis_instance_state":       instance_state.ResourceVestackRedisInstanceState(),
			//"vestack_redis_continuous_backup":    redisContinuousBackup.ResourceVestackRedisContinuousBackup(),

			// ================ CR ================
			"vestack_cr_registry":       cr_registry.ResourceVestackCrRegistry(),
			"vestack_cr_registry_state": cr_registry_state.ResourceVestackCrRegistryState(),
			"vestack_cr_namespace":      cr_namespace.ResourceVestackCrNamespace(),
			"vestack_cr_repository":     cr_repository.ResourceVestackCrRepository(),
			"vestack_cr_tag":            cr_tag.ResourceVestackCrTag(),
			"vestack_cr_endpoint":       cr_endpoint.ResourceVestackCrEndpoint(),
			"vestack_cr_vpc_endpoint":   cr_vpc_endpoint.ResourceVestackCrVpcEndpoint(),

			// ================ Shuttle ================
			//"vestack_shuttle":                shuttle.ResourceVestackShuttle(),
			//"vestack_shuttle_client":         shuttle_client.ResourceVestackShuttleClient(),
			//"vestack_shuttle_server":         shuttle_server.ResourceVestackShuttleServer(),
			//"vestack_shuttle_server_v1":      shuttle_server_v1.ResourceVestackShuttleServerV1(),
			//"vestack_shuttle_client_v1":      shuttle_client_v1.ResourceVestackShuttleClientV1(),
			//"vestack_shuttle_association_v1": shuttle_association_v1.ResourceVestackShuttleAssociationV1(),
			//
			//// ================ Veenedge ================
			//"vestack_veenedge_cloud_server": cloud_server.ResourceVestackCloudServer(),
			//"vestack_veenedge_instance":     veInstance.ResourceVestackInstance(),
			//"vestack_veenedge_vpc":          veVpc.ResourceVestackVpc(),
			//
			//// ================ Kafka ================
			//"vestack_kafka_sasl_user":      kafka_sasl_user.ResourceVestackKafkaSaslUser(),
			//"vestack_kafka_group":          kafka_group.ResourceVestackKafkaGroup(),
			//"vestack_kafka_topic":          kafka_topic.ResourceVestackKafkaTopic(),
			//"vestack_kafka_instance":       kafka_instance.ResourceVestackKafkaInstance(),
			//"vestack_kafka_public_address": kafka_public_address.ResourceVestackKafkaPublicAddress(),
			//
			//// ================ MongoDB ================
			//"vestack_mongodb_instance":             mongodbInstance.ResourceVestackMongoDBInstance(),
			//"vestack_mongodb_endpoint":             endpoint.ResourceVestackMongoDBEndpoint(),
			//"vestack_mongodb_allow_list":           allow_list.ResourceVestackMongoDBAllowList(),
			//"vestack_mongodb_instance_parameter":   instance_parameter.ResourceVestackMongoDBInstanceParameter(),
			//"vestack_mongodb_allow_list_associate": allow_list_associate.ResourceVestackMongodbAllowListAssociate(),
			//"vestack_mongodb_ssl_state":            ssl_state.ResourceVestackMongoDBSSLState(),
			//
			//// ================ Bioos ================
			//"vestack_bioos_cluster":      bioosCluster.ResourceVestackBioosCluster(),
			//"vestack_bioos_workspace":    workspace.ResourceVestackBioosWorkspace(),
			//"vestack_bioos_cluster_bind": cluster_bind.ResourceVestackBioosClusterBind(),
			//
			//// ================ Veenedge ================
			//"vestack_direct_connect_connection":        direct_connect_connection.ResourceVestackDirectConnectConnection(),
			//"vestack_direct_connect_gateway":           direct_connect_gateway.ResourceVestackDirectConnectGateway(),
			//"vestack_direct_connect_virtual_interface": direct_connect_virtual_interface.ResourceVestackDirectConnectVirtualInterface(),
			//"vestack_direct_connect_bgp_peer":          direct_connect_bgp_peer.ResourceVestackDirectConnectBgpPeer(),
			//"vestack_direct_connect_gateway_route":     direct_connect_gateway_route.ResourceVestackDirectConnectGatewayRoute(),
			//
			//// ================ TransitRouter =============
			//"vestack_transit_router":                         transit_router.ResourceVestackTransitRouter(),
			//"vestack_transit_router_vpc_attachment":          transit_router_vpc_attachment.ResourceVestackTransitRouterVpcAttachment(),
			//"vestack_transit_router_vpn_attachment":          transit_router_vpn_attachment.ResourceVestackTransitRouterVpnAttachment(),
			//"vestack_transit_router_route_table":             trTable.ResourceVestackTransitRouterRouteTable(),
			//"vestack_transit_router_route_entry":             trEntry.ResourceVestackTransitRouterRouteEntry(),
			//"vestack_transit_router_route_table_association": route_table_association.ResourceVestackTransitRouterRouteTableAssociation(),
			//"vestack_transit_router_route_table_propagation": route_table_propagation.ResourceVestackTransitRouterRouteTablePropagation(),
			//
			//// ================ PrivateLink ==================
			//"vestack_privatelink_vpc_endpoint":                    vpc_endpoint.ResourceVestackPrivatelinkVpcEndpoint(),
			//"vestack_privatelink_vpc_endpoint_service":            vpc_endpoint_service.ResourceVestackPrivatelinkVpcEndpointService(),
			//"vestack_privatelink_vpc_endpoint_service_resource":   vpc_endpoint_service_resource.ResourceVestackPrivatelinkVpcEndpointServiceResource(),
			//"vestack_privatelink_vpc_endpoint_service_permission": vpc_endpoint_service_permission.ResourceVestackPrivatelinkVpcEndpointServicePermission(),
			//"vestack_privatelink_security_group":                  plSecurityGroup.ResourceVestackPrivatelinkSecurityGroupService(),
			//"vestack_privatelink_vpc_endpoint_connection":         vpc_endpoint_connection.ResourceVestackPrivatelinkVpcEndpointConnectionService(),
			//"vestack_privatelink_vpc_endpoint_zone":               vpc_endpoint_zone.ResourceVestackPrivatelinkVpcEndpointZone(),
			//
			//// ================ RDS V2 ==============
			//"vestack_rds_instance_v2": rds_instance_v2.ResourceVestackRdsInstance(),
			//
			//// ================ RdsMysql ================
			//"vestack_rds_mysql_instance":               rds_mysql_instance.ResourceVestackRdsMysqlInstance(),
			//"vestack_rds_mysql_instance_readonly_node": rds_mysql_instance_readonly_node.ResourceVestackRdsMysqlInstanceReadonlyNode(),
			//"vestack_rds_mysql_account":                rds_mysql_account.ResourceVestackRdsMysqlAccount(),
			//"vestack_rds_mysql_database":               rds_mysql_database.ResourceVestackRdsMysqlDatabase(),
			//"vestack_rds_mysql_allowlist":              allowlist.ResourceVestackRdsMysqlAllowlist(),
			//"vestack_rds_mysql_allowlist_associate":    allowlist_associate.ResourceVestackRdsMysqlAllowlistAssociate(),
			//
			//// ================ TLS ================
			//"vestack_tls_kafka_consumer":     kafka_consumer.ResourceVestackTlsKafkaConsumer(),
			//"vestack_tls_host_group":         host_group.ResourceVestackTlsHostGroup(),
			//"vestack_tls_rule":               tlsRule.ResourceVestackTlsRule(),
			//"vestack_tls_rule_applier":       rule_applier.ResourceVestackTlsRuleApplier(),
			//"vestack_tls_alarm":              alarm.ResourceVestackTlsAlarm(),
			//"vestack_tls_alarm_notify_group": alarm_notify_group.ResourceVestackTlsAlarmNotifyGroup(),
			//"vestack_tls_host":               host.ResourceVestackTlsHost(),
			//"vestack_tls_project":            tlsProject.ResourceVestackTlsProject(),
			//"vestack_tls_topic":              tlsTopic.ResourceVestackTlsTopic(),
			//"vestack_tls_index":              tlsIndex.ResourceVestackTlsIndex(),
			//
			//// ================ VMP ================
			//"vestack_vmp_workspace":           vmp_workspace.ResourceVestackVmpWorkspace(),
			//"vestack_vmp_rule_file":           vmp_rule_file.ResourceVestackVmpRuleFile(),
			//"vestack_vmp_contact_group":       vmp_contact_group.ResourceVestackVmpContactGroup(),
			//"vestack_vmp_contact":             vmp_contact.ResourceVestackVmpContact(),
			//"vestack_vmp_alerting_rule":       vmp_alerting_rule.ResourceVestackVmpAlertingRule(),
			//"vestack_vmp_notify_group_policy": vmp_notify_group_policy.ResourceVestackVmpNotifyGroupPolicy(),
			//"vestack_vmp_notify_policy":       vmp_notify_policy.ResourceVestackVmpNotifyPolicy(),
			//
			//// ================ InternetTunnel ================
			//"vestack_internet_tunnel_gateway":           internet_tunnel_gateway.ResourceVestackInternetTunnelGateway(),
			//"vestack_internet_tunnel_virtual_interface": internet_tunnel_virtual_interface.ResourceVestackInternetTunnelVirtualInterface(),
			//"vestack_internet_tunnel_bgp_peer":          internet_tunnel_bgp_peer.ResourceVestackInternetTunnelBgpPeer(),
			//"vestack_internet_tunnel_bandwidth":         internet_tunnel_bandwidth.ResourceVestackInternetTunnelBandwidth(),
			//"vestack_internet_tunnel_address":           internet_tunnel_address.ResourceVestackInternetTunnelAddress(),
			//"vestack_internet_tunnel_address_associate": internet_tunnel_address_associate.ResourceVestackInternetTunnelAddressAssociate(),
			//
			//// ================ FastTrack ================
			//"vestack_fast_track_associate": fast_track_associate.ResourceVestackFastTrackAssociate(),
			//"vestack_fast_track":           fast_track.ResourceVestackFastTrack(),
			//
			//// ================ NAS ================
			//"vestack_nas_file_system":      nas_file_system.ResourceVestackNasFileSystem(),
			//"vestack_nas_snapshot":         nas_snapshot.ResourceVestackNasSnapshot(),
			//"vestack_nas_mount_point":      nas_mount_point.ResourceVestackNasMountPoint(),
			//"vestack_nas_permission_group": nas_permission_group.ResourceVestackNasPermissionGroup(),
			//
			//// ================ Cloudfs ================
			//"vestack_cloudfs_file_system": cloudfs_file_system.ResourceVestackCloudfsFileSystem(),
			//"vestack_cloudfs_access":      cloudfs_access.ResourceVestackCloudfsAccess(),
			//"vestack_cloudfs_namespace":   cloudfs_namespace.ResourceVestackCloudfsNamespace(),
		},
		ConfigureFunc: ProviderConfigure,
	}
}

func ProviderConfigure(d *schema.ResourceData) (interface{}, error) {
	config := ve.Config{
		AccessKey:         d.Get("access_key").(string),
		SecretKey:         d.Get("secret_key").(string),
		SessionToken:      d.Get("session_token").(string),
		Region:            d.Get("region").(string),
		Endpoint:          d.Get("endpoint").(string),
		DisableSSL:        d.Get("disable_ssl").(bool),
		CustomerHeaders:   map[string]string{},
		CustomerEndpoints: defaultCustomerEndPoints(),
		ProxyUrl:          d.Get("proxy_url").(string),
	}
	logger.Info("access_key: %+v", config.AccessKey)
	headers := d.Get("customer_headers").(string)
	if headers != "" {
		hs1 := strings.Split(headers, ",")
		for _, hh := range hs1 {
			hs2 := strings.Split(hh, ":")
			if len(hs2) == 2 {
				config.CustomerHeaders[hs2[0]] = hs2[1]
			}
		}
	}

	endpoints := d.Get("customer_endpoints").(string)
	if endpoints != "" {
		ends := strings.Split(endpoints, ",")
		for _, end := range ends {
			point := strings.Split(end, ":")
			if len(point) == 2 {
				config.CustomerEndpoints[point[0]] = point[1]
			}
		}
	}

	client, err := config.Client()
	return client, err
}

func defaultCustomerEndPoints() map[string]string {
	return map[string]string{
		"veenedge": "veenedge.volcengineapi.com",
	}
}
