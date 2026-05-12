# Manages a client-side script that runs in the browser in ServiceNow.
resource "servicenow_client_script" "example" {
  name        = "Example Client Script"
  table       = "incident"
  type        = "onChange"
  field_name  = "priority"
  active      = true
  description = "Shows a warning when the priority field is changed to critical."
  script      = <<-EOT
    function onChange(control, oldValue, newValue, isLoading, isTemplate) {
      if (isLoading || newValue === '') {
        return;
      }
      if (newValue === '1') {
        g_form.addInfoMessage('You have selected critical priority. Please confirm.');
      }
    }
  EOT
}
