# Manages a server-side Script Include in ServiceNow.
resource "servicenow_script_include" "example" {
  name   = "ExampleScriptInclude"
  active = true
  script = <<-EOT
    var ExampleScriptInclude = Class.create();
    ExampleScriptInclude.prototype = {
      initialize: function() {},
      type: 'ExampleScriptInclude'
    };
  EOT
}
