# Look up an existing encryption context record in ServiceNow by name.
data "servicenow_encryption_context" "example" {
  name = "example-encryption-context"
}
