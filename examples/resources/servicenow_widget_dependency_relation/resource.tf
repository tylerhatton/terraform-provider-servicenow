# Associates a widget with a widget dependency (many-to-many).
resource "servicenow_widget" "example" {
  identifier = "example-widget"
  name       = "Example Widget"
  template   = "<div></div>"
}

resource "servicenow_widget_dependency" "example" {
  name = "example-widget-dependency"
}

resource "servicenow_widget_dependency_relation" "example" {
  widget_id     = servicenow_widget.example.id
  dependency_id = servicenow_widget_dependency.example.id
}
