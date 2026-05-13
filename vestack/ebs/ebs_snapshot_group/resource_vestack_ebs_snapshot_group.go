package ebs_snapshot_group

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
EbsSnapshotGroup can be imported using the id, e.g.
```
$ terraform import vestack_ebs_snapshot_group.default resource_id
```

*/

func ResourceVestackEbsSnapshotGroup() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackEbsSnapshotGroupCreate,
		Read:   resourceVestackEbsSnapshotGroupRead,
		Update: resourceVestackEbsSnapshotGroupUpdate,
		Delete: resourceVestackEbsSnapshotGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"volume_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The volume id of the snapshot group. The status of the volume must be `attached`." +
					"If multiple volumes are specified, they need to be attached to the same ECS instance.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The instance id of the snapshot group.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the snapshot group.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The instance id of the snapshot group.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project name of the snapshot group.",
			},
			"tags": bp.TagsSchema(),

			// computed fields
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the snapshot group.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The image id of the snapshot group.",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the snapshot group.",
			},
		},
	}
	return resource
}

func resourceVestackEbsSnapshotGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEbsSnapshotGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackEbsSnapshotGroup())
	if err != nil {
		return fmt.Errorf("error on creating ebs_snapshot_group %q, %s", d.Id(), err)
	}
	return resourceVestackEbsSnapshotGroupRead(d, meta)
}

func resourceVestackEbsSnapshotGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEbsSnapshotGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackEbsSnapshotGroup())
	if err != nil {
		return fmt.Errorf("error on reading ebs_snapshot_group %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackEbsSnapshotGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEbsSnapshotGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackEbsSnapshotGroup())
	if err != nil {
		return fmt.Errorf("error on updating ebs_snapshot_group %q, %s", d.Id(), err)
	}
	return resourceVestackEbsSnapshotGroupRead(d, meta)
}

func resourceVestackEbsSnapshotGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEbsSnapshotGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackEbsSnapshotGroup())
	if err != nil {
		return fmt.Errorf("error on deleting ebs_snapshot_group %q, %s", d.Id(), err)
	}
	return err
}
