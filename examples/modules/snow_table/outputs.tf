output "id" {
  description = "sys_id of the underlying sys_db_object row."
  value       = servicenow_db_table.this.id
}

output "name" {
  description = "Auto-generated internal table name (e.g. `u_demo_table`). Use this as `table = module.x.name` on downstream `servicenow_record` resources."
  value       = servicenow_db_table.this.name
}

output "columns" {
  description = "Map of element name → sys_id of the underlying sys_dictionary entry."
  value       = { for k, c in servicenow_dictionary.columns : k => c.id }
}

output "choices" {
  description = "Map of `<column>__<value>` → sys_id of the underlying sys_choice entry."
  value       = { for k, c in servicenow_choice.choices : k => c.id }
}
