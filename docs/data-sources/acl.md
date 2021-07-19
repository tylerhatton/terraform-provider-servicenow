---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "servicenow_acl Data Source - terraform-provider-servicenow"
subcategory: ""
description: |-
  
---

# servicenow_acl (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) Enter the name of the object being secured, either the record name or the table and field names.

### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **active** (Boolean) Activates the ACL rule.
- **admin_overrides** (Boolean) Users with admin override this rule
- **advanced** (Boolean) Displays the Script field when active.
- **condition** (String) Selects the fields and values that must be true for users to access the object.
- **description** (String) Enter a description of the object or permissions this ACL rule secures.
- **operation** (String) Select the operation this ACL rule secures.
- **protection_policy** (String) Determines how application files are protected when downloaded or installed. Can be empty for no protection, 'read' for read-only protection or 'protected'.
- **scope** (String) Associates a resource to a specific application ID in ServiceNow.
- **script** (String) Custom script describing the permissions required to access the object.
- **type** (String) Select what kind of object this ACL rule secures.

