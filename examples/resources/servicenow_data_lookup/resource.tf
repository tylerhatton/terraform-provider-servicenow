# Manages a data lookup definition (dl_definition) in ServiceNow.
resource "servicenow_data_lookup" "example" {
  name               = "Incident Priority Lookup"
  table              = "incident"
  lookup_table       = "dl_matcher_definition"
  active             = true
  run_on_insert      = true
  run_on_update      = true
  run_on_form_change = false
}
