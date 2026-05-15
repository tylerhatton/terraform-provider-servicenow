# Look up a single row in any ServiceNow table.
# Supply either `sys_id` for a direct lookup or `query` for an encoded
# sysparm_query expression that must match exactly one record.

data "servicenow_record" "admin_group" {
  table = "sys_user_group"
  query = "name=Admin"
}

output "admin_group_sys_id" {
  value = data.servicenow_record.admin_group.id
}

output "admin_group_email" {
  value = data.servicenow_record.admin_group.output["email"]
}
