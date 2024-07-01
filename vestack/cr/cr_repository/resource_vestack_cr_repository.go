package cr_repository

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
CR Repository can be imported using the registry:namespace:name, e.g.
```
$ terraform import vestack_cr_repository.default cr-basic:namespace-1:repo-1
```

*/

func crRepositoryImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(d.Id(), ":")
	if len(items) != 3 {
		return []*schema.ResourceData{d}, fmt.Errorf("the format of import id must be 'registry:namespace:name'")
	}
	if err := d.Set("registry", items[0]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	if err := d.Set("namespace", items[1]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	if err := d.Set("name", items[2]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}

func ResourceVestackCrRepository() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackCrRepositoryCreate,
		Read:   resourceVestackCrRepositoryRead,
		Update: resourceVestackCrRepositoryUpdate,
		Delete: resourceVestackCrRepositoryDelete,
		Importer: &schema.ResourceImporter{
			State: crRepositoryImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"registry": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The CrRegistry name.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target namespace name.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of CrRepository.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of CrRepository.",
			},
			"access_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Private",
				ValidateFunc: validation.StringInSlice([]string{"Private", "Public"}, false),
				Description:  "The access level of CrRepository.",
			},
		},
	}
	dataSource := DataSourceVestackCrRepositories().Schema["repositories"].Elem.(*schema.Resource).Schema
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceVestackCrRepositoryCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRepositoryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceVestackCrRepository())
	if err != nil {
		return fmt.Errorf("error on creating CrRepository %q,%s", d.Id(), err)
	}
	return resourceVestackCrRepositoryRead(d, meta)
}

func resourceVestackCrRepositoryUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRepositoryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceVestackCrRepository())
	if err != nil {
		return fmt.Errorf("error on updating CrRepository  %q, %s", d.Id(), err)
	}
	return resourceVestackCrRepositoryRead(d, meta)
}

func resourceVestackCrRepositoryDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRepositoryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceVestackCrRepository())
	if err != nil {
		return fmt.Errorf("error on deleting CrRepository %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackCrRepositoryRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRepositoryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceVestackCrRepository())
	if err != nil {
		return fmt.Errorf("Error on reading CrRepository %q,%s", d.Id(), err)
	}
	return err
}
