package iam_policy_project

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	bp "github.com/volcengine/terraform-provider-vestack/common"
)

/*

Import
IamPolicyProject can be imported using the id, e.g.
```
$ terraform import vestack_iam_policy_project.default PrincipalType:PrincipalName:PolicyType:PolicyName:ProjectName
```

*/

func ResourceVestackIamPolicyProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceVestackIamPolicyProjectCreate,
		Read:   resourceVestackIamPolicyProjectRead,
		Delete: resourceVestackIamPolicyProjectDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				id := d.Id()
				// Remove trailing commas or spaces first
				id = strings.TrimRight(id, ", ")
				ids := strings.Split(id, ":")
				if len(ids) < 5 {
					return nil, fmt.Errorf("invalid id format, should be PrincipalType:PrincipalName:PolicyType:PolicyName:ProjectName1,ProjectName2...")
				}
				_ = d.Set("principal_type", ids[0])
				_ = d.Set("principal_name", ids[1])
				_ = d.Set("policy_type", ids[2])
				_ = d.Set("policy_name", ids[3])

				rawProjects := strings.Split(ids[4], ",")
				var projects []interface{}
				var cleanProjectNames []string
				for _, p := range rawProjects {
					p = strings.TrimSpace(p)
					if p != "" {
						projects = append(projects, p)
						cleanProjectNames = append(cleanProjectNames, p)
					}
				}
				if len(projects) > 0 {
					_ = d.Set("project_names", schema.NewSet(schema.HashString, projects))
				}

				// Reset ID to a clean format
				newId := fmt.Sprintf("%s:%s:%s:%s:%s", ids[0], ids[1], ids[2], ids[3], strings.Join(cleanProjectNames, ","))
				d.SetId(newId)

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"principal_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the principal. Valid values: User, Role, UserGroup.",
			},
			"principal_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the principal.",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the policy. Valid values: System, Custom.",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the policy.",
			},
			"project_names": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Description: "The list of project names, which is the scope of the policy.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceVestackIamPolicyProjectCreate(d *schema.ResourceData, meta interface{}) error {
	service := NewIamPolicyProjectService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Create(service, d, ResourceVestackIamPolicyProject())
}

func resourceVestackIamPolicyProjectRead(d *schema.ResourceData, meta interface{}) error {
	service := NewIamPolicyProjectService(meta.(*bp.SdkClient))
	err := bp.DefaultDispatcher().Read(service, d, ResourceVestackIamPolicyProject())
	if err != nil {
		if bp.ResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	return nil
}

func resourceVestackIamPolicyProjectDelete(d *schema.ResourceData, meta interface{}) error {
	service := NewIamPolicyProjectService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Delete(service, d, ResourceVestackIamPolicyProject())
}
