# Manages a REST message in ServiceNow.
resource "servicenow_rest_message" "example" {
  name          = "example-rest-message"
  rest_endpoint = "https://example.com/api"
  description   = "Example REST message"
}
