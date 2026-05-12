# Manages a MID server (ECC agent) record in ServiceNow.
#
# WARNING: MID server agents are normally installed on a host using the official
# ServiceNow MID server installer, which auto-registers the record with the
# instance. Creating a record purely through Terraform produces a placeholder row
# with no live agent backing it. Prefer the `servicenow_mid_server` data source
# to reference a MID server that the installer has already registered.

resource "servicenow_mid_server" "example" {
  name        = "example-mid-server"
  host_name   = "mid-host-01.example.com"
  description = "Placeholder MID server created via Terraform"
  agent_type  = "MIDServer"
  validated   = false
}
