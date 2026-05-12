# Manages a CSS include (external URL or service portal style sheet).
resource "servicenow_css_include" "example" {
  name   = "example-css-include"
  source = "url"
  url    = "https://example.com/style.css"
}
