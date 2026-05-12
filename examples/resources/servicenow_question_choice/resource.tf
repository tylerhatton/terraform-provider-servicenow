# Manages a choice option for a catalog item variable (Select Box, Lookup, etc.).
resource "servicenow_service_catalog" "example" {
  title  = "Example Catalog"
  active = true
}

resource "servicenow_service_catalog_item" "example" {
  name              = "example-catalog-item"
  short_description = "Example catalog item"
  service_catalogs  = servicenow_service_catalog.example.id
  active            = true
}

resource "servicenow_service_catalog_variable" "example" {
  name         = "example_choice_var"
  question     = "Pick an option"
  catalog_item = servicenow_service_catalog_item.example.id
  type         = "Select Box"
  order        = "100"
}

resource "servicenow_question_choice" "example" {
  text     = "Option 1"
  value    = "option_1"
  question = servicenow_service_catalog_variable.example.id
  order    = "100"
}
