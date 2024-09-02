terraform {
  required_providers {
    qwilt = {
      source = "Qwilt/qwilt"
    }
  }
}

provider "qwilt" {
}

data "qwilt_cdn_sites" "detail" {
  filter = {}
}
