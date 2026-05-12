# Look up an existing data lookup definition in ServiceNow by name.
data "servicenow_data_lookup" "example" {
  name = "Incident Priority Lookup"
}
