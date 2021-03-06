---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "servicenow_ui_macro Resource - terraform-provider-servicenow"
subcategory: ""
description: |-
  servicenow_ui_macro manages a UI Macro configuration within ServiceNow.
---

# servicenow_ui_macro (Resource)

`servicenow_ui_macro` manages a UI Macro configuration within ServiceNow.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String)
- **xml** (String) The body of the UI Macro. Must be in XML format.

### Optional

- **active** (Boolean) Whether or not this Macro is enabled.
- **api_name** (String) The scoped name of the macro. Normally contains the name field prefixed with the application scope.
- **description** (String)
- **id** (String) The ID of this resource.
- **protection_policy** (String) Determines how application files are protected when downloaded or installed. Can be empty for no protection, 'read' for read-only protection or 'protected'.
- **scope** (String) Associates a resource to a specific application ID in ServiceNow.


