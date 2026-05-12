# Look up an existing scripted extension point in ServiceNow by name.
data "servicenow_extension_point" "example" {
  name = "my-extension-point"
}
