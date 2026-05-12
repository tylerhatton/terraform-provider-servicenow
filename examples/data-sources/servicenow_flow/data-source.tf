# Look up an existing Flow Designer flow in ServiceNow by name.
data "servicenow_flow" "example" {
  name = "example-flow"
}
