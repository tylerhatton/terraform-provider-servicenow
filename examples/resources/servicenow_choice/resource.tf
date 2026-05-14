# Manages a choice list value (sys_choice entry) for a column on a ServiceNow table.
resource "servicenow_db_table" "example" {
  label     = "Example Choice Table"
  user_role = ""
}

resource "servicenow_dictionary" "status" {
  name          = servicenow_db_table.example.name
  element       = "u_status"
  column_label  = "Status"
  internal_type = "string"
  max_length    = 40
  choice        = 3
}

resource "servicenow_choice" "open" {
  name     = servicenow_dictionary.status.name
  element  = servicenow_dictionary.status.element
  value    = "open"
  label    = "Open"
  sequence = 100
}

resource "servicenow_choice" "closed" {
  name     = servicenow_dictionary.status.name
  element  = servicenow_dictionary.status.element
  value    = "closed"
  label    = "Closed"
  sequence = 200
}
