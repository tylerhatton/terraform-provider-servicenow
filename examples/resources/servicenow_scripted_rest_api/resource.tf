# Manages a scripted REST API in ServiceNow.
resource "servicenow_scripted_rest_api" "example" {
  name       = "Example API"
  service_id = "example_api"
  active     = true
}
