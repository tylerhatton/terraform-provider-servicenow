---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "servicenow_rest_method Resource - terraform-provider-servicenow"
subcategory: ""
description: |-
  servicenow_rest_method manages a REST method within ServiceNow.
---

# servicenow_rest_method (Resource)

`servicenow_rest_method` manages a REST method within ServiceNow.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **http_method** (String) The HTTP method this record implements. Can be 'get', 'post', 'put', 'patch' or 'delete'.
- **name** (String) A unique identifier for this HTTP method record.
- **rest_message_id** (String) The REST message record ID this method is based on.

### Optional

- **id** (String) The ID of this resource.
- **rest_endpoint** (String) The URL of the REST web service provider this method sends requests to. Can contain variables in the format '${variable}'.
- **scope** (String) Associates a resource to a specific application ID in ServiceNow.

### Read-Only

- **qualified_name** (String)


