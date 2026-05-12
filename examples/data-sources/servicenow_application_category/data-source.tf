# Look up an existing application category in ServiceNow by name.
data "servicenow_application_category" "example" {
  name = "Custom Applications"
}
