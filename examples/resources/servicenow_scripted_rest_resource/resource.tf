# Manages a resource (endpoint) under a scripted REST API.
resource "servicenow_scripted_rest_api" "example" {
  name       = "Example API"
  service_id = "example_api"
  active     = true
}

resource "servicenow_scripted_rest_resource" "example" {
  name                   = "example-resource"
  http_method            = "GET"
  relative_path          = "/items"
  web_service_definition = servicenow_scripted_rest_api.example.id
  operation_script       = "(function process(/*RESTAPIRequest*/ request, /*RESTAPIResponse*/ response) {\n  response.setBody({status: 'ok'});\n})(request, response);"
  active                 = true
}
