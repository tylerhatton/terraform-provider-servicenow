# Manages a header applied to all methods of a REST message.
resource "servicenow_rest_message" "example" {
  name          = "example-rest-message"
  rest_endpoint = "https://example.com/api"
}

resource "servicenow_rest_message_header" "example" {
  name            = "Content-Type"
  value           = "application/json"
  rest_message_id = servicenow_rest_message.example.id
}
