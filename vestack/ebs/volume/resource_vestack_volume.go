package volume

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Volume can be imported using the id, e.g.
```
$ terraform import vestack_volume.default vol-mizl7m1kqccg5smt1bdpijuj
```

*/

func ResourceVestackVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackVolumeCreate,
		Read:   resourceVestackVolumeRead,
		Update: resourceVestackVolumeUpdate,
		Delete: resourceVestackVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the Zone.",
			},
			"volume_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of Volume.",
			},
			"volume_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The type of Volume.",
				ValidateFunc: validation.StringInSlice([]string{"ESSD_PL0", "ESSD_PL1", "ESSD_PL2", "PTSSD"}, false),
			},
			"kind": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"data"}, false),
				Description:  "The kind of Volume.",
			},
			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(20), // 最小20GB
				Description:  "The size of Volume.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Volume.",
			},
			"volume_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PostPaid"}, false),
				Default:      "PostPaid",
				Description:  "The charge type of the Volume.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of Volume.",
			},
			"trade_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status of Trade.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of Volume.",
			},
			"billing_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Billing type of Volume.",
			},
			"pay_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Pay type of Volume.",
			},
			"delete_with_instance": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Delete Volume with Attached Instance.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
		},
	}
}

func resourceVestackVolumeCreate(d *schema.ResourceData, meta interface{}) (err error) {
	volumeService := NewVolumeService(meta.(*ve.SdkClient))
	err = volumeService.Dispatcher.Create(volumeService, d, ResourceVestackVolume())
	if err != nil {
		return fmt.Errorf("error on creating volume %q, %w", d.Id(), err)
	}
	return resourceVestackVolumeRead(d, meta)
}

func resourceVestackVolumeRead(d *schema.ResourceData, meta interface{}) (err error) {
	volumeService := NewVolumeService(meta.(*ve.SdkClient))
	err = volumeService.Dispatcher.Read(volumeService, d, ResourceVestackVolume())
	if err != nil {
		return fmt.Errorf("error on reading volume %q, %w", d.Id(), err)
	}
	return err
}

func resourceVestackVolumeUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	volumeService := NewVolumeService(meta.(*ve.SdkClient))
	err = volumeService.Dispatcher.Update(volumeService, d, ResourceVestackVolume())
	if err != nil {
		return fmt.Errorf("error on updating volume %q, %w", d.Id(), err)
	}
	return resourceVestackVolumeRead(d, meta)
}

func resourceVestackVolumeDelete(d *schema.ResourceData, meta interface{}) (err error) {
	volumeService := NewVolumeService(meta.(*ve.SdkClient))
	err = volumeService.Dispatcher.Delete(volumeService, d, ResourceVestackVolume())
	if err != nil {
		return fmt.Errorf("error on deleting volume %q, %w", d.Id(), err)
	}
	return err
}
