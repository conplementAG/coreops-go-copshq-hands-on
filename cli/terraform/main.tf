terraform {
  backend "azurerm" {}
}

resource "azurerm_resource_group" "app-rg" {
  name     = var.resource_group_name
  location = var.location
}
