# Manages a field mapping (sys_transform_entry) inside a transform map.
resource "servicenow_transform_map" "parent" {
  name         = "Example User Import"
  source_table = "u_user_import"
  target_table = "sys_user"
}

resource "servicenow_transform_entry" "user_name" {
  map          = servicenow_transform_map.parent.id
  source_field = "u_email"
  target_field = "user_name"
  coalesce     = true
}
