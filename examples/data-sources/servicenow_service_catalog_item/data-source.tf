# Look up an existing service catalog item in ServiceNow by name.
data "servicenow_service_catalog_item" "example" {
  name = "MyExistingServiceCatalogItem"
}
