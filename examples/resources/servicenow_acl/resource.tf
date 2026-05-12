# Manages a security access control rule (sys_security_acl) in ServiceNow.
resource "servicenow_db_table" "example" {
  label     = "Example Table"
  user_role = ""
}

resource "servicenow_acl" "read" {
  name            = "${servicenow_db_table.example.name}.*"
  operation       = "read"
  type            = "record"
  admin_overrides = true
  active          = true
  description     = "Grant read access to records on the example table."
}

resource "servicenow_acl" "write" {
  name            = "${servicenow_db_table.example.name}.*"
  operation       = "write"
  type            = "record"
  admin_overrides = true
  active          = true
  description     = "Grant write access to records on the example table."
}
