# Assigns a role directly to a user.
resource "servicenow_user" "example" {
  user_name  = "jane.doe"
  first_name = "Jane"
  last_name  = "Doe"
}

resource "servicenow_role" "example" {
  suffix      = "example_role"
  description = "Example role managed by Terraform."
}

resource "servicenow_user_role" "example" {
  user = servicenow_user.example.id
  role = servicenow_role.example.id
}
