terraform {
  required_providers {
    qwiltcdn = {
      source = "qwilt.com/qwiltinc/qwilt"
    }
  }
}

provider "qwiltcdn" {
  xapi_token = var.token
}

data "qwiltcdn_sites" "detail" {
  filter = {}
}
