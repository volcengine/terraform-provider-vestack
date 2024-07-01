package route_table

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
Route table can be imported using the id, e.g.
```
$ terraform import vestack_route_table.default vtb-274e0syt9av407fap8tle16kb
```

*/

func ResourceVestackRouteTable() *schema.Resource {
	return &schema.Resource{
		Delete: resourceVestackRouteTableDelete,
		Create: resourceVestackRouteTableCreate,
		Read:   resourceVestackRouteTableRead,
		Update: resourceVestackRouteTableUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the VPC.",
			},
			"route_table_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the route table.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the route table.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ProjectName of the route table.",
			},
		},
	}
}

func resourceVestackRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	routeTableService := NewRouteTableService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(routeTableService, d, ResourceVestackRouteTable()); err != nil {
		return fmt.Errorf("error on creating route table  %q, %w", d.Id(), err)
	}
	return resourceVestackRouteTableRead(d, meta)
}

func resourceVestackRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	routeTableService := NewRouteTableService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(routeTableService, d, ResourceVestackRouteTable()); err != nil {
		return fmt.Errorf("error on reading route table %q, %w", d.Id(), err)
	}
	return nil
}

func resourceVestackRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {
	routeTableService := NewRouteTableService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Update(routeTableService, d, ResourceVestackRouteTable()); err != nil {
		return fmt.Errorf("error on updating route table %q, %w", d.Id(), err)
	}
	return resourceVestackRouteTableRead(d, meta)
}

func resourceVestackRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	routeTableService := NewRouteTableService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(routeTableService, d, ResourceVestackRouteTable()); err != nil {
		return fmt.Errorf("error on deleting route table %q, %w", d.Id(), err)
	}
	return nil
}
