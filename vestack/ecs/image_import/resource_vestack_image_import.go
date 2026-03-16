package image_import

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
ImageImport can be imported using the id, e.g.
```
$ terraform import vestack_image_import.default resource_id
```

*/

func ResourceVestackImageImport() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackImageImportCreate,
		Read:   resourceVestackImageImportRead,
		Update: resourceVestackImageImportUpdate,
		Delete: resourceVestackImageImportDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"platform": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The platform of the custom image. Valid values: `CentOS`, `Debian`, `veLinux`, `Windows Server`, `Fedora`, `OpenSUSE`, `Ubuntu`, `Rocky Linux`, `AlmaLinux`.",
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The url of the custom image in tos bucket." +
					"When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"image_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the custom image.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the custom image.",
			},
			"architecture": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The architecture of the custom image. Valid values: `amd64`, `arm64`.",
			},
			"boot_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The boot mode of the custom image. Valid values: `BIOS`, `UEFI`.",
			},
			"license_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The license type of the custom image. Valid values: `VolcanoEngine`.",
			},
			"os_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The os type of the custom image. Valid values: `linux`, `Windows`.",
			},
			"platform_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The platform version of the custom image.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project name of the custom image.",
			},
			"tags": bp.TagsSchema(),

			// computed fields
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of Image.",
			},
			"visibility": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The visibility of Image.",
			},
			"os_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of Image operating system.",
			},
			"is_support_cloud_init": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the Image support cloud-init.",
			},
			"share_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The share mode of Image.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size(GiB) of Image.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create time of Image.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time of Image.",
			},
		},
	}
	return resource
}

func resourceVestackImageImportCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewImageImportService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackImageImport())
	if err != nil {
		return fmt.Errorf("error on creating image_import %q, %s", d.Id(), err)
	}
	return resourceVestackImageImportRead(d, meta)
}

func resourceVestackImageImportRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewImageImportService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackImageImport())
	if err != nil {
		return fmt.Errorf("error on reading image_import %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackImageImportUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewImageImportService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackImageImport())
	if err != nil {
		return fmt.Errorf("error on updating image_import %q, %s", d.Id(), err)
	}
	return resourceVestackImageImportRead(d, meta)
}

func resourceVestackImageImportDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewImageImportService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackImageImport())
	if err != nil {
		return fmt.Errorf("error on deleting image_import %q, %s", d.Id(), err)
	}
	return err
}
