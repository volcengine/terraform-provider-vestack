package dnat_entry

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Dnat entry can be imported using the id, e.g.
```
$ terraform import vestack_dnat_entry.default dnat-3fvhk47kf56****
```

*/

func ResourceVestackDnatEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackDnatEntryCreate,
		Update: resourceVestackDnatEntryUpdate,
		Read:   resourceVestackDnatEntryRead,
		Delete: resourceVestackDnatEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"nat_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the nat gateway to which the entry belongs.",
			},
			"external_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Provides the public IP address for public network access.",
			},
			"external_port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The port or port segment that receives requests from the public network. If InternalPort is passed into the port segment, ExternalPort must also be passed into the port segment.",
			},
			"internal_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Provides the internal IP address.",
			},
			"internal_port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The port or port segment on which the cloud server instance provides services to the public network.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
				Description:  "The network protocol.",
			},
			"dnat_entry_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the DNAT rule.",
			},
			"dnat_entry_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the DNAT rule.",
			},
		},
	}
}

func resourceVestackDnatEntryCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDnatEntryService(meta.(*ve.SdkClient))
	err = ve.DefaultDispatcher().Create(service, d, ResourceVestackDnatEntry())
	if err != nil {
		return fmt.Errorf("error on creating dnat entry: %q, %w", d.Id(), err)
	}
	return resourceVestackDnatEntryRead(d, meta)
}

func resourceVestackDnatEntryRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDnatEntryService(meta.(*ve.SdkClient))
	err = ve.DefaultDispatcher().Read(service, d, ResourceVestackDnatEntry())
	if err != nil {
		return fmt.Errorf("error on reading dnat entry: %q, %w", d.Id(), err)
	}
	return nil
}

func resourceVestackDnatEntryUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDnatEntryService(meta.(*ve.SdkClient))
	err = ve.DefaultDispatcher().Update(service, d, ResourceVestackDnatEntry())
	if err != nil {
		return fmt.Errorf("error on updating dnat entry: %q, %w", d.Id(), err)
	}
	return resourceVestackDnatEntryRead(d, meta)
}

func resourceVestackDnatEntryDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDnatEntryService(meta.(*ve.SdkClient))
	err = ve.DefaultDispatcher().Delete(service, d, ResourceVestackDnatEntry())
	if err != nil {
		return fmt.Errorf("error on deleting dnat entry: %q, %w", d.Id(), err)
	}
	return nil
}
