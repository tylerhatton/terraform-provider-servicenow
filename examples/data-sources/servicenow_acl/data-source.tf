# Look up an existing ACL in ServiceNow by name.
# Use operation and/or type to disambiguate when multiple ACLs share the same name.
data "servicenow_acl" "example" {
  name      = "sys_user"
  operation = "read"
}
