resource "qwilt_cdn_site" "example2" {
  site_name = "Terraform Advanced Example Site 2"
}

resource "qwilt_cdn_site_configuration" "example2" {
  site_id            = qwilt_cdn_site.example2.site_id
  host_index         = file("./examplesite2.json")
  change_description = "Advanced multi-config example 2 demonstrating the Terraform plugin"
}

resource "qwilt_cdn_certificate" "example2" {
  certificate       = filebase64("./tfm2.example.com.crt")
  certificate_chain = filebase64("./tfm2.example.com.crt")
  private_key       = filebase64("./tfm2.example.com.key")
  description       = "Certificate for the Terraform advanced multi-config example 2"
}

resource "qwilt_cdn_site_activation" "example2" {
  site_id        = qwilt_cdn_site_configuration.example2.site_id
  revision_id    = qwilt_cdn_site_configuration.example2.revision_id
  certificate_id = qwilt_cdn_certificate.example2.cert_id
}

output "examplesite2" {
  value = qwilt_cdn_site.example2
}

output "examplesiteconfig2" {
  value = qwilt_cdn_site_configuration.example2
}

output "examplecertificate2" {
  value     = qwilt_cdn_certificate.example2
  sensitive = true
}

output "examplesiteactivation2" {
  value = qwilt_cdn_site_activation.example2
}

output "examplesite2-host_index" {
  value = qwilt_cdn_site_configuration.example2.host_index
}
