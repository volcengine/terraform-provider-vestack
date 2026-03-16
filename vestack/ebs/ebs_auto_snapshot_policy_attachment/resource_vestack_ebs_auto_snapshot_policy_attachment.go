package ebs_auto_snapshot_policy_attachment

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
EbsAutoSnapshotPolicyAttachment can be imported using the auto_snapshot_policy_id:volume_id, e.g.
```
$ terraform import vestack_ebs_auto_snapshot_policy_attachment.default resource_id
```

*/

func ResourceVestackEbsAutoSnapshotPolicyAttachment() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackEbsAutoSnapshotPolicyAttachmentCreate,
		Read:   resourceVestackEbsAutoSnapshotPolicyAttachmentRead,
		Delete: resourceVestackEbsAutoSnapshotPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: importAutoSnapshotPolicyAttachment,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_snapshot_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the auto snapshot policy.",
			},
			"volume_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the volume.",
			},
		},
	}
	return resource
}

func resourceVestackEbsAutoSnapshotPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEbsAutoSnapshotPolicyAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackEbsAutoSnapshotPolicyAttachment())
	if err != nil {
		return fmt.Errorf("error on creating ebs_auto_snapshot_policy_attachment %q, %s", d.Id(), err)
	}
	return resourceVestackEbsAutoSnapshotPolicyAttachmentRead(d, meta)
}

func resourceVestackEbsAutoSnapshotPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEbsAutoSnapshotPolicyAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackEbsAutoSnapshotPolicyAttachment())
	if err != nil {
		return fmt.Errorf("error on reading ebs_auto_snapshot_policy_attachment %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackEbsAutoSnapshotPolicyAttachmentUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEbsAutoSnapshotPolicyAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackEbsAutoSnapshotPolicyAttachment())
	if err != nil {
		return fmt.Errorf("error on updating ebs_auto_snapshot_policy_attachment %q, %s", d.Id(), err)
	}
	return resourceVestackEbsAutoSnapshotPolicyAttachmentRead(d, meta)
}

func resourceVestackEbsAutoSnapshotPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEbsAutoSnapshotPolicyAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackEbsAutoSnapshotPolicyAttachment())
	if err != nil {
		return fmt.Errorf("error on deleting ebs_auto_snapshot_policy_attachment %q, %s", d.Id(), err)
	}
	return err
}

func importAutoSnapshotPolicyAttachment(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(d.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{d}, fmt.Errorf("import id must be of the form auto_snapshot_policy_id:volume_id")
	}
	err = d.Set("auto_snapshot_policy_id", items[0])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	err = d.Set("volume_id", items[1])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}

	return []*schema.ResourceData{d}, nil
}
