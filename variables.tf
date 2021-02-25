variable "stage" {
    type = string
    default = "LANDINGZONE"
  
}

variable "type" {
    type = string
    default = "module"
}

variable account_tier {
    type = string
}

variable account_replication_type {
    type = string
    default = "lrs"
}

locals {
    tags = {
        provisionedBy = "https://github.com/czelabueno/infrastructure-as-code-testing"
    }
}

data "azurerm_subscription" "current" {}