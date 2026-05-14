# Manages a column definition (sys_dictionary entry) for a table in ServiceNow.
resource "servicenow_db_table" "example" {
  label     = "Example Dictionary Table"
  user_role = ""
}

resource "servicenow_dictionary" "example" {
  name          = servicenow_db_table.example.name
  element       = "u_example_field"
  column_label  = "Example Field"
  internal_type = "string"
  max_length    = 255
  mandatory     = false
  active        = true
}
