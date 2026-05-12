# Manages a transform map (sys_transform_map) in ServiceNow.
resource "servicenow_transform_map" "example" {
  name                     = "Example User Import"
  source_table             = "u_user_import"
  target_table             = "sys_user"
  active                   = true
  run_business_rules       = true
  enforce_mandatory_fields = "no"
  copy_empty_fields        = false
  order                    = 100
  description              = "Imports staged user rows into sys_user."
}
