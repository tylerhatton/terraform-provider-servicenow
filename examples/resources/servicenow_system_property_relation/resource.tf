# Associates a system property with a system property category.
resource "servicenow_system_property" "example" {
  suffix = "example.property"
  type   = "string"
}

resource "servicenow_system_property_category" "example" {
  name = "Example Property Category"
}

resource "servicenow_system_property_relation" "example" {
  property_id = servicenow_system_property.example.id
  category_id = servicenow_system_property_category.example.id
}
