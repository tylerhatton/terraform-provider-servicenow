# Manages an email notification (sysevent_email_action) in ServiceNow.
resource "servicenow_notification" "example" {
  name        = "Example Incident Notification"
  table       = "incident"
  event_name  = "incident.inserted"
  active      = true
  description = "Sends an email when a new incident is inserted."
  subject     = "New incident ${number} has been created"
  message_html = <<-EOT
    <p>A new incident has been created.</p>
    <p>Number: $${number}</p>
    <p>Short description: $${short_description}</p>
  EOT
  weight = 100
  type   = "EMAIL"
}
