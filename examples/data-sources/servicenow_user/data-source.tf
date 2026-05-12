# Look up an existing user in ServiceNow by their user_name (login).
data "servicenow_user" "example" {
  user_name = "admin"
}
