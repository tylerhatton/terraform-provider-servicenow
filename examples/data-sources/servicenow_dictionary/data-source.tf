# Look up an existing column definition in ServiceNow by table name and field name.
data "servicenow_dictionary" "example" {
  name    = "sys_user"
  element = "user_name"
}
