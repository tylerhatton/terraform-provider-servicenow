# Manages an email template (sysevent_email_template) in ServiceNow.
resource "servicenow_email_template" "example" {
  name         = "Example Incident Template"
  table        = "incident"
  subject      = "Incident ${number} update"
  message_html = "<p>Incident has been updated.</p>"
  description  = "Reusable email template for incident updates."
  active       = true
}
