# Look up an existing MID server in ServiceNow by name.
data "servicenow_mid_server" "example" {
  name = "example-mid-server"
}
