# Look up an existing email template in ServiceNow by name.
data "servicenow_email_template" "example" {
  name = "Example Incident Template"
}
