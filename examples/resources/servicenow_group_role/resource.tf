# Assigns a role to a group so all members inherit it.
resource "servicenow_group" "example" {
  name        = "Example Engineering Group"
  description = "Engineering team responsible for example services."
}

resource "servicenow_role" "example" {
  suffix      = "example_role"
  description = "Example role managed by Terraform."
}

resource "servicenow_group_role" "example" {
  group    = servicenow_group.example.id
  role     = servicenow_role.example.id
  inherits = true
}
