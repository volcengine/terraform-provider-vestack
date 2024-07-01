package network_acl_associate

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
NetworkAcl associate can be imported using the network_acl_id:resource_id, e.g.
```
$ terraform import vestack_network_acl_associate.default nacl-172leak37mi9s4d1w33pswqkh:subnet-637jxq81u5mon3gd6ivc7rj
```

*/

func ResourceVestackNetworkAclAssociate() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackAclAssociateCreate,
		Read:   resourceVestackAclAssociateRead,
		Delete: resourceVestackAclAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: aclAssociateImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"network_acl_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of Network Acl.",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The resource id of Network Acl.",
			},
		},
	}
}

func resourceVestackAclAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	aclAssociateService := NewNetworkAclAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(aclAssociateService, d, ResourceVestackNetworkAclAssociate())
	if err != nil {
		return fmt.Errorf("error on creating acl Associate %q, %w", d.Id(), err)
	}
	return resourceVestackAclAssociateRead(d, meta)
}

func resourceVestackAclAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	aclAssociateService := NewNetworkAclAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(aclAssociateService, d, ResourceVestackNetworkAclAssociate())
	if err != nil {
		return fmt.Errorf("error on reading acl Associate %q, %w", d.Id(), err)
	}
	return err
}

func resourceVestackAclAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	aclAssociateService := NewNetworkAclAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(aclAssociateService, d, ResourceVestackNetworkAclAssociate())
	if err != nil {
		return fmt.Errorf("error on deleting acl Associate %q, %w", d.Id(), err)
	}
	return err
}
