
#keep an empty resource to import into
#after import is completed the user should manually set the required attributes in the resource from the imported state file
#it is advised to change site_id attribute with references to qwiltcdn_site resource to achieve implicit dependency
resource "qwiltcdn_site_configuration" "example" {
}

# site_configurations can be imported using their site ID, e.g.
terraform import qwiltcdn_site_configuration.example <site_id>

# By default, either:
# 1. the active revision will be imported
# 2. the latest published revision if no version is active
# 3. the latest configured revision, if no published revision
# Alternatively, a specific revision_id of the site configuration can be selected by appending an : followed by the revision_id to the site ID, e.g. */

terraform import qwiltcdn_site_configuration.example <site_id>:<revision_id>