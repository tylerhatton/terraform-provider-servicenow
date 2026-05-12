# Manages a navigator module (link) belonging to an application menu.
resource "servicenow_application_menu" "parent" {
  title = "Example Menu"
  order = 100
}

resource "servicenow_application_module" "example" {
  title               = "Example Module"
  application_menu_id = servicenow_application_menu.parent.id
  link_type           = "DIRECT"
  arguments           = "home.do"
  order               = 100
}
