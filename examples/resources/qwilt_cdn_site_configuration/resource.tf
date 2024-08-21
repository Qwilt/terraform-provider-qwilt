
resource "qwilt_cdn_site_configuration" "example" {
  site_id            = qwilt_cdn_site.example.site_id
  host_index         = <<-EOT
{
	"hosts": [
		{
			"host": "www.basicdemo2.kuku.com",
			"host-metadata": {
				"metadata": [
					{
						"generic-metadata-type": "MI.SourceMetadataExtended",
						"generic-metadata-value": {
							"sources": [
								{
									"protocol": "https/1.1",
									"endpoints": [
										"www.example-origin-host1.com"
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
  change_description = "Example demonstrating the Terraform plugin"
}