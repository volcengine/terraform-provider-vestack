package node_pool

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
NodePool can be imported using the id, e.g.
```
$ terraform import vestack_vke_node_pool.default pcabe57vqtofgrbln3dp0
```

*/

func ResourceVestackNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackNodePoolCreate,
		Read:   resourceVestackNodePoolRead,
		Update: resourceVestackNodePoolUpdate,
		Delete: resourceVestackNodePoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Name of NodePool.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ClusterId of NodePool.",
			},
			"client_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ClientToken of NodePool.",
			},
			"tags": bp.TagsSchema(),
			"instance_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 100,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"auto_scaling"},
				Description: "The list of existing ECS instance ids. Add existing instances with same type of security group under the same cluster VPC to the custom node pool.\n" +
					"Note that removing instance ids from the list will only remove the nodes from cluster and not release the ECS instances. But deleting node pool will release the ECS instances in it.\n" +
					"It is not recommended to use this field, it is recommended to use `volcengine_vke_node` resource to add an existing instance to a custom node pool.",
			},
			"keep_instance_name": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Whether to keep instance name when adding an existing instance to a custom node pool, the value is `true` or `false`.\n" +
					"This field is valid only when adding new instances to the custom node pool.",
			},
			"auto_scaling": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_ids"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Is Enabled of AutoScaling.",
						},
						"max_replicas": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      10,
							ValidateFunc: validation.IntBetween(1, 2000),
							Description:  "The MaxReplicas of AutoScaling, default 10, range in 1~2000.",
						},
						"min_replicas": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The MinReplicas of AutoScaling, default 0.",
						},
						"desired_replicas": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The DesiredReplicas of AutoScaling, default 0, range in min_replicas to max_replicas.",
						},
						"priority": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 100),
							Description:  "The Priority of AutoScaling, default 10, rang in 0~100.",
						},
						"subnet_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ZoneBalance",
								"Priority",
							}, false),
							Description: "Multi-subnet scheduling strategy for nodes. The value can be `ZoneBalance` or `Priority`.",
						},
					},
				},
				Description: "The node pool elastic scaling configuration information.",
			},
			"node_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type_ids": {
							Type:     schema.TypeList,
							Required: true,
							//ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The InstanceTypeIds of NodeConfig.",
						},
						"subnet_ids": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The SubnetIds of NodeConfig.",
						},
						"security": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_group_ids": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "The SecurityGroupIds of Security.",
									},
									"security_strategies": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "The SecurityStrategies of Security, the value can be empty or `Hids`.",
									},
									"login": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"password": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The Password of Security, this field must be encoded with base64.",
												},
												"ssh_key_pair_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The SshKeyPairName of Security.",
												},
											},
										},
										Description: "The Login of Security.",
									},
								},
							},
							Description: "The Security of NodeConfig.",
						},
						"system_volume": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The Type of SystemVolume, the value can be `PTSSD` or `ESSD_PL0` or `ESSD_FlexPL`.",
									},
									"size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The Size of SystemVolume, the value range in 20~2048.",
									},
								},
							},
							Description: "The SystemVolume of NodeConfig.",
						},
						"data_volumes": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "ESSD_PL0",
										Description: "The Type of DataVolumes, the value can be `PTSSD` or `ESSD_PL0` or `ESSD_FlexPL`. Default value is `ESSD_PL0`.",
									},
									"size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     20,
										Description: "The Size of DataVolumes, the value range in 20~32768. Default value is `20`.",
									},
									"mount_point": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The target mount directory of the disk. Must start with `/`.",
									},
								},
							},
							Description: "The DataVolumes of NodeConfig.",
						},
						"initialize_script": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The initializeScript of NodeConfig.",
						},
						"additional_container_storage_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "The AdditionalContainerStorageEnabled of NodeConfig.",
						},
						"image_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The ImageId of NodeConfig.",
						},
						"instance_charge_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "PostPaid",
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
							Description:  "The InstanceChargeType of PrePaid instance of NodeConfig. Valid values: PostPaid, PrePaid. Default value: PostPaid.",
						},
						"period": {
							Type:             schema.TypeInt,
							Optional:         true,
							Computed:         true,
							ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
							DiffSuppressFunc: prePaidDiffSuppressFunc,
							Description:      "The Period of PrePaid instance of NodeConfig. Valid values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36. Unit: month. when InstanceChargeType is PrePaid, default value is 12.",
						},
						"auto_renew": {
							Type:             schema.TypeBool,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: prePaidDiffSuppressFunc,
							Description:      "Is AutoRenew of PrePaid instance of NodeConfig. Valid values: true, false. when InstanceChargeType is PrePaid, default value is true.",
						},
						"auto_renew_period": {
							Type:             schema.TypeInt,
							Optional:         true,
							Computed:         true,
							ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12}),
							DiffSuppressFunc: prePaidAndAutoNewDiffSuppressFunc,
							Description:      "The AutoRenewPeriod of PrePaid instance of NodeConfig. Valid values: 1, 2, 3, 6, 12. Unit: month. when InstanceChargeType is PrePaid and AutoRenew enable, default value is 1.",
						},
						"name_prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The NamePrefix of NodeConfig.",
						},
						"ecs_tags": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Tags for Ecs.",
							Set:         bp.TagsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The Key of Tags.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The Value of Tags.",
									},
								},
							},
						},
						"hpc_cluster_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The IDs of HpcCluster, only one ID is supported currently.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},
						"project_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The project name of the ecs instance.",
						},
					},
				},
				Description: "The Config of NodePool.",
			},
			"kubernetes_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"labels": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The Key of Labels.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The Value of Labels.",
									},
								},
							},
							Set:         kubernetesConfigLabelHash,
							Description: "The Labels of KubernetesConfig.",
						},
						"taints": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The Key of Taints.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The Value of Taints.",
									},
									"effect": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "NoSchedule",
										Description: "The Effect of Taints, the value can be `NoSchedule` or `NoExecute` or `PreferNoSchedule`.",
									},
								},
							},
							Description: "The Taints of KubernetesConfig.",
						},
						"cordon": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "The Cordon of KubernetesConfig.",
						},
						"name_prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The NamePrefix of node metadata.",
						},
						"auto_sync_disabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether to disable the function of automatically synchronizing labels and taints to existing nodes. Default is false.",
						},
						"kubelet_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The KubeletConfig of KubernetesConfig. After adding parameters, deleting parameters does not take effect.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"topology_manager_scope": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The TopologyManagerScope of KubeletConfig. Valid values: `container`.",
									},
									"topology_manager_policy": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The TopologyManagerPolicy of KubeletConfig. Valid values: `none`, `restricted`, `best-effort`, `single-numa-node`. Default is `none`.",
									},
									"feature_gates": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "The FeatureGates of KubeletConfig.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"qos_resource_manager": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to enable QoSResourceManager. Default is false.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Description: "The KubernetesConfig of NodeConfig.",
			},

			// computed fields
			"node_statistics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The TotalCount of Node.",
						},
						"creating_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The CreatingCount of Node.",
						},
						"running_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The RunningCount of Node.",
						},
						"updating_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The UpdatingCount of Node.",
						},
						"deleting_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The DeletingCount of Node.",
						},
						"failed_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The FailedCount of Node.",
						},
						"stopped_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Deprecated:  "This field has been deprecated and is not recommended for use.",
							Description: "The StoppedCount of Node.",
						},
						"stopping_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Deprecated:  "This field has been deprecated and is not recommended for use.",
							Description: "The StoppingCount of Node.",
						},
						"starting_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Deprecated:  "This field has been deprecated and is not recommended for use.",
							Description: "The StartingCount of Node.",
						},
					},
				},
				Description: "The NodeStatistics of NodeConfig.",
			},
		},
	}
}

func resourceVestackNodePoolCreate(d *schema.ResourceData, meta interface{}) (err error) {
	nodePoolService := NewNodePoolService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(nodePoolService, d, ResourceVestackNodePool())
	if err != nil {
		return fmt.Errorf("error on creating nodePoolService  %q, %w", d.Id(), err)
	}
	return resourceVestackNodePoolRead(d, meta)
}

func resourceVestackNodePoolRead(d *schema.ResourceData, meta interface{}) (err error) {
	nodePoolService := NewNodePoolService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(nodePoolService, d, ResourceVestackNodePool())
	if err != nil {
		return fmt.Errorf("error on reading nodePoolService %q, %w", d.Id(), err)
	}
	return err
}

func resourceVestackNodePoolUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	nodePoolService := NewNodePoolService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(nodePoolService, d, ResourceVestackNodePool())
	if err != nil {
		return fmt.Errorf("error on updating nodePoolService  %q, %w", d.Id(), err)
	}
	return resourceVestackNodePoolRead(d, meta)
}

func resourceVestackNodePoolDelete(d *schema.ResourceData, meta interface{}) (err error) {
	nodePoolService := NewNodePoolService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(nodePoolService, d, ResourceVestackNodePool())
	if err != nil {
		return fmt.Errorf("error on deleting nodePoolService %q, %w", d.Id(), err)
	}
	return err
}
