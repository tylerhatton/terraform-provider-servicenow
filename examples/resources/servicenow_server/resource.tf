# Manages a server CMDB entry in ServiceNow.
resource "servicenow_server" "example" {
  name       = "example-server"
  ip_address = "192.168.1.100"
}
