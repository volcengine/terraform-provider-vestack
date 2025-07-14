package ecs_instance

import (
	"fmt"
	"log"
	"math"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
ECS Instance can be imported using the id, e.g.
If Import,The data_volumes is sort by volume name
```
$ terraform import vestack_ecs_instance.default i-mizl7m1kqccg5smt1bdpijuj
```

*/

func ResourceVestackEcsInstance() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackEcsInstanceCreate,
		Read:   resourceVestackEcsInstanceRead,
		Update: resourceVestackEcsInstanceUpdate,
		Delete: resourceVestackEcsInstanceDelete,
		Exists: resourceVestackEcsInstanceExist,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The available zone ID of ECS instance.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Image ID of ECS instance.",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance type of ECS instance.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of ECS instance.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of ECS instance.",
			},
			"host_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The host name of ECS instance.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password of ECS instance.",
			},
			"key_pair_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The ssh key name of ECS instance.",
			},
			"keep_image_credential": {
				Type:     schema.TypeBool,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.HasChange("image_id")
				},
				Description: "Whether to keep the mirror settings. Only custom images and shared images support this field.\n When the value of this field is true, the Password and KeyPairName cannot be specified.\n When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"instance_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PostPaid",
					"PrePaid",
				}, false),
				Description: "The charge type of ECS instance, the value can be `PrePaid` or `PostPaid`.",
			},
			"spot_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NoSpot",
					"SpotAsPriceGo",
				}, false),
				Description: "The spot strategy will auto" +
					"remove instance in some conditions.Please make sure you can maintain instance lifecycle before " +
					"auto remove.The spot strategy of ECS instance, the value can be `NoSpot` or `SpotAsPriceGo`.",
			},
			"user_data": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: UserDateImportDiffSuppress,
				Description:      "The user data of ECS instance, this field must be encrypted with base64.",
			},
			"security_enhancement_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Active",
					"InActive",
				}, false),
				Default:     "Active",
				Description: "The security enhancement strategy of ECS instance. The value can be Active or InActive. Default is Active.When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"hpc_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The hpc cluster ID of ECS instance.",
			},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          12,
				DiffSuppressFunc: EcsInstanceImportDiffSuppress,
				Description:      "The period of ECS instance.Only effective when instance_charge_type is PrePaid. Default is 12. Unit is Month.",
			},
			//"period_unit": {
			//	Type:     schema.TypeString,
			//	Optional: true,
			//	Default:  "Month",
			//	ValidateFunc: validation.StringInSlice([]string{
			//		"Month",
			//	}, false),
			//	DiffSuppressFunc: bp.EcsInstanceImportDiffSuppress,
			//	Description:      "The period unit of ECS instance.Only effective when instance_charge_type is PrePaid. Default is Month.",
			//},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
				//ForceNew: true,
				Default:          true,
				DiffSuppressFunc: AutoRenewDiffSuppress,
				Description:      "The auto renew flag of ECS instance.Only effective when instance_charge_type is PrePaid. Default is true.When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"auto_renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
				//ForceNew: true,
				Default:          1,
				DiffSuppressFunc: AutoRenewDiffSuppress,
				Description:      "The auto renew period of ECS instance.Only effective when instance_charge_type is PrePaid. Default is 1.When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},

			"include_data_volumes": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: EcsInstanceImportDiffSuppress,
				Description:      "The include data volumes flag of ECS instance.Only effective when change instance charge type.include_data_volumes.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The subnet ID of primary networkInterface.",
			},

			"security_group_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				MaxItems:    5,
				MinItems:    1,
				Description: "The security group ID set of primary networkInterface.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"network_interface_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of primary networkInterface.",
			},

			"primary_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The private ip address of primary networkInterface.",
			},

			"system_volume_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The type of system volume, the value is `PTSSD` or `ESSD_PL0` or `ESSD_PL1` or `ESSD_PL2` or `ESSD_FlexPL`.",
			},

			"system_volume_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "The size of system volume. " +
					"The value range of the system volume size is ESSD_PL0: 20~2048, ESSD_FlexPL: 20~2048, PTSSD: 10~500.",
			},

			"system_volume_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of system volume.",
			},

			"deployment_set_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of Ecs Deployment Set.",
			},

			"ipv6_address_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Description:   "The number of IPv6 addresses to be automatically assigned from within the CIDR block of the subnet that hosts the ENI. Valid values: 1 to 10.",
				ValidateFunc:  validation.IntBetween(1, 10),
				ConflictsWith: []string{"ipv6_addresses"},
			},

			"ipv6_addresses": {
				Type:        schema.TypeSet,
				MaxItems:    10,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Set:         schema.HashString,
				Description: "One or more IPv6 addresses selected from within the CIDR block of the subnet that hosts the ENI. Support up to 10.\n You cannot specify both the ipv6_addresses and ipv6_address_count parameters.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"ipv6_address_count"},
			},

			"cpu_options": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				MinItems:    1,
				Description: "The option of cpu.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"threads_per_core": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "The per core of threads.",
						},
					},
				},
			},

			"data_volumes": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    15,
				MinItems:    1,
				Computed:    true,
				Description: "The data volumes collection of  ECS instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The type of volume, the value is `PTSSD` or `ESSD_PL0` or `ESSD_PL1` or `ESSD_PL2` or `ESSD_FlexPL`.",
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
							Description: "The size of volume. " +
								"The value range of the data volume size is ESSD_PL0: 10~32768, ESSD_FlexPL: 10~32768, PTSSD: 20~8192.",
						},
						"delete_with_instance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							ForceNew:    true,
							Description: "The delete with instance flag of volume.",
						},
					},
				},
			},

			"secondary_network_interfaces": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				MinItems:    1,
				Description: "The secondary networkInterface detail collection of ECS instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The subnet ID of secondary networkInterface.",
						},
						"security_group_ids": {
							Type:        schema.TypeSet,
							Required:    true,
							ForceNew:    true,
							MaxItems:    5,
							MinItems:    1,
							Description: "The security group ID set of secondary networkInterface.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},
						"primary_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The private ip address of secondary networkInterface.",
						},
					},
				},
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ProjectName of the ecs instance.",
			},
			"tags": bp.TagsSchema(),
			"ha_strategy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether the instance is turned on the high available mode, the value can be `offsite_rebuild` or empty string.",
			},

			"bms_system_disk_config": {
				Type:     schema.TypeList,
				Optional: true,
				//MaxItems:    1,
				//MinItems:    1,
				Description: "For bms only",
				//Set:         resourceBmsSystemDiskConfigHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity_gb": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The size of CapacityGB.",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The disk type.",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								log.Printf("[DEBUG] Comparing disk_type: old=%q, new=%q", old, new)
								return strings.EqualFold(old, new)
							},
						},
						"partitions": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Partitions configuration for BMS system disk",
							//Set:         resourceBmsPartitionHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"file_system": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "File system type of the partition.",
									},
									"mount_point": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Mount point of the partition.",
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											log.Printf("[DEBUG] Comparing mount_point: old=%q, new=%q", old, new)
											return path.Clean(old) == path.Clean(new)
										},
									},
									"size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Size of the partition.",
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											oldVal, _ := strconv.Atoi(old)
											newVal, _ := strconv.Atoi(new)
											log.Printf("[DEBUG] Comparing size: old=%d, new=%d, diff=%.2f%%",
												oldVal, newVal, 100*math.Abs(float64(oldVal-newVal))/float64(oldVal))
											return math.Abs(float64(oldVal-newVal))/float64(oldVal) < 0.05
										},
									},
								},
							},
						},
					},
				},
			},

			"bms_clean_data_disk": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether clean disk",
			},

			"bms_delete_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "1.Detach 2. WholeDisksErase 3. SystemDiskErase",
			},
		},
	}
	dataSource := DataSourceVestackEcsInstances().Schema["instances"].Elem.(*schema.Resource).Schema
	delete(dataSource, "network_interfaces")
	delete(dataSource, "volumes")
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceVestackEcsInstanceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	instanceService := NewEcsService(meta.(*bp.SdkClient))
	err = bp.NewRateLimitDispatcher(rateInfo).Create(instanceService, d, ResourceVestackEcsInstance())
	if err != nil {
		return fmt.Errorf("error on creating ecs instance  %q, %s", d.Id(), err)
	}
	return resourceVestackEcsInstanceRead(d, meta)
}

func resourceVestackEcsInstanceRead(d *schema.ResourceData, meta interface{}) (err error) {
	instanceService := NewEcsService(meta.(*bp.SdkClient))
	err = bp.NewRateLimitDispatcher(rateInfo).Read(instanceService, d, ResourceVestackEcsInstance())
	if err != nil {
		return fmt.Errorf("error on reading ecs instance %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackEcsInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	instanceService := NewEcsService(meta.(*bp.SdkClient))
	err = bp.NewRateLimitDispatcher(rateInfo).Update(instanceService, d, ResourceVestackEcsInstance())
	if err != nil {
		return fmt.Errorf("error on updating ecs instance  %q, %s", d.Id(), err)
	}
	return resourceVestackEcsInstanceRead(d, meta)
}

func resourceVestackEcsInstanceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	instanceService := NewEcsService(meta.(*bp.SdkClient))
	err = bp.NewRateLimitDispatcher(rateInfo).Delete(instanceService, d, ResourceVestackEcsInstance())
	if err != nil {
		return fmt.Errorf("error on deleting ecs instance %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackEcsInstanceExist(d *schema.ResourceData, meta interface{}) (flag bool, err error) {
	err = resourceVestackEcsInstanceRead(d, meta)
	if err != nil {
		if strings.Contains(err.Error(), "notfound") || strings.Contains(err.Error(), "not found") ||
			strings.Contains(err.Error(), "not exist") || strings.Contains(err.Error(), "not associate") ||
			strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "not_found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
