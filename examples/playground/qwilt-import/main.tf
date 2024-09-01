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

resource "qwilt_cdn_site" "example" {
  site_name = ""
}

resource "qwilt_cdn_site_configuration" "example" {
  site_id = ""
  #host_index = file("./examplesitebasic.json")
  host_index         = ""
  change_description = ""
}

resource "qwilt_cdn_certificate" "example" {
  certificate       = ""
  certificate_chain = ""
  private_key       = ""
  description       = ""
}

resource "qwilt_cdn_site_activation" "example" {
  site_id        = ""
  revision_id    = ""
  certificate_id = ""
}

output "examplesite" {
  value = qwilt_cdn_site.example
}

output "examplesiteconfig" {
  value = qwilt_cdn_site_configuration.example
}

output "examplecertificate" {
  value = qwilt_cdn_certificate.example
}

output "examplesiteactivation" {
  value = qwilt_cdn_site_activation.example
}
