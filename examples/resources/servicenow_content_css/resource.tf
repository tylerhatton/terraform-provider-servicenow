# Manages a Content Management style sheet (CSS) in ServiceNow.
resource "servicenow_content_css" "example" {
  name  = "example-content-css"
  type  = "local"
  style = "body { color: #333; }"
}
