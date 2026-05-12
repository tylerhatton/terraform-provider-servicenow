# Associates a JS include with a widget dependency for load ordering.
resource "servicenow_js_include" "example" {
  display_name = "example-js-include"
  source       = "url"
  url          = "https://example.com/script.js"
}

resource "servicenow_widget_dependency" "example" {
  name = "example-widget-dependency"
}

resource "servicenow_js_include_relation" "example" {
  js_include_id = servicenow_js_include.example.id
  dependency_id = servicenow_widget_dependency.example.id
  order         = 100
}
