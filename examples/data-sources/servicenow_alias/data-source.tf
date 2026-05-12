# Look up an existing connection or credential alias in ServiceNow by name.
data "servicenow_alias" "example" {
  name = "my-existing-alias"
}
