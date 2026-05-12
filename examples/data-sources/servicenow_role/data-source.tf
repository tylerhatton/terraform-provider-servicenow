# Look up an existing role in ServiceNow by its suffix.
data "servicenow_role" "example" {
  suffix = "admin"
}
