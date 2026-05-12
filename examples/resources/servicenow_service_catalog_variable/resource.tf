# Manages a variable (input field) on a service catalog item.
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
  name         = "example_variable"
  question     = "Please enter a value"
  catalog_item = servicenow_service_catalog_item.example.id
  type         = "Single Line Text"
  order        = "100"
  mandatory    = false
  active       = true
}
