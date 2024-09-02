resource "qwilt_cdn_site" "example" {
  site_name = "Terraform Advanced Example Site 1"
}

resource "qwilt_cdn_site_configuration" "example" {
  site_id            = qwilt_cdn_site.example.site_id
  host_index         = file("./examplesite.json")
  change_description = "Advanced multi-config example demonstrating the Terraform plugin"
}

resource "qwilt_cdn_certificate" "example" {
  certificate       = filebase64("./tfm.example.com.crt")
  certificate_chain = filebase64("./tfm.example.com.crt")
  private_key       = filebase64("./tfm.example.com.key")
  description       = "Certificate for the Terraform advanced multi-config example"
}

resource "qwilt_cdn_site_activation" "example" {
  site_id        = qwilt_cdn_site_configuration.example.site_id
  revision_id    = qwilt_cdn_site_configuration.example.revision_id
  certificate_id = qwilt_cdn_certificate.example.cert_id
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

output "examplesite-host_index" {
  value = qwilt_cdn_site_configuration.example.host_index
}
