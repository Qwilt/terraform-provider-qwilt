terraform {
  required_providers {
    qwilt = {
      source = "Qwilt/qwilt"
    }
  }
}

provider "qwilt" {
}

resource "qwilt_cdn_site" "example" {
  site_name = "Terraform Basic Example Site"
}

resource "qwilt_cdn_site_configuration" "example" {
  site_id = qwilt_cdn_site.example.site_id
  #host_index = file("./examplesitebasic.json")
  host_index         = <<-EOT
{
	"hosts": [
		{
			"host": "tf.example.com",
			"host-metadata": {
				"metadata": [
					{
						"generic-metadata-type": "MI.SourceMetadataExtended",
						"generic-metadata-value": {
							"sources": [
								{
									"protocol": "https/1.1",
									"endpoints": [
										"origin-host.example.com"
									]
								}
							]
						}
					}
				],
				"paths": []
			}
		}
	]
}
EOT
  change_description = "Basic example demonstrating the Terraform plugin"
}

resource "qwilt_cdn_certificate" "example" {
  certificate       = filebase64("./tf.example.com.crt")
  certificate_chain = filebase64("./tf.example.com.crt")
  private_key       = filebase64("./tf.example.com.key")
  description       = "Certificate for the Terraform basic example configuration"
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
