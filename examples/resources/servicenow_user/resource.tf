# Manages a user record within ServiceNow.
resource "servicenow_user" "example" {
  user_name  = "jane.doe"
  first_name = "Jane"
  last_name  = "Doe"
  email      = "jane.doe@example.com"
  title      = "Software Engineer"
  active     = true
}
