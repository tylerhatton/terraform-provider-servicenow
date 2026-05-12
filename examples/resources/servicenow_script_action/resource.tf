# Manages an event-driven script action (sysevent_script_action) in ServiceNow.
resource "servicenow_script_action" "example" {
  name        = "Example Script Action"
  event_name  = "incident.commented"
  description = "Reacts to incident.commented events."
  active      = true
  order       = 100
  script      = <<-EOT
    gs.info("Event " + event.name + " fired with parm1=" + event.parm1);
  EOT
}
