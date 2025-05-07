package direct_connect_bgp_peer

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
DirectConnectBgpPeer can be imported using the id, e.g.
```
$ terraform import volcengine_direct_connect_bgp_peer.default bgp-2752hz4teko3k7fap8u4c****
```

*/

func ResourceVestackDirectConnectBgpPeer() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackDirectConnectBgpPeerCreate,
		Read:   resourceVestackDirectConnectBgpPeerRead,
		Update: resourceVestackDirectConnectBgpPeerUpdate,
		Delete: resourceVestackDirectConnectBgpPeerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bgp_peer_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of bgp peer.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of bgp peer.",
			},
			"auth_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The auth key of bgp peer.",
			},
			"remote_asn": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The remote asn of bgp peer.",
			},
			"virtual_interface_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of virtual interface.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The account id of direct connect gateway routes.",
			},
		},
	}
	dataSource := DataSourceVestackDirectConnectBgpPeers().Schema["bgp_peers"].Elem.(*schema.Resource).Schema
	ve.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceVestackDirectConnectBgpPeerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectBgpPeerService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackDirectConnectBgpPeer())
	if err != nil {
		return fmt.Errorf("error on creating direct_connect_bgp_peer %q, %s", d.Id(), err)
	}
	return resourceVestackDirectConnectBgpPeerRead(d, meta)
}

func resourceVestackDirectConnectBgpPeerRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectBgpPeerService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackDirectConnectBgpPeer())
	if err != nil {
		return fmt.Errorf("error on reading direct_connect_bgp_peer %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackDirectConnectBgpPeerUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectBgpPeerService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackDirectConnectBgpPeer())
	if err != nil {
		return fmt.Errorf("error on updating direct_connect_bgp_peer %q, %s", d.Id(), err)
	}
	return resourceVestackDirectConnectBgpPeerRead(d, meta)
}

func resourceVestackDirectConnectBgpPeerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectBgpPeerService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackDirectConnectBgpPeer())
	if err != nil {
		return fmt.Errorf("error on deleting direct_connect_bgp_peer %q, %s", d.Id(), err)
	}
	return err
}
