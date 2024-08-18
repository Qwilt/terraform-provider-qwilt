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
  site_name = "Terraform Basic Example Site 88"
}

resource "qwiltcdn_site_configuration" "example" {
  site_id = qwiltcdn_site.example.site_id
  #host_index = file("./examplesitebasic.json")
  host_index         = <<-EOT
{
	"hosts": [
		{
			"host": "tf1.example.com",
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

resource "qwiltcdn_certificate" "example" {
  certificate       = filebase64("./tf.example.com.crt")
  certificate_chain = filebase64("./tf.example.com.crt")
  private_key       = filebase64("./tf.example.com.key")
  description       = "Certificate for the Terraform basic example configuration"
}

resource "qwiltcdn_site_activation" "example" {
  site_id        = qwiltcdn_site_configuration.example.site_id
  revision_id    = qwiltcdn_site_configuration.example.revision_id
  certificate_id = null
}

#output "examplesite" {
#  value = qwiltcdn_site.example
#}
#
#output "examplesiteconfig" {
#  value = qwiltcdn_site_configuration.example
#}
#
#output "examplecertificate" {
#  value = qwiltcdn_certificate.example
#}
#
#output "examplesiteactivation" {
#  value = qwiltcdn_site_activation.example
#}
