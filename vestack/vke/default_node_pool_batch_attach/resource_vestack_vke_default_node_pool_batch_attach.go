package default_node_pool_batch_attach

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	bp "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/vestack/vke/default_node_pool"
)

/*

The resource not support import

*/

func ResourceVestackDefaultNodePoolBatchAttach() *schema.Resource {
	m := map[string]*schema.Schema{
		"cluster_id": default_node_pool.ResourceVestackDefaultNodePool().Schema["cluster_id"],
		"default_node_pool_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The default NodePool ID.",
		},
		"instances": default_node_pool.ResourceVestackDefaultNodePool().Schema["instances"],
		"kubernetes_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"labels": {
						Type:     schema.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:        schema.TypeString,
									Required:    true,
									ForceNew:    true,
									Description: "The Key of Labels.",
								},
								"value": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    true,
									Description: "The Value of Labels.",
								},
							},
						},
						Description: "The Labels of KubernetesConfig.",
					},
					"taints": {
						Type:     schema.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:        schema.TypeString,
									Required:    true,
									ForceNew:    true,
									Description: "The Key of Taints.",
								},
								"value": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    true,
									Description: "The Value of Taints.",
								},
								"effect": {
									Type:     schema.TypeString,
									Optional: true,
									ForceNew: true,
									ValidateFunc: validation.StringInSlice([]string{
										"NoSchedule",
										"NoExecute",
										"PreferNoSchedule",
									}, false),
									Description: "The Effect of Taints. The value can be one of the following: `NoSchedule`, `NoExecute`, `PreferNoSchedule`, default value is `NoSchedule`.",
								},
							},
						},
						Description: "The Taints of KubernetesConfig.",
					},
					"cordon": {
						Type:        schema.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Description: "The Cordon of KubernetesConfig.",
					},
				},
			},
			Description: "The KubernetesConfig of NodeConfig. Please note that this field is the configuration of the node. The same key is subject to the config of the node pool. Different keys take effect together.",
		},
	}
	bp.MergeDateSourceToResource(default_node_pool.ResourceVestackDefaultNodePool().Schema, &m)

	// logger.Debug(logger.RespFormat, "ATTACH_TEST", m)

	return &schema.Resource{
		Create: resourceVestackDefaultNodePoolBatchAttachCreate,
		Update: resourceVestackDefaultNodePoolBatchAttachUpdate,
		Read:   resourceVestackDefaultNodePoolBatchAttachUpdate,
		Delete: resourceVestackNodePoolBatchAttachDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				return nil, fmt.Errorf("The resource not support import ")
			},
		},
		Schema: m,
	}
}

func resourceVestackDefaultNodePoolBatchAttachCreate(d *schema.ResourceData, meta interface{}) (err error) {
	nodePoolService := NewVestackVkeDefaultNodePoolBatchAttachService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(nodePoolService, d, ResourceVestackDefaultNodePoolBatchAttach())
	if err != nil {
		return fmt.Errorf("error on creating DefaultNodePoolBatchAttach  %q, %w", d.Id(), err)
	}
	return resourceVestackDefaultNodePoolBatchAttachRead(d, meta)
}

func resourceVestackDefaultNodePoolBatchAttachRead(d *schema.ResourceData, meta interface{}) (err error) {
	nodePoolService := NewVestackVkeDefaultNodePoolBatchAttachService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(nodePoolService, d, ResourceVestackDefaultNodePoolBatchAttach())
	if err != nil {
		return fmt.Errorf("error on reading DefaultNodePoolBatchAttach %q, %w", d.Id(), err)
	}
	return err
}

func resourceVestackDefaultNodePoolBatchAttachUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	nodePoolService := NewVestackVkeDefaultNodePoolBatchAttachService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(nodePoolService, d, ResourceVestackDefaultNodePoolBatchAttach())
	if err != nil {
		return fmt.Errorf("error on updating DefaultNodePoolBatchAttach  %q, %w", d.Id(), err)
	}
	return resourceVestackDefaultNodePoolBatchAttachRead(d, meta)
}

func resourceVestackNodePoolBatchAttachDelete(d *schema.ResourceData, meta interface{}) (err error) {
	nodePoolService := NewVestackVkeDefaultNodePoolBatchAttachService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(nodePoolService, d, ResourceVestackDefaultNodePoolBatchAttach())
	if err != nil {
		return fmt.Errorf("error on deleting DefaultNodePoolBatchAttach %q, %w", d.Id(), err)
	}
	return err
}
