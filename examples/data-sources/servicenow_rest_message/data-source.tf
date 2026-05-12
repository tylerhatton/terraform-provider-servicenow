# Look up an existing REST message in ServiceNow by name.
data "servicenow_rest_message" "example" {
  name = "my-rest-message"
}
