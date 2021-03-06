package listener

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ve "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Listener can be imported using the id, e.g.
```
$ terraform import vestack_listener.default lsn-273yv0mhs5xj47fap8sehiiso
```

*/

func ResourceVestackListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackListenerCreate,
		Read:   resourceVestackListenerRead,
		Update: resourceVestackListenerUpdate,
		Delete: resourceVestackListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the Listener.",
			},
			"load_balancer_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region of the request.",
			},
			"listener_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the Listener.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The protocol of the Listener. Optional choice contains `TCP`, `UDP`, `HTTP`, `HTTPS`.",
				ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "HTTP", "HTTPS"}, false),
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The port receiving request of the Listener.",
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The scheduling algorithm of the Listener. Optional choice contains `wrr`, `wlc`, `sh`.",
				ValidateFunc: validation.StringInSlice([]string{"wrr", "wlc", "sh"}, false),
			},
			"enabled": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The enable status of the Listener. Optional choice contains `on`, `off`.",
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
			"established_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The connection timeout of the Listener.",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The certificate id associated with the listener.",
			},
			"server_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The server group id associated with the listener.",
			},
			"acl_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The enable status of Acl. Optional choice contains `on`, `off`.",
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
			"acl_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The type of the Acl. Optional choice contains `white`, `black`.",
				ValidateFunc: validation.StringInSlice([]string{"white", "black"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("acl_status").(string) == "off"
				},
			},
			"acl_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The id list of the Acl.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("acl_status").(string) == "off"
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Listener.",
			},
			"health_check": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "The config of health check.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "The enable status of health check function. Optional choice contains `on`, `off`.",
							ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
						},
						"interval": {
							Type:             schema.TypeInt,
							Optional:         true,
							Description:      "The interval executing health check.",
							DiffSuppressFunc: HealthCheckFieldDiffSuppress,
						},
						"timeout": {
							Type:             schema.TypeInt,
							Optional:         true,
							Description:      "The response timeout of health check.",
							DiffSuppressFunc: HealthCheckFieldDiffSuppress,
						},
						"healthy_threshold": {
							Type:             schema.TypeInt,
							Optional:         true,
							Description:      "The healthy threshold of health check.",
							DiffSuppressFunc: HealthCheckFieldDiffSuppress,
						},
						"un_healthy_threshold": {
							Type:             schema.TypeInt,
							Optional:         true,
							Description:      "The unhealthy threshold of health check.",
							DiffSuppressFunc: HealthCheckFieldDiffSuppress,
						},
						"method": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "The method of health check.",
							DiffSuppressFunc: HealthCheckHTTPOnlyFieldDiffSuppress,
						},
						"domain": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "The domain of health check.",
							DiffSuppressFunc: HealthCheckHTTPOnlyFieldDiffSuppress,
						},
						"uri": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "The uri of health check.",
							DiffSuppressFunc: HealthCheckHTTPOnlyFieldDiffSuppress,
						},
						"http_code": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "The normal http status code of health check.",
							DiffSuppressFunc: HealthCheckHTTPOnlyFieldDiffSuppress,
						},
					},
				},
			},
		},
	}
}

func resourceVestackListenerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	listenerService := NewListenerService(meta.(*ve.SdkClient))
	err = listenerService.Dispatcher.Create(listenerService, d, ResourceVestackListener())
	if err != nil {
		return fmt.Errorf("error on creating listener  %q, %w", d.Id(), err)
	}
	return resourceVestackListenerRead(d, meta)
}

func resourceVestackListenerRead(d *schema.ResourceData, meta interface{}) (err error) {
	listenerService := NewListenerService(meta.(*ve.SdkClient))
	err = listenerService.Dispatcher.Read(listenerService, d, ResourceVestackListener())
	if err != nil {
		return fmt.Errorf("error on reading listener %q, %w", d.Id(), err)
	}
	return err
}

func resourceVestackListenerUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	listenerService := NewListenerService(meta.(*ve.SdkClient))
	err = listenerService.Dispatcher.Update(listenerService, d, ResourceVestackListener())
	if err != nil {
		return fmt.Errorf("error on updating listener  %q, %w", d.Id(), err)
	}
	return resourceVestackListenerRead(d, meta)
}

func resourceVestackListenerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	listenerService := NewListenerService(meta.(*ve.SdkClient))
	err = listenerService.Dispatcher.Delete(listenerService, d, ResourceVestackListener())
	if err != nil {
		return fmt.Errorf("error on deleting listener %q, %w", d.Id(), err)
	}
	return err
}
