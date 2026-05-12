# Look up an existing UI policy in ServiceNow by its short description.
data "servicenow_ui_policy" "example" {
  short_description = "Make assignment_group mandatory when priority is critical"
}
