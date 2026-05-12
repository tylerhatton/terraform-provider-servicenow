# Manages a JavaScript include (external URL or UI Script reference).
resource "servicenow_js_include" "example" {
  display_name = "example-js-include"
  source       = "url"
  url          = "https://example.com/script.js"
}
