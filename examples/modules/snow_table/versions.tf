terraform {
  required_version = ">= 1.3"

  required_providers {
    servicenow = {
      source  = "tylerhatton/servicenow"
      version = ">= 0.11.0"
    }
  }
}
