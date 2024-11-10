terraform {
  required_providers {
    qwilt = {
      source = "qwilt.com/qwiltinc/qwilt"
    }
  }
}

provider "qwilt" {
}

resource "qwilt_cdn_certificate_template" "exmaple" {
  common_name = "example.com"
  auto_managed_certificate_template = true
}