# Look up an existing JDBC connection in ServiceNow by name.
data "servicenow_jdbc_connection" "example" {
  name = "example-jdbc-connection"
}
