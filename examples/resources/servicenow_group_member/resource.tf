# Adds a user as a member of a group.
resource "servicenow_user" "example" {
  user_name  = "jane.doe"
  first_name = "Jane"
  last_name  = "Doe"
}

resource "servicenow_group" "example" {
  name        = "Example Engineering Group"
  description = "Engineering team responsible for example services."
}

resource "servicenow_group_member" "example" {
  user  = servicenow_user.example.id
  group = servicenow_group.example.id
}
