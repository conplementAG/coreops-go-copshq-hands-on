resource "azurerm_key_vault" "app-keyvault" {
  name                      = var.keyvault_name
  location                  = azurerm_resource_group.app-rg.location
  resource_group_name       = azurerm_resource_group.app-rg.name
  tenant_id                 = var.tenant_id
  sku_name                  = "standard"
  enable_rbac_authorization = true

  # The number of days that items should be retained for once soft-deleted. Default 90 cannot be updated.
  soft_delete_retention_days = 7
  # Once enabled it's not possible to disable it. Deletion will be done by Azure as scheduled: 90 day.
  purge_protection_enabled   = false
}

resource "azurerm_role_assignment" "keyvault-admin" {
  scope                = azurerm_key_vault.app-keyvault.id
  role_definition_name = "Key Vault Administrator"
  principal_id         = var.principal_id
}