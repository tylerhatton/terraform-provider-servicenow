# Look up an existing assignment rule in ServiceNow by name.
data "servicenow_assignment_rule" "example" {
  name = "Network Incidents To Network Group"
}
