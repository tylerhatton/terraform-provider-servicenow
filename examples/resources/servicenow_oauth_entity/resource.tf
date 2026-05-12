# Manages an OAuth application entity in ServiceNow.
resource "servicenow_oauth_entity" "example" {
  name                   = "example-oauth-entity"
  redirect_url           = "https://example.com/callback"
  access_token_lifespan  = 1800
  refresh_token_lifespan = 8640000
}
