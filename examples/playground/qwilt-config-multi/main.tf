terraform {
  required_providers {
    qwilt = {
      source = "Qwilt/qwilt"
      version = "0.1.4"
    }
  }
}

provider "qwilt" {
  xapi_token = var.token
}
