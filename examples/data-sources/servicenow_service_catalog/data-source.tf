# Look up an existing service catalog in ServiceNow by title.
data "servicenow_service_catalog" "example" {
  title = "Service Catalog"
}
