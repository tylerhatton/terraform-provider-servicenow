# Look up an existing certificate record in ServiceNow by name.
data "servicenow_certificate" "example" {
  name = "example-trust-cert"
}
