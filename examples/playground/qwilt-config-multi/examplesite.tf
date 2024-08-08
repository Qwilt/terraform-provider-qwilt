resource "qwiltcdn_site" "example" {
  site_name      = "Terraform Advanced Example Site 1"
  routing_method = "HTTP"
}

resource "qwiltcdn_site_configuration" "example" {
  site_id            = qwiltcdn_site.example.site_id
  host_index         = file("./examplesite.json")
  change_description = "Advanced multi-config example demonstrating the Terraform plugin"
}

resource "qwiltcdn_certificate" "example" {
  certificate       = filebase64("./tfm.example.com.crt")
  certificate_chain = filebase64("./tfm.example.com.crt")
  private_key       = filebase64("./tfm.example.com.key")
  description       = "Certificate for the Terraform advanced multi-config example"
}

resource "qwiltcdn_site_activation" "example" {
  site_id        = qwiltcdn_site_configuration.example.site_id
  revision_id    = qwiltcdn_site_configuration.example.revision_id
  certificate_id = qwiltcdn_certificate.example.cert_id
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

output "examplesite-host_index" {
  value = qwiltcdn_site_configuration.example.host_index
}
