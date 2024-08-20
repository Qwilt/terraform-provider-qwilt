terraform {
  required_providers {
    qwilt = {
      source = "qwilt.com/qwiltinc/qwilt"
    }
  }
}

provider "qwilt" {
  xapi_token = var.token
}
