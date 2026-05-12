# Look up an existing choice list value in ServiceNow by table name, field name and value.
data "servicenow_choice" "example" {
  name    = "incident"
  element = "state"
  value   = "1"
}
