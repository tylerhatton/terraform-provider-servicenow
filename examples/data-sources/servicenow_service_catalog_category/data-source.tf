# Look up an existing service catalog category in ServiceNow by title.
data "servicenow_service_catalog_category" "example" {
  title = "Hardware"
}
