# resource "azurerm_user_assigned_identity" "workload_identity" {
#   name                = var.identity_name
#   location            = var.location
#   resource_group_name = azurerm_resource_group.app-rg.name
# }

# resource "azurerm_role_assignment" "keyvault-user" {
#   scope                = azurerm_key_vault.app-keyvault.id
#   role_definition_name = "Key Vault Secrets User"
#   principal_id         = azurerm_user_assigned_identity.workload_identity.principal_id
# }