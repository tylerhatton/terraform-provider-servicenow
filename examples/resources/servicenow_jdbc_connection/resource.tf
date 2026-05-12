# Manages a JDBC connection tied to a connection alias and (optionally) a MID server.
resource "servicenow_alias" "example" {
  name            = "example-jdbc-connection-alias"
  type            = "connection"
  connection_type = "jdbc_connection"
}

resource "servicenow_jdbc_connection" "example" {
  name             = "example-jdbc-connection"
  connection_alias = servicenow_alias.example.id
  connection_url   = "jdbc:mysql://db.example.com:3306/app"
  database_name    = "app"
  database_type    = "mysql"
  active           = true
  use_mid_server   = false
  mid_selection    = "auto_select"
}
