package traffic_mirror_filter

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
TrafficMirrorFilter can be imported using the id, e.g.
```
$ terraform import volcengine_traffic_mirror_filter.default resource_id
```

*/

func ResourceVestackTrafficMirrorFilter() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackTrafficMirrorFilterCreate,
		Read:   resourceVestackTrafficMirrorFilterRead,
		Update: resourceVestackTrafficMirrorFilterUpdate,
		Delete: resourceVestackTrafficMirrorFilterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"traffic_mirror_filter_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the traffic mirror filter.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the traffic mirror filter.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The project name of the traffic mirror filter.",
			},
			"tags": ve.TagsSchema(),

			// computed fields
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of traffic mirror filter.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create time of traffic mirror filter.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update time of traffic mirror filter.",
			},
		},
	}
	return resource
}

func resourceVestackTrafficMirrorFilterCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTrafficMirrorFilterService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackTrafficMirrorFilter())
	if err != nil {
		return fmt.Errorf("error on creating traffic_mirror_filter %q, %s", d.Id(), err)
	}
	return resourceVestackTrafficMirrorFilterRead(d, meta)
}

func resourceVestackTrafficMirrorFilterRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTrafficMirrorFilterService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackTrafficMirrorFilter())
	if err != nil {
		return fmt.Errorf("error on reading traffic_mirror_filter %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackTrafficMirrorFilterUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTrafficMirrorFilterService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackTrafficMirrorFilter())
	if err != nil {
		return fmt.Errorf("error on updating traffic_mirror_filter %q, %s", d.Id(), err)
	}
	return resourceVestackTrafficMirrorFilterRead(d, meta)
}

func resourceVestackTrafficMirrorFilterDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTrafficMirrorFilterService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackTrafficMirrorFilter())
	if err != nil {
		return fmt.Errorf("error on deleting traffic_mirror_filter %q, %s", d.Id(), err)
	}
	return err
}
