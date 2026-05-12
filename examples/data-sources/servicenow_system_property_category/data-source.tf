# Look up an existing system property category in ServiceNow by name.
data "servicenow_system_property_category" "example" {
  name = "System"
}
