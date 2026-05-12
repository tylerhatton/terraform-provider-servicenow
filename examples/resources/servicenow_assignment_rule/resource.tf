# Manages an assignment rule (sysrule_assignment) in ServiceNow.
resource "servicenow_assignment_rule" "example" {
  name             = "Network Incidents To Network Group"
  table            = "incident"
  description      = "Routes network category incidents to the Network group."
  active           = true
  order            = 100
  condition        = "category=network^EQ"
  assignment_group = "287ebd7da9fe198100f92cc8d1d2154e"
}
