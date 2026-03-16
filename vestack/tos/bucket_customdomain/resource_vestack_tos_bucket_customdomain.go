package tos_bucket_customdomain

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
TosBucketCustomDomain can be imported using the bucketName:domain, e.g.
```
$ terraform import vestack_tos_bucket_customdomain.default bucket_name:custom_domain
```

*/

func ResourceVestackTosBucketCustomDomain() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackTosBucketCustomDomainCreate,
		Read:   resourceVestackTosBucketCustomDomainRead,
		Update: resourceVestackTosBucketCustomDomainUpdate,
		Delete: resourceVestackTosBucketCustomDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the TOS bucket.",
			},
			"custom_domain_rule": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The custom domain role for the bucket.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The custom domain name for the bucket.",
						},
						"cert_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The certificate id.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom domain access protocol.tos|s3.",
						},
					},
				},
			},
		},
	}
	return resource
}

func resourceVestackTosBucketCustomDomainCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketCustomDomainService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackTosBucketCustomDomain())
	if err != nil {
		return fmt.Errorf("error on creating tos_bucket_customdomain %q, %s", d.Id(), err)
	}
	return resourceVestackTosBucketCustomDomainRead(d, meta)
}

func resourceVestackTosBucketCustomDomainRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketCustomDomainService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackTosBucketCustomDomain())
	if err != nil {
		return fmt.Errorf("error on reading tos_bucket_customdomain %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackTosBucketCustomDomainUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketCustomDomainService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackTosBucketCustomDomain())
	if err != nil {
		return fmt.Errorf("error on updating tos_bucket_customdomain %q, %s", d.Id(), err)
	}
	return resourceVestackTosBucketCustomDomainRead(d, meta)
}

func resourceVestackTosBucketCustomDomainDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketCustomDomainService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackTosBucketCustomDomain())
	if err != nil {
		return fmt.Errorf("error on deleting tos_bucket_customdomain %q, %s", d.Id(), err)
	}
	return nil
}

func parseImportId(id string) (string, string, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid import ID format, expected format: bucket_name:custom_domain")
	}
	return parts[0], parts[1], nil
}
