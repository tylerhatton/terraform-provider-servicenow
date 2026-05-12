# Manages an encryption context record in ServiceNow Edge Encryption.
#
# NOTE: Requires the Edge Encryption plugin to be installed and activated on
# the target ServiceNow instance. On instances without the plugin the
# sys_encryption_context table does not exist and this resource cannot be used.

resource "servicenow_encryption_context" "example" {
  name        = "example-encryption-context"
  type        = "standard"
  description = "Example encryption context managed by Terraform"
  algorithm   = "AES_256"
  active      = true

  # Sensitive: not read back from ServiceNow once set.
  encryption_key = "base64-encoded-key-material"
}
