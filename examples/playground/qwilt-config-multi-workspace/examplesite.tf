resource "qwilt_cdn_site" "example" {
  site_name = terraform.workspace == "prod" ? "Terraform Advanced Workspace Example Site 1" : "Terraform Advanced Workspace Example Site 1 Testing"
}

resource "qwilt_cdn_site_configuration" "example" {
  site_id            = qwilt_cdn_site.example.site_id
  host_index         = templatefile("./examplesite.json", { hostname = terraform.workspace == "prod" ? "tfws.example.com" : "tfws-test.example.com" })
  change_description = "Advanced multi-config example demonstrating the Terraform plugin with workspaces"
}

resource "qwilt_cdn_certificate" "example" {
  certificate       = terraform.workspace == "prod" ? filebase64("./tfws.example.com.crt") : filebase64("./tfws-test.example.com.crt")
  certificate_chain = terraform.workspace == "prod" ? filebase64("./tfws.example.com.crt") : filebase64("./tfws-test.example.com.crt")
  private_key       = terraform.workspace == "prod" ? filebase64("./tfws.example.com.key") : filebase64("./tfws-test.example.com.key")
  description       = "Certificate for the Terraform advanced multi-config example with workspaces"
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
  value     = qwilt_cdn_certificate.example
  sensitive = true
}

output "examplesiteactivation" {
  value = qwilt_cdn_site_activation.example
}

output "examplesite-host_index" {
  value = qwilt_cdn_site_configuration.example.host_index
}
