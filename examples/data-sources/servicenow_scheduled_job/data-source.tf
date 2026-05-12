# Look up an existing scheduled job in ServiceNow by name.
data "servicenow_scheduled_job" "example" {
  name = "Example Daily Job"
}
