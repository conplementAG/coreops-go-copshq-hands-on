resource "azurerm_storage_account" "app-storage" {
  name                     = var.storage_account_name
  resource_group_name      = azurerm_resource_group.app-rg.name
  location                 = azurerm_resource_group.app-rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_table" "app-storage-table" {
  name                 = "Notes"
  storage_account_name = azurerm_storage_account.app-storage.name
}

resource "azurerm_key_vault_secret" "app-storage-connection-string" {
  name         = "connection-string"
  value        = "DefaultEndpointsProtocol=https;AccountName=${azurerm_storage_account.app-storage.name};AccountKey=${azurerm_storage_account.app-storage.primary_access_key};BlobEndpoint=https://${azurerm_storage_account.app-storage.name}.blob.core.windows.net/;QueueEndpoint=https://${azurerm_storage_account.app-storage.name}.queue.core.windows.net/;TableEndpoint=https://${azurerm_storage_account.app-storage.name}.table.core.windows.net/;FileEndpoint=https://${azurerm_storage_account.app-storage.name}.file.core.windows.net/;"
  key_vault_id = azurerm_key_vault.app-keyvault.id

  depends_on = [azurerm_role_assignment.keyvault-admin]
}

resource "azurerm_key_vault_secret" "app-storage-accountname" {
  name         = "accountname"
  value        = azurerm_storage_account.app-storage.name
  key_vault_id = azurerm_key_vault.app-keyvault.id

  depends_on = [azurerm_role_assignment.keyvault-admin]
}

resource "azurerm_key_vault_secret" "app-storage-accesskey" {
  name         = "accesskey"
  value        = azurerm_storage_account.app-storage.primary_access_key
  key_vault_id = azurerm_key_vault.app-keyvault.id

  depends_on = [azurerm_role_assignment.keyvault-admin]
}