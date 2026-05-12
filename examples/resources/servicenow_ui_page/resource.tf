# Manages a UI Page (HTML + client/server scripts) in ServiceNow.
resource "servicenow_ui_page" "example" {
  name              = "example_page"
  category          = "general"
  html              = "<h1>Example Page</h1>"
  client_script     = "// client-side script"
  processing_script = "// server-side processing script"
}
