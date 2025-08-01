package cluster

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

func DataSourceVestackVkeVkeClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVestackVkeClustersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Cluster IDs.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of Cluster.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of Cluster query.",
			},
			"page_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The page number of clusters query.",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The page size of clusters query.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the cluster.",
			},
			"delete_protection_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The delete protection of the cluster, the value is `true` or `false`.",
			},
			"pods_config_pod_network_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The container network model of the cluster, the value is `Flannel` or `VpcCniShared`. Flannel: Flannel network model, an independent Underlay container network solution, combined with the global routing capability of VPC, to achieve a high-performance network experience for the cluster. VpcCniShared: VPC-CNI network model, an Underlay container network solution based on the ENI of the private network elastic network card, with high network communication performance.",
			},
			"statuses": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Array of cluster states to filter. (The elements of the array are logically ORed. A maximum of 15 state array elements can be filled at a time).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"phase": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The status of cluster. the value contains `Creating`, `Running`, `Updating`, `Deleting`, `Stopped`, `Failed`.",
						},
						"conditions_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The state condition in the current main state of the cluster, that is, the reason for entering the main state, there can be multiple reasons, the value contains `Progressing`, `Ok`, `Degraded`, `SetByProvider`, `Balance`, `Security`, `CreateError`, `ResourceCleanupFailed`, `LimitedByQuota`, `StockOut`,`Unknown`.",
						},
					},
				},
			},
			"create_client_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ClientToken when the cluster is created successfully. ClientToken is a string that guarantees the idempotency of the request. This string is passed in by the caller.",
			},
			"update_client_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ClientToken when the last cluster update succeeded. ClientToken is a string that guarantees the idempotency of the request. This string is passed in by the caller.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The project name of the cluster.",
			},
			"tags": bp.TagsSchema(),
			"clusters": {
				Description: "The collection of VkeCluster query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true, // tf中不支持写值
							Description: "The ID of the Cluster.",
						},
						"create_client_token": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ClientToken on successful creation. ClientToken is a string that guarantees the idempotency of the request. This string is passed in by the caller.",
						},
						"update_client_token": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ClientToken when the last update was successful. ClientToken is a string that guarantees the idempotency of the request. This string is passed in by the caller.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster creation time. UTC+0 time in standard RFC3339 format.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last time a request was accepted by the cluster and executed or completed. UTC+0 time in standard RFC3339 format.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the cluster.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the cluster.",
						},
						"delete_protection_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The delete protection of the cluster, the value is `true` or `false`.",
						},
						"kubernetes_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Kubernetes version information corresponding to the cluster, specific to the patch version.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project name of the cluster.",
						},
						"status": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Computed:    true,
							Description: "The status of the cluster.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"phase": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The status of cluster. the value contains `Creating`, `Running`, `Updating`, `Deleting`, `Stopped`, `Failed`.",
									},
									"conditions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The state condition in the current primary state of the cluster, that is, the reason for entering the primary state.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The state condition in the current main state of the cluster, that is, the reason for entering the main state, there can be multiple reasons, the value contains `Progressing`, `Ok`, `Balance`, `CreateError`, `ResourceCleanupFailed`, `Unknown`.",
												},
											},
										},
									},
								},
							},
						},
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
																			Computed:    false,
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
							Computed:    true,
							Description: "The config of the cluster.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the private network (VPC) where the network of the cluster control plane and some nodes is located.",
									},
									"subnet_ids": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Set:         schema.HashString,
										Description: "The subnet ID for the cluster control plane to communicate within the private network.",
									},
									"security_group_ids": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Set:         schema.HashString,
										Description: "The security group used by the cluster control plane and nodes.",
									},
									"api_server_public_access_enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Cluster API Server public network access configuration, the value is `true` or `false`.",
									},
									"api_server_public_access_config": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Computed:    true,
										Description: "Cluster API Server public network access configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"public_access_network_config": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Computed:    true,
													Description: "Public network access network configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"billing_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Billing type of public IP, the value is `PostPaidByBandwidth` or `PostPaidByTraffic`.",
															},
															"bandwidth": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The peak bandwidth of the public IP, unit: Mbps.",
															},
															"isp": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The ISP of public IP.",
															},
														},
													},
												},
												"access_source_ipsv4": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Set:         schema.HashString,
													Description: "IPv4 public network access whitelist. A null value means all network segments (0.0.0.0/0) are allowed to pass.",
												},
												"ip_family": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "[SkipDoc]The IpFamily configuration,the value is `Ipv4` or `DualStack`.",
												},
											},
										},
									},
									"resource_public_access_default_enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Node public network access configuration, the value is `true` or `false`.",
									},
									"api_server_endpoints": {
										Type:        schema.TypeList,
										Computed:    true,
										MaxItems:    1,
										Description: "Endpoint information accessed by the cluster API Server.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"private_ip": {
													Type:        schema.TypeList,
													Computed:    true,
													MaxItems:    1,
													Description: "Endpoint address of the cluster API Server private network.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Ipv4 address.",
															},
														},
													},
												},
												"public_ip": {
													Type:        schema.TypeList,
													Computed:    true,
													MaxItems:    1,
													Description: "Endpoint address of the cluster API Server public network.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Ipv4 address.",
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
						"pods_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Computed:    true,
							Description: "The config of the pods.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pod_network_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Container Pod Network Type (CNI), the value is `Flannel` or `VpcCniShared`.",
									},
									"flannel_config": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Computed:    true,
										Description: "Flannel network configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"pod_cidrs": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Set:         schema.HashString,
													Description: "Pod CIDR for the Flannel container network.",
												},
												"max_pods_per_node": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The maximum number of single-node Pod instances for a Flannel container network.",
												},
											},
										},
									},
									"vpc_cni_config": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Computed:    true,
										Description: "VPC-CNI network configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The private network where the cluster control plane network resides.",
												},
												"subnet_ids": {
													Type:     schema.TypeSet,
													Computed: true,
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
										Computed:    true,
										Description: "Calico network configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"pod_cidrs": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Set:         schema.HashString,
													Description: "Pod CIDR for the Flannel container network.",
												},
												"max_pods_per_node": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The maximum number of single-node Pod instances for a Flannel container network, the value can be `16` or `32` or `64` or `128` or `256`, default value is `64`.",
												},
												"bgp_config": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": {
																Type:     schema.TypeString,
																Computed: true,
																Description: "BGP mode, optional values are: FullMesh | RouteReflectors: " +
																	"- FullMesh: all nodes serve as RR peers. " +
																	"- RouteReflectors: The Master node as the RR contains the routing information of all nodes, and the node only has routes pointing to the RR node.",
															},
															"as_number": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Value range [64512, 65534].",
															},
															"external_route_reflector_enabled": {
																Type:     schema.TypeBool,
																Computed: true,
																Description: "Whether to enable external RoutReflector, optional values true | false: " +
																	"- true: Enable external RoutReflector instead of using Master node as RoutReflector. " +
																	"- false: Use Master node as RoutReflector. " +
																	"Default value: false.",
															},
															"route_reflector_peer_points": {
																Type:     schema.TypeSet,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ip_address": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "IP address of RouteReflector.",
																		},
																		"port": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Port of RouteReflector.",
																		},
																		"as_number": {
																			Type:        schema.TypeInt,
																			Computed:    true,
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
							Computed:    true,
							Description: "The config of the services.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_cidrsv4": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Set:         schema.HashString,
										Description: "The IPv4 private network address exposed by the service.",
									},
								},
							},
						},
						"node_statistics": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Computed:    true,
							Description: "Statistics on the number of nodes corresponding to each master state in the cluster.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total number of nodes.",
									},
									"creating_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Phase=Creating total number of nodes.",
									},
									"running_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Phase=Running total number of nodes.",
									},
									"stopped_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Deprecated:  "This field has been deprecated and is not recommended for use.",
										Description: "Phase=Stopped total number of nodes.",
									},
									"updating_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Phase=Updating total number of nodes.",
									},
									"deleting_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Phase=Deleting total number of nodes.",
									},
									"failed_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Phase=Failed total number of nodes.",
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
						"tags": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Tags of the Cluster.",
							Set:         bp.VkeTagsResponseHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The Key of Tags.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The Value of Tags.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The Type of Tags.",
									},
								},
							},
						},
						"logging_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Cluster log configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"log_project_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The TLS log item ID of the collection target.",
									},
									"log_setups": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Cluster logging options.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"log_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The currently enabled log type.",
												},
												"log_ttl": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The storage time of logs in Log Service. After the specified log storage time is exceeded, the expired logs in this log topic will be automatically cleared. The unit is days, and the default is 30 days. The value range is 1 to 3650, specifying 3650 days means permanent storage.",
												},
												"enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to enable the log option, true means enable, false means not enable, the default is false. When Enabled is changed from false to true, a new Topic will be created.",
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
		},
	}
}

func dataSourceVestackVkeClustersRead(d *schema.ResourceData, meta interface{}) error {
	clusterService := NewVkeClusterService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(clusterService, d, DataSourceVestackVkeVkeClusters())
}
