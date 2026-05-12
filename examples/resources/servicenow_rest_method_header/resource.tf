# Manages a header applied to a single REST method.
resource "servicenow_rest_message" "example" {
  name          = "example-rest-message"
  rest_endpoint = "https://example.com/api"
}

resource "servicenow_rest_method" "example" {
  name            = "get_items"
  rest_message_id = servicenow_rest_message.example.id
  http_method     = "get"
}

resource "servicenow_rest_method_header" "example" {
  name           = "Accept"
  value          = "application/json"
  rest_method_id = servicenow_rest_method.example.id
}
