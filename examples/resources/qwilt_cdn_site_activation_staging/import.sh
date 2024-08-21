> ⚠️ **Disclaimer**: this resource is experimental and should not be used!

#keep an empty resource to import into
#after import is completed the user should manually set the required attributes in the resource from the imported state file
#it is advised to change site_id and site_id and revision_id attributes with references to qwilt_cdn_site_configuration resource to achieve implicit dependency
resource "qwilt_cdn_site_activation" "example" {
}

# site_activation_staging can be imported using their site ID, e.g.
terraform import qwilt_cdn_site_activation_staging.example <site_id>

# By default, either:
# 1. the active publish operation will be imported
# 2. the latest published operation if no publish is active
# Alternatively, a specific publish_id of the site configuration can be selected by appending an : followed by the publish_id to the site ID, e.g. */
terraform import qwilt_cdn_site_activation_staging.example <site_id>:<publish_id>