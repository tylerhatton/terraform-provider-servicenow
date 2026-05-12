# Manages a UI Macro (Jelly template) in ServiceNow.
resource "servicenow_ui_macro" "example" {
  name = "example_macro"
  xml  = "<j:jelly xmlns:j='jelly:core' xmlns:g='glide' xmlns:j2='null' xmlns:g2='null'><h1>Example Macro</h1></j:jelly>"
}
