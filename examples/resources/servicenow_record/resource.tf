# Manages an arbitrary row in any ServiceNow table.
#
# The `fields` map declares the columns Terraform owns and will drift-check.
# All other columns ServiceNow populates (number, sys_created_on, state, ...)
# are returned in the read-only `output` map.

# Look up the assignment group's sys_id first so the example is portable.
data "servicenow_record" "network_team" {
  table = "sys_user_group"
  query = "name=Network"
}

resource "servicenow_record" "outage" {
  table = "incident"

  fields = {
    short_description = "Edge router rebooted unexpectedly"
    urgency           = "1"
    impact            = "1"
    assignment_group  = data.servicenow_record.network_team.id
  }
}

output "incident_number" {
  value = servicenow_record.outage.output["number"]
}
