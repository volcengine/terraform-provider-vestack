package ecs_key_pair

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
ECS key pair can be imported using the id, e.g.
```
$ terraform import vestack_ecs_key_pair.default kp-mizl7m1kqccg5smt1bdpijuj
```

*/

func ResourceVestackEcsKeyPair() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackEcsKeyPairCreate,
		Read:   resourceVestackEcsKeyPairRead,
		Update: resourceVestackEcsKeyPairUpdate,
		Delete: resourceVestackEcsKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"key_pair_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 64),
				Description:  "The name of key pair.",
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch ele := v.(type) {
					case string:
						return strings.TrimSpace(ele)
					default:
						return ""
					}
				},
				Description: "Public key string.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of key pair.",
			},
			"key_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Target file to save private key. It is recommended that the value not be empty. " +
					"You only have one chance to download the private key, the vestack will not save your private key, please keep it safe. " +
					"In the TF import scenario, this field will not write the private key locally.",
			},
			"finger_print": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The finger print info.",
			},
			"key_pair_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of key pair.",
			},
		},
	}
	return resource
}

func resourceVestackEcsKeyPairCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEcsKeyPairService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceVestackEcsKeyPair())
	if err != nil {
		return fmt.Errorf("error on creating ecs key pair  %q, %s", d.Id(), err)
	}
	return resourceVestackEcsKeyPairRead(d, meta)
}

func resourceVestackEcsKeyPairRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEcsKeyPairService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceVestackEcsKeyPair())
	if err != nil {
		return fmt.Errorf("error on reading ecs key pair %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackEcsKeyPairUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEcsKeyPairService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceVestackEcsKeyPair())
	if err != nil {
		return fmt.Errorf("error on updating ecs key pair  %q, %s", d.Id(), err)
	}
	return resourceVestackEcsKeyPairRead(d, meta)
}

func resourceVestackEcsKeyPairDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEcsKeyPairService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceVestackEcsKeyPair())
	if err != nil {
		return fmt.Errorf("error on deleting ecs key pair %q, %s", d.Id(), err)
	}
	return err
}
