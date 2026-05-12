# Manages a PEM trust store certificate in ServiceNow.
resource "servicenow_certificate" "example" {
  name              = "example-trust-cert"
  short_description = "Trust certificate uploaded by Terraform"
  format            = "pem"
  type              = "trust_store"
  active            = true

  # Sensitive fields are write-only - ServiceNow does not return them on read.
  pem_certificate = <<-EOT
    -----BEGIN CERTIFICATE-----
    MIIDazCCAlOgAwIBAgIUExampleCertificateBodyHere==
    -----END CERTIFICATE-----
  EOT
}
