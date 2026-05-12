# Manages a service catalog item that users can request.
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
