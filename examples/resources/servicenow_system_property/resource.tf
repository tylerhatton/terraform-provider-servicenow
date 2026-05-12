# Manages a system property (sys_properties) in ServiceNow.
resource "servicenow_system_property" "example" {
  suffix      = "example.property"
  type        = "string"
  description = "Example system property managed by Terraform."
}
