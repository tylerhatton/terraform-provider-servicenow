# Look up an existing server CMDB entry in ServiceNow by name.
data "servicenow_server" "example" {
  name = "MyExistingServer"
}
