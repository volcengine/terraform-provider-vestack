package security_group_rule

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
SecurityGroupRule can be imported using the id, e.g.
```
$ terraform import vestack_security_group_rule.default ID is a string concatenated with colons(SecurityGroupId:Protocol:PortStart:PortEnd:CidrIp:SourceGroupId:Direction:Policy:Priority)
```

*/

func ResourceVestackSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackSecurityGroupRuleCreate,
		Read:   resourceVestackSecurityGroupRuleRead,
		Update: resourceVestackSecurityGroupRuleUpdate,
		Delete: resourceVestackSecurityGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: importSecurityGroupRule,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ingress",
					"egress",
				}, false),
				Description: "Direction of rule, ingress (inbound) or egress (outbound).",
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"tcp",
					"udp",
					"icmp",
					"all",
					"icmpv6",
				}, false),
				Description: "Protocol of the SecurityGroup, the value can be `tcp` or `udp` or `icmp` or `all` or `icmpv6`.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of SecurityGroup.",
			},
			"port_start": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(-1, 65535),
				Description:  "Port start of egress/ingress Rule.",
			},
			"port_end": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(-1, 65535),
				Description:  "Port end of egress/ingress Rule.",
			},
			"cidr_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_group_id"},
				Description:   "Cidr ip of egress/ingress Rule.",
			},
			"source_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cidr_ip"},
				Description:   "ID of the source security group whose access permission you want to set.",
			},
			"policy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"accept",
					"drop",
				}, false),
				Default:     "accept",
				Description: "Access strategy.",
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 100),
				Description:  "Priority of a security group rule.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "description of a egress rule.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of SecurityGroup.",
			},
		},
	}
}

func importSecurityGroupRule(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(d.Id(), ":")
	if len(items) != 9 {
		return []*schema.ResourceData{d}, fmt.Errorf("import id must be of the form " +
			"SecurityGroupId:Protocol:PortStart:PortEnd:CidrIp:SourceGroupId:Direction:Policy:Priority")
	}
	err = d.Set("security_group_id", items[0])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	err = d.Set("protocol", items[1])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}

	if len(items[2]) > 0 {
		ps, err := strconv.Atoi(items[2])
		if err != nil {
			return []*schema.ResourceData{d}, err
		}
		err = d.Set("port_start", ps)
		if err != nil {
			return []*schema.ResourceData{d}, err
		}
	}

	if len(items[3]) > 0 {
		pn, err := strconv.Atoi(items[3])
		if err != nil {
			return []*schema.ResourceData{d}, err
		}
		err = d.Set("port_end", pn)
		if err != nil {
			return []*schema.ResourceData{d}, err
		}
	}

	err = d.Set("cidr_ip", items[4])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}

	err = d.Set("source_group_id", items[5])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}

	err = d.Set("direction", items[6])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}

	err = d.Set("policy", items[7])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}

	if len(items[8]) > 0 {
		pr, err := strconv.Atoi(items[8])
		if err != nil {
			return []*schema.ResourceData{d}, err
		}
		err = d.Set("priority", pr)
		if err != nil {
			return []*schema.ResourceData{d}, err
		}
	}
	return []*schema.ResourceData{d}, nil
}

func resourceVestackSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) (err error) {
	securityGroupRuleService := NewSecurityGroupRuleService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(securityGroupRuleService, d, ResourceVestackSecurityGroupRule())
	if err != nil {
		return fmt.Errorf("error on creating securityGroupRuleService  %q, %w", d.Id(), err)
	}
	return resourceVestackSecurityGroupRuleRead(d, meta)
}

func resourceVestackSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) (err error) {
	securityGroupRuleService := NewSecurityGroupRuleService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(securityGroupRuleService, d, ResourceVestackSecurityGroupRule())
	if err != nil {
		return fmt.Errorf("error on reading securityGroupRuleService %q, %w", d.Id(), err)
	}
	return err
}

func resourceVestackSecurityGroupRuleUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	securityGroupRuleService := NewSecurityGroupRuleService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(securityGroupRuleService, d, ResourceVestackSecurityGroupRule())
	if err != nil {
		return fmt.Errorf("error on updating securityGroupRuleService  %q, %w", d.Id(), err)
	}
	return resourceVestackSecurityGroupRuleRead(d, meta)
}

func resourceVestackSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) (err error) {
	securityGroupRuleService := NewSecurityGroupRuleService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(securityGroupRuleService, d, ResourceVestackSecurityGroupRule())
	if err != nil {
		return fmt.Errorf("error on deleting securityGroupRuleService %q, %w", d.Id(), err)
	}
	return err
}
