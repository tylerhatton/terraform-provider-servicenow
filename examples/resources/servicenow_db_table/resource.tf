# Manages a database table within ServiceNow.
resource "servicenow_db_table" "example" {
  label     = "Example Table"
  user_role = ""
}
