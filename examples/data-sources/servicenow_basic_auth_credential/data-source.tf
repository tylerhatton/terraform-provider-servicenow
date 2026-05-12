# Look up an existing basic auth credential in ServiceNow by name.
data "servicenow_basic_auth_credential" "example" {
  name = "my-credential"
}
