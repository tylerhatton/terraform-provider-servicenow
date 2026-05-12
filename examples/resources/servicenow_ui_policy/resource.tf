# Manages a UI policy in ServiceNow to dynamically change form field behavior.
resource "servicenow_ui_policy" "example" {
  short_description = "Make assignment_group mandatory when priority is critical"
  table             = "incident"
  active            = true
  conditions        = "priority=1^EQ"
  reverse_if_false  = true
  on_load           = true
  order             = 100
}
