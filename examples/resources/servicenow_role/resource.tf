# Manages a role within ServiceNow.
resource "servicenow_role" "example" {
  suffix      = "example_role"
  description = "Example role managed by Terraform."
}
