package addon

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
VkeAddon can be imported using the clusterId:Name, e.g.
```
$ terraform import vestack_vke_addon.default cc9l74mvqtofjnoj5****:nginx-ingress
```

Notice
Some kind of VKEAddon can not be removed from vestack, and it will make a forbidden error when try to destroy.
If you want to remove it from terraform state, please use command
```
$ terraform state rm vestack_vke_addon.${name}
```

*/

func ResourceVestackVkeAddon() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackVkeAddonCreate,
		Read:   resourceVestackVkeAddonRead,
		Update: resourceVestackVkeAddonUpdate,
		Delete: resourceVestackVkeAddonDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				items := strings.Split(data.Id(), ":")
				if len(items) != 2 {
					return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
				}
				if err := data.Set("cluster_id", items[0]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				if err := data.Set("name", items[1]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				return []*schema.ResourceData{data}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The cluster id of the addon.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the addon.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The version info of the cluster.",
			},
			"deploy_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The deploy mode.",
			},
			"deploy_node_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The deploy node type.",
			},
			"config": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "The config info of addon. " +
					"Please notice that `ingress-nginx` component prohibits updating config, can only works on the web console.",
			},
		},
	}
}

func resourceVestackVkeAddonCreate(d *schema.ResourceData, meta interface{}) (err error) {
	clusterService := NewVkeAddonService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(clusterService, d, ResourceVestackVkeAddon())
	if err != nil {
		return fmt.Errorf("error on creating addon  %q, %w", d.Id(), err)
	}
	return resourceVestackVkeAddonRead(d, meta)
}

func resourceVestackVkeAddonRead(d *schema.ResourceData, meta interface{}) (err error) {
	clusterService := NewVkeAddonService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(clusterService, d, ResourceVestackVkeAddon())
	if err != nil {
		return fmt.Errorf("error on reading addon %q, %w", d.Id(), err)
	}
	return err
}

func resourceVestackVkeAddonUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	clusterService := NewVkeAddonService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(clusterService, d, ResourceVestackVkeAddon())
	if err != nil {
		return fmt.Errorf("error on updating addon  %q, %w", d.Id(), err)
	}
	return resourceVestackVkeAddonRead(d, meta)
}

func resourceVestackVkeAddonDelete(d *schema.ResourceData, meta interface{}) (err error) {
	clusterService := NewVkeAddonService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(clusterService, d, ResourceVestackVkeAddon())
	if err != nil {
		return fmt.Errorf("error on deleting addon %q, %w", d.Id(), err)
	}
	return err
}
