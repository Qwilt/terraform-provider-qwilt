terraform {
  required_providers {
    qwiltcdn = {
      source  = "qwilt.com/qwiltinc/qwilt"
      version = "1.0.0"
    }
  }
}

provider "qwiltcdn" {
  xapi_token = var.token
}
