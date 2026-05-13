package health_check_log_topic

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
HealthCheckLogTopic can be imported using the id, e.g.
```
$ terraform import vestack_health_check_log_topic.default log_topic_id:load_balancer_id
```

*/

func ResourceVestackHealthCheckLogTopic() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVestackHealthCheckLogTopicCreate,
		Read:   resourceVestackHealthCheckLogTopicRead,
		Delete: resourceVestackHealthCheckLogTopicDelete,
		Importer: &schema.ResourceImporter{State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			parts := strings.Split(d.Id(), ":")
			if len(parts) < 2 {
				return []*schema.ResourceData{d}, fmt.Errorf("import id must be 'log_topic_id:load_balancer_id'")
			}
			_ = d.Set("log_topic_id", parts[0])
			_ = d.Set("load_balancer_id", parts[1])
			return []*schema.ResourceData{d}, nil
		}},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"log_topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the log topic.",
			},
			"load_balancer_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the CLB instance.",
			},
		},
	}
	return resource
}

func resourceVestackHealthCheckLogTopicCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewHealthCheckLogTopicService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVestackHealthCheckLogTopic())
	if err != nil {
		return fmt.Errorf("error on creating health_check_log_topic %q, %s", d.Id(), err)
	}
	return resourceVestackHealthCheckLogTopicRead(d, meta)
}

func resourceVestackHealthCheckLogTopicRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewHealthCheckLogTopicService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVestackHealthCheckLogTopic())
	if err != nil {
		return fmt.Errorf("error on reading health_check_log_topic %q, %s", d.Id(), err)
	}
	return err
}

func resourceVestackHealthCheckLogTopicDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewHealthCheckLogTopicService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVestackHealthCheckLogTopic())
	if err != nil {
		return fmt.Errorf("error on deleting health_check_log_topic %q, %s", d.Id(), err)
	}
	return err
}
