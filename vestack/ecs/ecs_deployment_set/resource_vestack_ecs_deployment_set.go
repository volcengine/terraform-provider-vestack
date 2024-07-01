package ecs_deployment_set

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
ECS deployment set can be imported using the id, e.g.
```
$ terraform import vestack_ecs_deployment_set.default i-mizl7m1kqccg5smt1bdpijuj
```

*/

func ResourceVestackEcsDeploymentSet() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackEcsDeploymentSetCreate,
		Read:   resourceVestackEcsDeploymentSetRead,
		Update: resourceVestackEcsDeploymentSetUpdate,
		Delete: resourceVestackEcsDeploymentSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"deployment_set_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of ECS DeploymentSet.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of ECS DeploymentSet.",
			},
			"granularity": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"switch",
					"host",
					"rack",
				}, false),
				Default:     "host",
				Description: "The granularity of ECS DeploymentSet.Valid values: switch, host, rack,Default is host.",
			},
			"strategy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Availability",
				}, false),
				Default:     "Availability",
				Description: "The strategy of ECS DeploymentSet.Valid values: Availability.Default is Availability.",
			},
			"deployment_set_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of ECS DeploymentSet.",
			},
		},
	}
	return resource
}

func resourceVestackEcsDeploymentSetCreate(d *schema.ResourceData, meta interface{}) (err error) {
	deploymentSetService := NewEcsDeploymentSetService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(deploymentSetService, d, ResourceVestackEcsDeploymentSet())
	if err != nil {
		return fmt.Errorf("error on creating ecs deployment set  %q, %s", d.Id(), err)
	}
	return resourceVestackEcsDeploymentSetRead(d, meta)
}

func resourceVestackEcsDeploymentSetRead(d *schema.ResourceData, meta interface{}) (err error) {
	deploymentSetService := NewEcsDeploymentSetService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(deploymentSetService, d, ResourceVestackEcsDeploymentSet())
	if err != nil {
		return fmt.Errorf("error on reading ecs deployment set %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackEcsDeploymentSetUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	deploymentSetService := NewEcsDeploymentSetService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(deploymentSetService, d, ResourceVestackEcsDeploymentSet())
	if err != nil {
		return fmt.Errorf("error on updating ecs deployment set  %q, %s", d.Id(), err)
	}
	return resourceVestackEcsDeploymentSetRead(d, meta)
}

func resourceVestackEcsDeploymentSetDelete(d *schema.ResourceData, meta interface{}) (err error) {
	deploymentSetService := NewEcsDeploymentSetService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(deploymentSetService, d, ResourceVestackEcsDeploymentSet())
	if err != nil {
		return fmt.Errorf("error on deleting ecs deployment set %q, %s", d.Id(), err)
	}
	return err
}
