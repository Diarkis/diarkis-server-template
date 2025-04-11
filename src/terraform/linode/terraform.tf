terraform {
  required_version = "~> 1.11.0"
  required_providers {
    linode = {
      source  = "linode/linode"
      version = "~> 2.35.0"
    }
  }
}

provider "linode" {
  token = var.token
}
