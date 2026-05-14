# Manages a scheduled background script job (sysauto_script) in ServiceNow.
resource "servicenow_scheduled_job" "example" {
  name     = "Example Daily Job"
  run_type = "daily"
  run_time = "1970-01-01 04:00:00"
  active   = true
  script   = <<-EOT
    gs.info("Example scheduled job ran at " + new GlideDateTime());
  EOT
}
