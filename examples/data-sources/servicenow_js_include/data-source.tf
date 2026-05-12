# Look up an existing JS include in ServiceNow by display name.
data "servicenow_js_include" "example" {
  display_name = "my-js-include"
}
