# Manages a basic auth credential attached to a credential alias.
resource "servicenow_alias" "example" {
  name = "example-credential-alias"
  type = "credential"
}

resource "servicenow_basic_auth_credential" "example" {
  name             = "example-basic-auth"
  username         = "service-account"
  password         = "change-me"
  credential_alias = servicenow_alias.example.id
}
