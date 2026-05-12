# Manages a server-side business rule in ServiceNow.
resource "servicenow_business_rule" "example" {
  name           = "Example Business Rule"
  table          = "incident"
  when           = "before"
  order          = 100
  active         = true
  action_insert  = true
  action_update  = true
  description    = "Validates the incident short description on insert and update."
  script         = <<-EOT
    (function executeRule(current, previous /*null when async*/) {
      if (!current.short_description) {
        gs.addErrorMessage('Short description is required');
        current.setAbortAction(true);
      }
    })(current, previous);
  EOT
}
