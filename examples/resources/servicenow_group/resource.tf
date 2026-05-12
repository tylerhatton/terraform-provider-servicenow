# Manages a user group within ServiceNow.
resource "servicenow_group" "example" {
  name        = "Example Engineering Group"
  description = "Engineering team responsible for example services."
  email       = "engineering@example.com"
  active      = true
}
