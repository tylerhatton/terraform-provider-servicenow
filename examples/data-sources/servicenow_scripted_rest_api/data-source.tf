# Look up an existing scripted REST API in ServiceNow by name.
data "servicenow_scripted_rest_api" "example" {
  name = "MyExistingScriptedRestApi"
}
