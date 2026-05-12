# Manages an application menu section in the ServiceNow application navigator.
resource "servicenow_application_menu" "example" {
  title  = "Example Menu"
  order  = 100
  active = true
}
