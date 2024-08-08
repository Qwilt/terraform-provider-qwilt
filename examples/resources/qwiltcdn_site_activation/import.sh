#keep an empty resource to import into
#after import is completed the user manually sets the required attributes in the resource from the imported state file
#it is advised to change site_id and site_id and revision_id attributes with references to qwiltcdn_site_configuration resource to achieve implicit dependency
resource "qwiltcdn_site_activation" "example" {
}

# qwiltcdn_site_activation can be imported using their site ID, e.g.
terraform import qwiltcdn_site_activation.example <site_id>

# By default, either:
# 1. the active publish operation will be imported
# 2. the latest published operation if no publish is active
# Alternatively, a specific publish_id of the site configuration can be selected by appending an : followed by the publish_id to the site ID, e.g. */terraform import qwiltcdn_site_activation.example xxxxxxxxxxxxxxxxxxxx:yyyyyyyyyyyyyyyyyyy
terraform import qwiltcdn_site_activation.example <site_id>:<publish_id>