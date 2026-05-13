package ebs_snapshot

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
EbsSnapshot can be imported using the id, e.g.
```
$ terraform import vestack_ebs_snapshot.default resource_id
```

*/

func ResourceVestackEbsSnapshot() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackEbsSnapshotCreate,
		Read:   resourceVestackEbsSnapshotRead,
		Update: resourceVestackEbsSnapshotUpdate,
		Delete: resourceVestackEbsSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The volume id to create snapshot.",
			},
			"snapshot_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the snapshot.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the snapshot.",
			},
			"retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "The retention days of the snapshot. Valid values: 1~65536. Not specifying this field means permanently preserving the snapshot." +
					"When modifying this field, the retention days only supports extension and not shortening. The value range is N+1~65536, where N is the retention days set during snapshot creation.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project name of the snapshot.",
			},
			"tags": bp.TagsSchema(),

			// computed fields
			"snapshot_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the snapshot.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the snapshot.",
			},
			"volume_kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The volume kind of the snapshot.",
			},
			"volume_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The volume name of the snapshot.",
			},
			"volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The volume size of the snapshot.",
			},
			"volume_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The volume status of the snapshot.",
			},
			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The volume type of the snapshot.",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The zone id of the snapshot.",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the snapshot.",
			},
		},
	}
	return resource
}

func resourceVestackEbsSnapshotCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEbsSnapshotService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackEbsSnapshot())
	if err != nil {
		return fmt.Errorf("error on creating ebs_snapshot %q, %s", d.Id(), err)
	}
	return resourceVestackEbsSnapshotRead(d, meta)
}

func resourceVestackEbsSnapshotRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEbsSnapshotService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackEbsSnapshot())
	if err != nil {
		return fmt.Errorf("error on reading ebs_snapshot %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackEbsSnapshotUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEbsSnapshotService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackEbsSnapshot())
	if err != nil {
		return fmt.Errorf("error on updating ebs_snapshot %q, %s", d.Id(), err)
	}
	return resourceVestackEbsSnapshotRead(d, meta)
}

func resourceVestackEbsSnapshotDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEbsSnapshotService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackEbsSnapshot())
	if err != nil {
		return fmt.Errorf("error on deleting ebs_snapshot %q, %s", d.Id(), err)
	}
	return err
}
