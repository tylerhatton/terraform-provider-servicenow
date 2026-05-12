# Manages a connection or credential alias used by Flow Designer actions.
resource "servicenow_alias" "example" {
  name = "example-alias"
  type = "credential"
}
