terraform {
  required_providers {
    qwilt_cdn_ = {
      source = "qwilt.com/qwiltinc/qwilt"
    }
  }
}

provider "qwilt" {
  xapi_token = var.token
}

data "qwilt_cdn_sites" "detail" {
  filter = {}
}
