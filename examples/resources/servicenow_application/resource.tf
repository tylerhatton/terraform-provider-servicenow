# Manages a scoped application in ServiceNow.
resource "servicenow_application" "example" {
  name    = "Example Application"
  scope   = "x_example_app"
  version = "1.0.0"
}
