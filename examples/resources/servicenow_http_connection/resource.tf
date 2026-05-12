# Manages an HTTP connection configuration tied to a connection alias.
resource "servicenow_alias" "example" {
  name            = "example-http-connection-alias"
  type            = "connection"
  connection_type = "http_connection"
}

resource "servicenow_http_connection" "example" {
  name             = "example-http-connection"
  connection_alias = servicenow_alias.example.id
  connection_url   = "https://example.com"
  active           = true
}
