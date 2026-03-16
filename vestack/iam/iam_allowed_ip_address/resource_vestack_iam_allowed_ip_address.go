package iam_allowed_ip_address

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*
Import
Iam AllowedIpAddress key don't support import
*/
func ResourceVestackIamAllowedIpAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackIamAllowedIpAddressCreate,
		Read:   resourceVestackIamAllowedIpAddressRead,
		Update: resourceVestackIamAllowedIpAddressUpdate,
		Delete: resourceVestackIamAllowedIpAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"enable_ip_list": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable the IP whitelist.",
			},
			"ip_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The IP whitelist list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The IP address.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The description of the IP address.",
						},
					},
				},
			},
		},
	}
}

func resourceVestackIamAllowedIpAddressCreate(d *schema.ResourceData, meta interface{}) error {
	service := NewIamAllowedIpAddressService(meta.(*bp.SdkClient))
	err := bp.DefaultDispatcher().Create(service, d, ResourceVestackIamAllowedIpAddress())
	if err != nil {
		return fmt.Errorf("error on creating allowed ip addresses: %s", err)
	}
	return resourceVestackIamAllowedIpAddressRead(d, meta)
}

func resourceVestackIamAllowedIpAddressRead(d *schema.ResourceData, meta interface{}) error {
	service := NewIamAllowedIpAddressService(meta.(*bp.SdkClient))
	err := bp.DefaultDispatcher().Read(service, d, ResourceVestackIamAllowedIpAddress())
	if err != nil {
		return fmt.Errorf("error on reading allowed ip addresses: %s", err)
	}
	return nil
}

func resourceVestackIamAllowedIpAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	service := NewIamAllowedIpAddressService(meta.(*bp.SdkClient))
	err := bp.DefaultDispatcher().Update(service, d, ResourceVestackIamAllowedIpAddress())
	if err != nil {
		return fmt.Errorf("error on updating allowed ip addresses: %s", err)
	}
	return resourceVestackIamAllowedIpAddressRead(d, meta)
}

func resourceVestackIamAllowedIpAddressDelete(d *schema.ResourceData, meta interface{}) error {
	service := NewIamAllowedIpAddressService(meta.(*bp.SdkClient))
	err := bp.DefaultDispatcher().Delete(service, d, ResourceVestackIamAllowedIpAddress())
	if err != nil {
		return fmt.Errorf("error on deleting allowed ip addresses: %s", err)
	}
	return nil
}
