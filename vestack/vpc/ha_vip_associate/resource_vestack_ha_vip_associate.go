package ha_vip_associate

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
HaVipAssociate can be imported using the ha_vip_id:instance_id, e.g.
```
$ terraform import vestack_ha_vip_associate.default havip-2byzv8icq1b7k2dx0eegb****:eni-2d5wv84h7onpc58ozfeeu****
```

*/

func ResourceVestackHaVipAssociate() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackHaVipAssociateCreate,
		Read:   resourceVestackHaVipAssociateRead,
		Delete: resourceVestackHaVipAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: haVipAssociateImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ha_vip_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the Ha Vip.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the associated instance.",
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "EcsInstance",
				ValidateFunc: validation.StringInSlice([]string{"EcsInstance", "NetworkInterface"}, false),
				Description:  "The type of the associated instance. Valid values: `EcsInstance`, `NetworkInterface`.",
			},
		},
	}
	return resource
}

func resourceVestackHaVipAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewHaVipAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackHaVipAssociate())
	if err != nil {
		return fmt.Errorf("error on creating ha_vip_associate %q, %s", d.Id(), err)
	}
	return resourceVestackHaVipAssociateRead(d, meta)
}

func resourceVestackHaVipAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewHaVipAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackHaVipAssociate())
	if err != nil {
		return fmt.Errorf("error on reading ha_vip_associate %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackHaVipAssociateUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewHaVipAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackHaVipAssociate())
	if err != nil {
		return fmt.Errorf("error on updating ha_vip_associate %q, %s", d.Id(), err)
	}
	return resourceVestackHaVipAssociateRead(d, meta)
}

func resourceVestackHaVipAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewHaVipAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackHaVipAssociate())
	if err != nil {
		return fmt.Errorf("error on deleting ha_vip_associate %q, %s", d.Id(), err)
	}
	return err
}

var haVipAssociateImporter = func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
	}
	if err := data.Set("ha_vip_id", items[0]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	if err := data.Set("instance_id", items[1]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
