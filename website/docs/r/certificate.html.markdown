---
subcategory: "CLB"
layout: "vestack"
page_title: "Vestack: vestack_certificate"
sidebar_current: "docs-vestack-resource-certificate"
description: |-
  Provides a resource to manage certificate
---
# vestack_certificate
Provides a resource to manage certificate
## Example Usage
```hcl

```
## Argument Reference
The following arguments are supported:
* `private_key` - (Required, ForceNew) The private key of the Certificate. When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.
* `public_key` - (Required, ForceNew) The public key of the Certificate. When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.
* `certificate_name` - (Optional) The name of the Certificate.
* `description` - (Optional) The description of the Certificate.
* `project_name` - (Optional, ForceNew) The ProjectName of the Certificate.
* `tags` - (Optional) Tags.

The `tags` object supports the following:

* `key` - (Required) The Key of Tags.
* `value` - (Required) The Value of Tags.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
Certificate can be imported using the id, e.g.
```
$ terraform import vestack_certificate.default cert-2fe5k****c16o5oxruvtk3qf5
```

