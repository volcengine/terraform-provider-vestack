package direct_connect_connection

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
DirectConnectConnection can be imported using the id, e.g.
```
$ terraform import vestack_direct_connect_connection.default dcc-7qthudw0ll6jmc****
```

*/

func ResourceVestackDirectConnectConnection() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackDirectConnectConnectionCreate,
		Read:   resourceVestackDirectConnectConnectionRead,
		Update: resourceVestackDirectConnectConnectionUpdate,
		Delete: resourceVestackDirectConnectConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"direct_connect_connection_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of direct connect.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of direct connect.",
			},
			"port_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The direct connect access point port id.",
			},
			"owner_account_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The direct connect connection owner account id.",
			},
			"owner_project_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The direct connect connection owner project name.",
			},
			"line_operator": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The physical leased line operator.valid value contains `ChinaTelecom`,`ChinaMobile`,`ChinaUnicom`,`ChinaOther`.",
			},
			"port_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The physical leased line port type and spec.valid value contains `1000Base-T`,`10GBase-T`,`1000Base`,`10GBase`,`40GBase`,`100GBase`.",
			},
			"port_spec": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The physical leased line port spec.valid value contains `1G`,`10G`.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The line band width,unit:Mbps.",
			},
			"peer_location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The local IDC address.",
			},
			"customer_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dedicated line contact name.",
			},
			"customer_contact_phone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dedicated line contact phone.",
			},
			"customer_contact_email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dedicated line contact email.",
			},
		},
	}
	return resource
}

func resourceVestackDirectConnectConnectionCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectConnectionService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackDirectConnectConnection())
	if err != nil {
		return fmt.Errorf("error on creating direct_connect_connection %q, %s", d.Id(), err)
	}
	return resourceVestackDirectConnectConnectionRead(d, meta)
}

func resourceVestackDirectConnectConnectionRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectConnectionService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackDirectConnectConnection())
	if err != nil {
		return fmt.Errorf("error on reading direct_connect_connection %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackDirectConnectConnectionUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectConnectionService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackDirectConnectConnection())
	if err != nil {
		return fmt.Errorf("error on updating direct_connect_connection %q, %s", d.Id(), err)
	}
	return resourceVestackDirectConnectConnectionRead(d, meta)
}

func resourceVestackDirectConnectConnectionDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectConnectionService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackDirectConnectConnection())
	if err != nil {
		return fmt.Errorf("error on deleting direct_connect_connection %q, %s", d.Id(), err)
	}
	return err
}
