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

resource "qwiltcdn_site" "example" {
  site_name      = ""
}

resource "qwiltcdn_site_configuration" "example" {
  site_id = ""
  #host_index = file("./examplesitebasic.json")
  host_index         = ""
  change_description = ""
}

resource "qwiltcdn_certificate" "example" {
  certificate       = ""
  certificate_chain = ""
  private_key       = ""
  description       = ""
}

resource "qwiltcdn_site_activation" "example" {
  site_id        = ""
  revision_id    = ""
  certificate_id = ""
}

output "examplesite" {
  value = qwiltcdn_site.example
}

output "examplesiteconfig" {
  value = qwiltcdn_site_configuration.example
}

output "examplecertificate" {
  value = qwiltcdn_certificate.example
}

output "examplesiteactivation" {
  value = qwiltcdn_site_activation.example
}
