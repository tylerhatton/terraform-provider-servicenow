# Look up an existing email notification in ServiceNow by name.
data "servicenow_notification" "example" {
  name = "Example Incident Notification"
}
