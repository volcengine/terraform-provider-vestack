package ecs_hpc_cluster

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
EcsHpcCluster can be imported using the id, e.g.
```
$ terraform import vestack_ecs_hpc_cluster.default resource_id
```

*/

func ResourceVestackEcsHpcCluster() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackEcsHpcClusterCreate,
		Read:   resourceVestackEcsHpcClusterRead,
		Update: resourceVestackEcsHpcClusterUpdate,
		Delete: resourceVestackEcsHpcClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The zone id of the hpc cluster.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the hpc cluster.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the hpc cluster.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The project name of the hpc cluster.",
			},
			"tags": bp.TagsSchema(),
		},
	}
	return resource
}

func resourceVestackEcsHpcClusterCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEcsHpcClusterService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackEcsHpcCluster())
	if err != nil {
		return fmt.Errorf("error on creating ecs_hpc_cluster %q, %s", d.Id(), err)
	}
	return resourceVestackEcsHpcClusterRead(d, meta)
}

func resourceVestackEcsHpcClusterRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEcsHpcClusterService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackEcsHpcCluster())
	if err != nil {
		return fmt.Errorf("error on reading ecs_hpc_cluster %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackEcsHpcClusterUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEcsHpcClusterService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackEcsHpcCluster())
	if err != nil {
		return fmt.Errorf("error on updating ecs_hpc_cluster %q, %s", d.Id(), err)
	}
	return resourceVestackEcsHpcClusterRead(d, meta)
}

func resourceVestackEcsHpcClusterDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewEcsHpcClusterService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackEcsHpcCluster())
	if err != nil {
		return fmt.Errorf("error on deleting ecs_hpc_cluster %q, %s", d.Id(), err)
	}
	return err
}
