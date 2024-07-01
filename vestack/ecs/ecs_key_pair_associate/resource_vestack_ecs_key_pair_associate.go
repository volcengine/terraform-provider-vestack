package ecs_key_pair_associate

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
ECS key pair associate can be imported using the id, e.g.

After binding the key pair, the instance needs to be restarted for the key pair to take effect.

After the key pair is bound, the password login method will automatically become invalid. If your instance has been set for password login, after the key pair is bound, you will no longer be able to use the password login method.

```
$ terraform import vestack_ecs_key_pair_associate.default kp-ybti5tkpkv2udbfolrft:i-mizl7m1kqccg5smt1bdpijuj
```

*/

func ResourceVestackEcsKeyPairAssociate() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackEcsKeyPairAssociateCreate,
		Read:   resourceVestackEcsKeyPairAssociateRead,
		Delete: resourceVestackEcsKeyPairAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				items := strings.Split(data.Id(), ":")
				if len(items) != 2 {
					return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
				}
				if err := data.Set("key_pair_id", items[0]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				if err := data.Set("instance_id", items[1]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				return []*schema.ResourceData{data}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"key_pair_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of ECS KeyPair Associate.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of ECS Instance.",
			},
		},
	}
	return resource
}

func resourceVestackEcsKeyPairAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEcsKeyPairAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceVestackEcsKeyPairAssociate())
	if err != nil {
		return fmt.Errorf("error on creating ecs key pair Associate %q, %s", d.Id(), err)
	}
	return resourceVestackEcsKeyPairAssociateRead(d, meta)
}

func resourceVestackEcsKeyPairAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEcsKeyPairAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceVestackEcsKeyPairAssociate())
	if err != nil {
		return fmt.Errorf("error on reading ecs key pair Associate %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackEcsKeyPairAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEcsKeyPairAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceVestackEcsKeyPairAssociate())
	if err != nil {
		return fmt.Errorf("error on deleting ecs key pair Associate %q, %s", d.Id(), err)
	}
	return err
}
