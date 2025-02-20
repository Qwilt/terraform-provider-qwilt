---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "qwilt_cdn_site_configuration Resource - qwilt"
subcategory: ""
description: |-
  Manages a Qwilt CDN site Configuration. The site configuration determines how the CDN processes client requests and delivers content.Learn how to prepare the configuration JSON. https://docs.qwilt.com/docs/terraform-user-guide#site-configuration-json
---

# qwilt_cdn_site_configuration (Resource)

Manages a Qwilt CDN site Configuration. The site configuration determines how the CDN processes client requests and delivers content.<br><br>[Learn how to prepare the configuration JSON.](https://docs.qwilt.com/docs/terraform-user-guide#site-configuration-json)

## Example Usage

```terraform
resource "qwilt_cdn_site_configuration" "example" {
  site_id            = qwilt_cdn_site.example.site_id
  host_index         = <<-EOT
{
	"hosts": [
		{
			"host": "www.basicdemo2.example.com",
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `change_description` (String) Comments added by the user to the configuration JSON payload.
- `host_index` (String) The SVTA metadata objects that define the delivery service configuration, in application/json format.
- `site_id` (String) The unique identifier of the Site.

### Read-Only

- `id` (String) For internal use only, for testing. Equals site_id:revision_id.
- `last_update_time_milli` (Number) When the site configuration was last updated, in epoch time.
- `owner_org_id` (String) The organization that owns the site.
- `revision_id` (String) The unique identifier of the configuration version.
- `revision_num` (Number) The unique revision number of the configuration version.

## Import

Import is supported using the following syntax:

```shell
#Create an empty resource to import into.

#After the import is complete, manually set the required attributes 
#in the resource based on the imported state.

#We recommend changing the site_id attribute to a reference 
#to the qwilt_cdn_site resource to achieve implicit dependency.


resource "qwilt_cdn_site_configuration" "example" {
}

  # You can import the qwilt_cdn_site_configuration resource by 
  # specifying just the site_id. 

      #For example: terraform import qwilt_cdn_site_configuration.example <site_id>
    
          # The process determines which saved site version (represented by the 
          # revisionId) to import based on the following conditions: 
          # - If there is an active published version, it is imported.
          # - If not, the most recently published version is imported. 
          # - If the site has never been published, the most recently saved 
          #   configuration version is imported.

  # Alternatively, you can specify a particular revision_id of the site 
  # configuration by adding a colon (:) and the revision_id.

terraform import qwilt_cdn_site_configuration.example <site_id>:<revision_id>
```
