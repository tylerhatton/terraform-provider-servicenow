# Manages a Service Portal Widget in ServiceNow.
resource "servicenow_widget" "example" {
  identifier  = "example-widget"
  name        = "Example Widget"
  template    = "<div>Example Widget Body</div>"
  description = "An example Service Portal widget."
}
