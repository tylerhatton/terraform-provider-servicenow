---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "servicenow_oauth_entity Resource - terraform-provider-servicenow"
subcategory: ""
description: |-
  servicenow_js_include manages a javascript script within ServiceNow.
---

# servicenow_oauth_entity (Resource)

`servicenow_js_include` manages a javascript script within ServiceNow.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) Name of the OAuth app.

### Optional

- **access** (String) Whether this Script can be accessed from only this application scope or all application scopes. Values can be 'package_private' or 'public'.
- **access_token_lifespan** (Number) Number of seconds a newly created access token is good for.
- **id** (String) The ID of this resource.
- **logo_url** (String)
- **redirect_url** (String) The OAuth app's endpoint to receive the authorization code.
- **refresh_token_lifespan** (Number) Number of seconds the refresh token is good for.
- **scope** (String) Associates a resource to a specific application ID in ServiceNow.

### Read-Only

- **client_id** (String) OAuth Client ID required during handshake.
- **client_uuid** (String) Internal unique identifier of the entity.


