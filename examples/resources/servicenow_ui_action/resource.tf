# Manages a UI action button on the incident form in ServiceNow.
resource "servicenow_ui_action" "example" {
  name         = "Example UI Action"
  table        = "incident"
  action_name  = "example_ui_action"
  active       = true
  form_button  = true
  show_insert  = true
  show_update  = true
  hint         = "Click to perform the example action."
  order        = 100
  script       = <<-EOT
    gs.addInfoMessage('Example UI action triggered for: ' + current.number);
    action.setRedirectURL(current);
  EOT
}
