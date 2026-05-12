# Manages a service catalog in ServiceNow.
resource "servicenow_service_catalog" "example" {
  title       = "Example Catalog"
  description = "An example service catalog managed by Terraform."
  active      = true
}
