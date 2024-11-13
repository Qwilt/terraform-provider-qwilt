# The "required_providers" block declares the provider:

terraform {
  required_providers {
    qwilt = {
      source = "Qwilt/qwilt"
    }
  }
}


# The "provider" config sets authentication.
# Read about authentication, above.

provider "qwilt" {
}


# A named site resource:

resource "qwilt_cdn_site" "example" {
  site_name = "Terraform Basic Example Site"
}

# A named site configuration resource:

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
  change_description = "Basic example demonstrating the Terraform plugin."
}


# A named certificate resource:

resource "qwilt_cdn_certificate" "example" {
  certificate       = filebase64("./tf.example.com.crt")
  certificate_chain = filebase64("./tf.example.com.crt")
  private_key       = filebase64("./tf.example.com.key")
  description       = "Certificate for the Terraform basic example configuration"
}


# A named site activation resource:

resource "qwilt_cdn_site_activation" "example" {
  site_id        = qwilt_cdn_site_configuration.example.site_id
  revision_id    = qwilt_cdn_site_configuration.example.revision_id
  certificate_id = qwilt_cdn_certificate.example.cert_id
}


