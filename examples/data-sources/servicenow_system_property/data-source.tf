# Look up an existing system property in ServiceNow by its suffix.
data "servicenow_system_property" "example" {
  suffix = "glide.servlet.uri"
}
