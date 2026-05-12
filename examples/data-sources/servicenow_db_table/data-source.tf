# Look up an existing database table in ServiceNow by internal name.
data "servicenow_db_table" "example" {
  name = "sys_user"
}
