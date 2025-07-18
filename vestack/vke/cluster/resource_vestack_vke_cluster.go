package cluster

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
VkeCluster can be imported using the id, e.g.
```
$ terraform import vestack_vke_cluster.default cc9l74mvqtofjnoj5****
```

*/

func ResourceVestackVkeCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackVkeClusterCreate,
		Read:   resourceVestackVkeClusterRead,
		Update: resourceVestackVkeClusterUpdate,
		Delete: resourceVestackVkeClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"client_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ClientToken is a case-sensitive string of no more than 64 ASCII characters passed in by the caller.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the cluster.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the Cluster.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the cluster.",
			},
			"delete_protection_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The delete protection of the cluster, the value is `true` or `false`.",
			},
			"kubernetes_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if k == "kubernetes_version" && strings.Contains(old, new) {
						return true
					}
					return false
				},
				Description: "The version of Kubernetes specified when creating a VKE cluster (specified to patch version), if not specified, the latest Kubernetes version supported by VKE is used by default, which is a 3-segment version format starting with a lowercase v, that is, KubernetesVersion with IsLatestVersion=True in the return value of ListSupportedVersions.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project name of the cluster.",
			},
			"tags": bp.TagsSchema(),
			"control_plane_nodes_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "The control plane node information for the VKE cluster instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Node resource provider name, available values: VeStack: Resources built on veStack full-stack version.",
						},
						"ve_stack": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The resources in veStack are used for the master node in the VKE cluster.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"new_node_configs": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Configuration for auto create new nodes.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"count": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "numbers of master, must be 1 3 5 7.",
												},
												"subnet_ids": {
													Type:     schema.TypeSet,
													Required: true,
													ForceNew: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Set:         schema.HashString,
													Description: "The subnet ID for the master node.",
												},
												"instance_type_id": {
													Type:     schema.TypeString,
													Required: true,
												},
												"system_volume": {
													Type:     schema.TypeList,
													Required: true,
													MaxItems: 1,
													MinItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:         schema.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringInSlice([]string{"ESSD_PL0", "ESSD_FlexPL"}, false),
																Description:  "The Type of SystemVolume.",
															},
															"size": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Disk size, unit GB, value range is 40~2048, default value is 40.",
															},
														},
													},
													Description: "The SystemVolume of NodeConfig.",
												},
												"data_volumes": {
													Type:     schema.TypeList,
													Optional: true,
													ForceNew: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:         schema.TypeString,
																Optional:     true,
																ForceNew:     true,
																ValidateFunc: validation.StringInSlice([]string{"ESSD_FlexPL", "ESSD_PL0"}, false),
																Description:  "The Type of DataVolumes, the value can be `ESSD_PL0` or `ESSD_FlexPL`.",
															},
															"size": {
																Type:         schema.TypeInt,
																Optional:     true,
																ForceNew:     true,
																ValidateFunc: validation.IntBetween(20, 32768),
																Description:  "Disk size, unit GB, value range is 20~32768, default value is 20.",
															},
															"mount_point": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The target mounting directory after disk formatting.",
															},
														},
													},
													Description: "The DataVolumes of NodeConfig.",
												},
												"initialize_script": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"security": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"security_group_ids": {
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Set:         schema.HashString,
																Description: "The security group id.",
															},
															"login": {
																Type:     schema.TypeList,
																Required: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"password": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									"existed_node_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Use an existing node as the cluster master node configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instances": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Node resources, ECS instance ID list.",
															},
														},
													},
													Description: "ECS node information list.",
												},
												"keep_instance_name": {
													Type:     schema.TypeBool,
													Optional: true,
													Description: "Keep the node name to join the cluster, the priority is higher than NamePrefix, the value is:" +
														"false: (Default) Do not maintain node names." +
														"true: Keep the node name as the original host instance name.",
												},
												"name_prefix": {
													Type:     schema.TypeString,
													Optional: true,
													Description: "The node naming prefix has a lower priority than KeepInstanceName. When the value is empty, it means that the node naming prefix is not enabled. Among them, the prefix verification rules:" +
														"Supports English letters, numbers and dashes -, dashes - cannot be used continuously." +
														"It can only start with an English letter and end with an English letter or number." +
														"The length is 2 to 51 characters.",
												},
												"additional_container_storage_enabled": {
													Type:     schema.TypeBool,
													Optional: true,
													Description: "Select the data disk of the configuration node and format it and mount it as the storage directory for container images and logs. The value is:" +
														"false: (default) off." +
														"true: enable.",
												},
												"container_storage_path": {
													Type:     schema.TypeString,
													Optional: true,
													Description: "Use this data disk device to mount the container and image storage directory /var/lib/containerd. It is only valid when AdditionalContainerStorageEnabled=true and cannot be empty." +
														"The following conditions must be met, otherwise the initialization will fail:" +
														"Only cloud server instances with mounted data disks are supported." +
														"When specifying the data disk device name, please ensure that the data disk device exists, and the problem will be automatically initialized." +
														"When specifying a data disk partition or logical volume name, make sure that the partition or logical volume exists and is an exct4 file system." +
														"Notice" +
														"When specifying a data disk device, it will be automatically formatted and mounted directly. Please be sure to back up the data in advance." +
														"When specifying a data disk partition or logical volume name, no formatting is required.",
												},
												"initialize_script": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cript that is executed after ECS nodes are created and Kubernetes components are deployed. Supports Shell format, the length after Base64 encoding does not exceed 16 KB.",
												},
												"security": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"security_group_ids": {
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Set: schema.HashString,
																Description: "List of security group IDs in which the node network is located." +
																	"Call the DescribeSecurityGroups interface of the private network to obtain the security group ID." +
																	"" +
																	"Notice" +
																	"Must be in the same private network as the cluster." +
																	"When the value is empty, the default security group of the cluster node is used by default (the naming format is <cluster ID>-common)." +
																	"A single node pool supports up to 5 security groups (including the default security group of cluster nodes).",
															},
															"login": {
																Type:     schema.TypeList,
																Required: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"password": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
																Description: "Node access mode configuration." +
																	"Support password mode or key pair mode. When they are passed in at the same time, the key pair will be used first.",
															},
														},
													},
													Description: "Node security configuration.",
												},
											},
										},
									},
									"deployment_set_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Deployment set ID. If specified, the master node will be added to the deployment set group. Currently, only the new node method is supported.",
									},
								},
							},
						},
					},
				},
			},
			"cluster_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "The config of the cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_ids": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
							Description: "The subnet ID for the cluster control plane to communicate within the private network.\n" +
								"Up to 3 subnets can be selected from each available zone, and a maximum of 2 subnets can be added to each available zone.\n" +
								"Cannot support deleting configured subnets.",
						},
						"api_server_public_access_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Cluster API Server public network access configuration, the value is `true` or `false`.",
						},
						"api_server_public_access_config": {
							Type:             schema.TypeList,
							MaxItems:         1,
							Optional:         true,
							DiffSuppressFunc: ApiServerPublicAccessConfigFieldDiffSuppress,
							Description:      "Cluster API Server public network access configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"public_access_network_config": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										ForceNew:    true,
										Description: "Public network access network configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"billing_type": {
													Type:         schema.TypeString,
													Optional:     true,
													Description:  "Billing type of public IP, the value is `PostPaidByBandwidth` or `PostPaidByTraffic`.",
													ValidateFunc: validation.StringInSlice([]string{"PostPaidByBandwidth", "PostPaidByTraffic"}, false),
												},
												"bandwidth": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The peak bandwidth of the public IP, unit: Mbps.",
												},
												"isp": {
													Type:     schema.TypeString,
													Optional: true,
													Description: "Line type of public network IP, value:" +
														"BGP: (Default) BGP circuit." +
														"ChinaMobile: China Mobile." +
														"ChinaUnicom: China Unicom." +
														"ChinaTelecom: China Telecom.",
													ValidateFunc: validation.StringInSlice([]string{"BGP", "ChinaMobile", "ChinaUnicom", "ChinaTelecom"}, false),
												},
											},
										},
									},
								},
							},
						},
						"resource_public_access_default_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "Node public network access configuration, the value is `true` or `false`.",
						},
						"ip_family": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "[SkipDoc]The IpFamily configuration,the value is `Ipv4` or `DualStack`.",
						},
					},
				},
			},
			"pods_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "The config of the pods.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pod_network_mode": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							Description: "The container network model of the cluster, the value is `Flannel` or `VpcCniShared` or `VpcCniHybrid` or `CalicoVxlan` or `CalicoBgp`. " +
								"Flannel: Flannel network model, an independent Underlay container network solution, combined with the global routing capability of VPC, to achieve a high-performance network experience for the cluster. " +
								"VpcCniShared: VPC-CNI network model, an Underlay container network solution based on the ENI of the private network elastic network card, with high network communication performance. " +
								"CalicoVxlan: Calico network Vxlan mode, an overlay container network solution independent of the control plane. " +
								"CalicoBgp: Calico network BGP mode, configure BGP between nodes or peer network infrastructure to distribute routing information (OnPremise cluster supported only).",
						},
						"flannel_config": {
							Type:             schema.TypeList,
							MaxItems:         1,
							ForceNew:         true,
							Optional:         true,
							Description:      "Flannel network configuration.",
							DiffSuppressFunc: FlannelFieldDiffSuppress,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pod_cidrs": {
										Type:     schema.TypeSet,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Set:         schema.HashString,
										Description: "Pod CIDR for the Flannel container network.",
									},
									"max_pods_per_node": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Description: "The maximum number of single-node Pod instances for a Flannel container network, the value can be `16` or `32` or `64` or `128` or `256`.",
									},
								},
							},
						},
						"vpc_cni_config": {
							Type:             schema.TypeList,
							MaxItems:         1,
							Optional:         true,
							Description:      "VPC-CNI network configuration.",
							DiffSuppressFunc: VpcCniConfigFieldDiffSuppress,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											return true
										},
										Description: "The private network where the cluster control plane network resides.",
									},
									"subnet_ids": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Set:         schema.HashString,
										Description: "A list of Pod subnet IDs for the VPC-CNI container network.",
									},
								},
							},
						},
						"calico_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Calico network configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pod_cidrs": {
										Type:     schema.TypeSet,
										Required: true,
										ForceNew: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Set:         schema.HashString,
										Description: "Pod CIDR for the Flannel container network.",
									},
									"max_pods_per_node": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Description: "The maximum number of single-node Pod instances for a Flannel container network, the value can be `16` or `32` or `64` or `128` or `256`, default value is `64`.",
									},
									"bgp_config": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mode": {
													Type:     schema.TypeString,
													Required: true,
													ForceNew: true,
													Description: "BGP mode, optional values are: FullMesh | RouteReflectors: " +
														"- FullMesh: all nodes serve as RR peers. " +
														"- RouteReflectors: The Master node as the RR contains the routing information of all nodes, and the node only has routes pointing to the RR node.",
												},
												"as_number": {
													Type:        schema.TypeInt,
													Required:    true,
													ForceNew:    true,
													Description: "Value range [64512, 65534].",
												},
												"external_route_reflector_enabled": {
													Type:     schema.TypeBool,
													Optional: true,
													ForceNew: true,
													Description: "Whether to enable external RoutReflector, optional values true | false: " +
														"- true: Enable external RoutReflector instead of using Master node as RoutReflector. " +
														"- false: Use Master node as RoutReflector. " +
														"Default value: false.",
												},
												"route_reflector_peer_points": {
													Type:     schema.TypeSet,
													Required: true,
													ForceNew: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ip_address": {
																Type:        schema.TypeString,
																Required:    true,
																ForceNew:    true,
																Description: "IP address of RouteReflector.",
															},
															"port": {
																Type:        schema.TypeInt,
																Required:    true,
																ForceNew:    true,
																Description: "Port of RouteReflector.",
															},
															"as_number": {
																Type:        schema.TypeInt,
																Optional:    true,
																ForceNew:    true,
																Description: "If ExternalRouteReflectorEnabled=true, this parameter is optional, otherwise this parameter cannot be empty [64512, 65534].",
															},
														},
													},
													Description: "Pod CIDR for the Flannel container network.",
												},
											},
										},
										Description: "Configuration information of BGP mode under Calico. " +
											"Only supported in Onpremise cluster & CalicoBgp mode.",
									},
								},
							},
						},
					},
				},
			},
			"services_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: "The config of the services.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_cidrsv4": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set:         schema.HashString,
							Description: "The IPv4 private network address exposed by the service.",
						},
					},
				},
			},
			"kubeconfig_public": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kubeconfig data with public network access, returned in BASE64 encoding, it is suggested to use vke_kubeconfig instead.",
			},
			"kubeconfig_private": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kubeconfig data with private network access, returned in BASE64 encoding, it is suggested to use vke_kubeconfig instead.",
			},
			"eip_allocation_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Eip allocation Id.",
			},
			"logging_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Cluster log configuration information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_project_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The TLS log item ID of the collection target.",
						},
						"log_setups": {
							Type:        schema.TypeSet,
							Optional:    true,
							Set:         logSetupsHash,
							Description: "Cluster logging options. This structure can only be modified and added, and cannot be deleted. When encountering a `cannot be deleted` error, please query the log setups of the current cluster and fill in the current `tf` file.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"log_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The currently enabled log type.",
									},
									"log_ttl": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      30,
										ValidateFunc: validation.IntBetween(1, 3650),
										Description:  "The storage time of logs in Log Service. After the specified log storage time is exceeded, the expired logs in this log topic will be automatically cleared. The unit is days, and the default is 30 days. The value range is 1 to 3650, specifying 3650 days means permanent storage.",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether to enable the log option, true means enable, false means not enable, the default is false. When Enabled is changed from false to true, a new Topic will be created.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceVestackVkeClusterCreate(d *schema.ResourceData, meta interface{}) (err error) {
	clusterService := NewVkeClusterService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(clusterService, d, ResourceVestackVkeCluster())
	if err != nil {
		return fmt.Errorf("error on creating cluster  %q, %w", d.Id(), err)
	}
	return resourceVestackVkeClusterRead(d, meta)
}

func resourceVestackVkeClusterRead(d *schema.ResourceData, meta interface{}) (err error) {
	clusterService := NewVkeClusterService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(clusterService, d, ResourceVestackVkeCluster())
	if err != nil {
		return fmt.Errorf("error on reading cluster %q, %w", d.Id(), err)
	}
	return err
}

func resourceVestackVkeClusterUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	clusterService := NewVkeClusterService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(clusterService, d, ResourceVestackVkeCluster())
	if err != nil {
		return fmt.Errorf("error on updating cluster  %q, %w", d.Id(), err)
	}
	return resourceVestackVkeClusterRead(d, meta)
}

func resourceVestackVkeClusterDelete(d *schema.ResourceData, meta interface{}) (err error) {
	clusterService := NewVkeClusterService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(clusterService, d, ResourceVestackVkeCluster())
	if err != nil {
		return fmt.Errorf("error on deleting cluster %q, %w", d.Id(), err)
	}
	return err
}

func logSetupsHash(i interface{}) int {
	if i == nil {
		return hashcode.String("")
	}
	m := i.(map[string]interface{})
	var (
		buf bytes.Buffer
	)
	buf.WriteString(fmt.Sprintf("%v#%v#%v", m["log_type"], m["log_ttl"], m["enabled"]))
	return hashcode.String(buf.String())
}
