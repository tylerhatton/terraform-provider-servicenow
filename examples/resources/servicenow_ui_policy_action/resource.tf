# Manages an action attached to a UI policy in ServiceNow.
resource "servicenow_ui_policy" "example" {
  short_description = "Make assignment_group mandatory when priority is critical"
  table             = "incident"
  conditions        = "priority=1^EQ"
}

resource "servicenow_ui_policy_action" "example" {
  ui_policy  = servicenow_ui_policy.example.id
  field_name = "assignment_group"
  visible    = "true"
  mandatory  = "true"
  disabled   = "ignore"
}
