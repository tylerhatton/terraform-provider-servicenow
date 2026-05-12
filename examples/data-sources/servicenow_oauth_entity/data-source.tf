# Look up an existing OAuth application entity in ServiceNow by name.
data "servicenow_oauth_entity" "example" {
  name = "my-oauth-entity"
}
