# Look up an existing HTTP connection in ServiceNow by name.
data "servicenow_http_connection" "example" {
  name = "my-http-connection"
}
