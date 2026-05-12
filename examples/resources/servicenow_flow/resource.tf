# Manages the high-level metadata of a Flow Designer flow.
#
# NOTE: The full Flow Designer schema (steps, variables, triggers) is not
# exposed by Terraform. Design the body of flows in the Flow Designer UI and
# reference them here, or use the `servicenow_flow` data source to look up an
# existing flow by name.

resource "servicenow_flow" "example" {
  name         = "example-flow"
  description  = "Example flow created by Terraform"
  active       = true
  trigger_type = "record_created"
  category     = "flow"
  status       = "draft"
}
