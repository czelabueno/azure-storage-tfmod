terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 2.46.0"
    }
  }
  backend "remote" {
    organization = "devopsperu-demo"

    workspaces {
      name = "az-storage-module"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "random_integer" "iacexample" {
  min = 1
  max = 10
}

resource "azurerm_resource_group" "iacexample" {
  name     = upper("${var.type}_${var.stage}")
  location = "eastus2"
  tags = local.tags
}

resource "azurerm_storage_account" "iacexample" {
  name                     = "iacstorage${random_integer.iacexample.result}"
  resource_group_name      = azurerm_resource_group.iacexample.name
  location                 = azurerm_resource_group.iacexample.location
  account_tier             = var.account_tier
  account_replication_type = var.account_replication_type
  account_kind              = "StorageV2"
  tags  = local.tags
  
  static_website {
    index_document = "index.html"
  }

}

resource "azurerm_storage_container" "iacexample" {
  name                  = "contentexample"
  storage_account_name  = azurerm_storage_account.iacexample.name
  container_access_type = "private"
}