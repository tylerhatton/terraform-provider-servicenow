---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "servicenow_rest_method_header Resource - terraform-provider-servicenow"
subcategory: ""
description: |-
  servicenow_rest_message_header manages a header to be applied to a REST method within ServiceNow.
---

# servicenow_rest_method_header (Resource)

`servicenow_rest_message_header` manages a header to be applied to a REST method within ServiceNow.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) The name of the header to add to the HTTP request.
- **rest_method_id** (String) The REST method record ID this header will be applied to.
- **value** (String) The value of the header to add to the HTTP request.

### Optional

- **id** (String) The ID of this resource.
- **scope** (String) Associates a resource to a specific application ID in ServiceNow.


