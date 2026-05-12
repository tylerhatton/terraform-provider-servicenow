# Manages a UI Script (client-side JavaScript) in ServiceNow.
resource "servicenow_ui_script" "example" {
  name   = "example_ui_script"
  script = "// example UI script\nconsole.log('hello');"
  active = true
}
