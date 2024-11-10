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
  common_name = "yuval.weisz.com"
  auto_managed_certificate_template = true
  organization_name = "Qwilt"
  locality = "Hod Hasharon"
  country = "IL"
  state = "Israel"
  sans = ["www.yuval.weisz.com"]
}