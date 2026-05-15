variable "label" {
  description = "Display label for the table (ServiceNow auto-generates the internal `name` from this — e.g. \"Demo Table\" becomes `u_demo_table` in the global scope)."
  type        = string
}

variable "user_role" {
  description = "Role sys_id required for end-user access to the table. Empty string means anyone."
  type        = string
  default     = ""
}

variable "extendable" {
  description = "Allow other tables to extend this one."
  type        = bool
  default     = false
}

variable "super_class" {
  description = "sys_id of a parent table this one extends. Leave empty for a root table."
  type        = string
  default     = ""
}

# columns maps a column's element name (e.g. "u_environment") to the column's
# configuration. Use the `choices` sub-attribute to ship choice list values
# inline; for non-choice columns leave it empty.
variable "columns" {
  description = "Map of element name → column definition. Each entry expands to a `servicenow_dictionary` row, and any `choices` sub-map expands to `servicenow_choice` rows."
  type = map(object({
    column_label  = string
    internal_type = optional(string, "string")
    max_length    = optional(number, 0)
    mandatory     = optional(bool, false)
    read_only     = optional(bool, false)
    active        = optional(bool, true)
    display       = optional(bool, false)
    unique        = optional(bool, false)
    default_value = optional(string, "")
    comments      = optional(string, "")
    reference     = optional(string, "")
    # choice: 0 = none, 1 = suggestion, 3 = dropdown (ServiceNow convention).
    choice  = optional(number, 0)
    choices = optional(map(string), {})
  }))
  default = {}
}
