package tos_bucket_website

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
TosBucketWebsite can be imported using the bucketName, e.g.
```
$ terraform import vestack_tos_bucket_website.default bucket_name
```

*/

func ResourceVestackTosBucketWebsite() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackTosBucketWebsiteCreate,
		Read:   resourceVestackTosBucketWebsiteRead,
		Update: resourceVestackTosBucketWebsiteUpdate,
		Delete: resourceVestackTosBucketWebsiteDelete,
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
			"index_document": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The index document configuration for the website.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"suffix": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The suffix of the index document, e.g., index.html.",
						},
						"support_sub_dir": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to support subdirectory indexing. Default is false.",
						},
					},
				},
			},
			"error_document": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The error document configuration for the website.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The key of the error document object, e.g., error.html.",
						},
					},
				},
			},
			"redirect_all_requests_to": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The redirect configuration for all requests.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The target host name for redirect.",
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "The protocol for redirect. Valid values: http, https.",
							ValidateFunc: validation.StringInSlice([]string{"http", "https"}, false),
						},
					},
				},
			},
			"routing_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The routing rules for the website.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "The condition for the routing rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_prefix_equals": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The key prefix that must match for the rule to apply.",
									},
									"http_error_code_returned_equals": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The HTTP error code that must match for the rule to apply, e.g., 404.",
									},
								},
							},
						},
						"redirect": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "The redirect configuration for the routing rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "The protocol to use for the redirect. Valid values: http, https.",
										ValidateFunc: validation.StringInSlice([]string{"http", "https"}, false),
									},
									"host_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The host name to redirect to.",
									},
									"replace_key_with": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The key to replace the original key with.",
									},
									"replace_key_prefix_with": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The key prefix to replace the original key prefix with.",
									},
									"http_redirect_code": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The HTTP redirect code to use, e.g., 301, 302.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return resource
}

func resourceVestackTosBucketWebsiteCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketWebsiteService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackTosBucketWebsite())
	if err != nil {
		return fmt.Errorf("error on creating tos_bucket_website %q, %s", d.Id(), err)
	}
	return resourceVestackTosBucketWebsiteRead(d, meta)
}

func resourceVestackTosBucketWebsiteRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketWebsiteService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackTosBucketWebsite())
	if err != nil {
		return fmt.Errorf("error on reading tos_bucket_website %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackTosBucketWebsiteUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketWebsiteService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVestackTosBucketWebsite())
	if err != nil {
		return fmt.Errorf("error on updating tos_bucket_website %q, %s", d.Id(), err)
	}
	return resourceVestackTosBucketWebsiteRead(d, meta)
}

func resourceVestackTosBucketWebsiteDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketWebsiteService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackTosBucketWebsite())
	if err != nil {
		return fmt.Errorf("error on deleting tos_bucket_website %q, %s", d.Id(), err)
	}
	return nil
}
