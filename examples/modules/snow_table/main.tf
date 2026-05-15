# snow_table — convenience wrapper around the three ServiceNow resources that
# together define a custom table:
#
#   1. servicenow_db_table  — the table itself (sys_db_object row)
#   2. servicenow_dictionary — one per column on the table
#   3. servicenow_choice    — choice-list values for choice-typed columns
#
# Why a module rather than a single provider resource? Each layer is a
# distinct ServiceNow JSONv2 endpoint with its own lifecycle, and other
# resources (e.g. servicenow_record) reference columns individually by
# element name. Keeping each layer as a Terraform resource preserves
# fine-grained diff handling and rollback while the module gives users a
# concise way to declare a table with its columns in one block.

resource "servicenow_db_table" "this" {
  label       = var.label
  user_role   = var.user_role
  extendable  = var.extendable
  super_class = var.super_class
}

resource "servicenow_dictionary" "columns" {
  for_each = var.columns

  name          = servicenow_db_table.this.name
  element       = each.key
  column_label  = each.value.column_label
  internal_type = each.value.internal_type
  max_length    = each.value.max_length
  mandatory     = each.value.mandatory
  read_only     = each.value.read_only
  active        = each.value.active
  display       = each.value.display
  unique        = each.value.unique
  default_value = each.value.default_value
  comments      = each.value.comments
  reference     = each.value.reference
  choice        = each.value.choice
}

# Flatten {column → {value → label}} into a single keyed map so for_each can
# walk all choices across all columns in one resource block.
locals {
  flat_choices = flatten([
    for col_name, col in var.columns : [
      for value, label in col.choices : {
        column = col_name
        value  = value
        label  = label
      }
    ]
  ])
  choices_map = {
    for c in local.flat_choices : "${c.column}__${c.value}" => c
  }
}

resource "servicenow_choice" "choices" {
  for_each = local.choices_map

  name    = servicenow_db_table.this.name
  element = servicenow_dictionary.columns[each.value.column].element
  value   = each.value.value
  label   = each.value.label
}
