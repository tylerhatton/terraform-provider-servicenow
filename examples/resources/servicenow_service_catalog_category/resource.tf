# Manages a category that groups items within a service catalog.
resource "servicenow_service_catalog" "example" {
  title  = "Example Catalog"
  active = true
}

resource "servicenow_service_catalog_category" "example" {
  title   = "Example Category"
  catalog = servicenow_service_catalog.example.id
  active  = true
}
