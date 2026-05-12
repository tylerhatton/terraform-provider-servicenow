# Associates a CSS include with a widget dependency for load ordering.
resource "servicenow_css_include" "example" {
  name   = "example-css-include"
  source = "url"
  url    = "https://example.com/style.css"
}

resource "servicenow_widget_dependency" "example" {
  name = "example-widget-dependency"
}

resource "servicenow_css_include_relation" "example" {
  css_include_id = servicenow_css_include.example.id
  dependency_id  = servicenow_widget_dependency.example.id
  order          = 100
}
