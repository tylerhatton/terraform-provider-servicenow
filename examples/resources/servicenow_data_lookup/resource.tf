# Manages a data lookup definition (dl_definition) in ServiceNow.
resource "servicenow_data_lookup" "example" {
  name           = "Incident Priority Lookup"
  table          = "incident"
  lookup_table   = "dl_u_priority"
  matcher_fields = "impact,urgency"
  setter_fields  = "priority"
  active         = true
  description    = "Sets incident priority based on impact and urgency."
  order          = 100
}
