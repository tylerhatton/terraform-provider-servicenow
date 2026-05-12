# Look up an existing user group in ServiceNow by its name.
data "servicenow_group" "example" {
  name = "Application Development"
}
