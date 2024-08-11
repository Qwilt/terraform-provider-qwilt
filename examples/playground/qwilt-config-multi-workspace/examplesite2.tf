resource "qwiltcdn_site" "example2" {
  site_name = terraform.workspace == "prod" ? "Terraform Advanced Workspace Example Site 2" : "Terraform Advanced Workspace Example Site 2 Testing"
}

resource "qwiltcdn_site_configuration" "example2" {
  site_id            = qwiltcdn_site.example2.site_id
  host_index         = templatefile("./examplesite2.json", { hostname = terraform.workspace == "prod" ? "tfws2.example.com" : "tfws2-test.example.com" })
  change_description = "Advanced multi-config example 2 demonstrating the Terraform plugin with workspaces"
}

resource "qwiltcdn_certificate" "example2" {
  certificate       = terraform.workspace == "prod" ? filebase64("./tfws2.example.com.crt") : filebase64("./tfws2-test.example.com.crt")
  certificate_chain = terraform.workspace == "prod" ? filebase64("./tfws2.example.com.crt") : filebase64("./tfws2-test.example.com.crt")
  private_key       = terraform.workspace == "prod" ? filebase64("./tfws2.example.com.key") : filebase64("./tfws2-test.example.com.key")
  description       = "Certificate for the Terraform advanced multi-config example 2 with workspaces"
}

resource "qwiltcdn_site_activation" "example2" {
  site_id        = qwiltcdn_site_configuration.example2.site_id
  revision_id    = qwiltcdn_site_configuration.example2.revision_id
  certificate_id = qwiltcdn_certificate.example2.cert_id
}

output "examplesite2" {
  value = qwiltcdn_site.example2
}

output "examplesiteconfig2" {
  value = qwiltcdn_site_configuration.example2
}

output "examplecertificate2" {
  value = qwiltcdn_certificate.example2
}

output "examplesiteactivation2" {
  value = qwiltcdn_site_activation.example2
}

output "examplesite2-host_index" {
  value = qwiltcdn_site_configuration.example2.host_index
}
