# Look up an existing application in ServiceNow by name.
data "servicenow_application" "example" {
  name = "Global"
}
