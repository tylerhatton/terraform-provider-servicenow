# Manages a REST method (HTTP verb) under a REST message.
resource "servicenow_rest_message" "example" {
  name          = "example-rest-message"
  rest_endpoint = "https://example.com/api"
}

resource "servicenow_rest_method" "example" {
  name            = "get_items"
  rest_message_id = servicenow_rest_message.example.id
  http_method     = "get"
}
