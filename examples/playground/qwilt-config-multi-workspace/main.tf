terraform {
  required_providers {
    qwilt_cdn_ = {
      source  = "qwilt.com/qwiltinc/qwilt"
      version = "1.0.0"
    }
  }
}

provider "qwilt" {
  xapi_token = var.token
}
