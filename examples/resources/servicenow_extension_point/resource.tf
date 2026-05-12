# Manages a scripted extension point in ServiceNow.
resource "servicenow_extension_point" "example" {
  name        = "Example Extension Point"
  description = "An example scripted extension point."
}
