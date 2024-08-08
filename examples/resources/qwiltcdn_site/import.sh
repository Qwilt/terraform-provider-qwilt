#keep an empty resource to import into
#after import is completed the user should manually set the required attributes in the resource from the imported state file
resource "qwiltcdn_site" "example" {
}

terraform import qwiltcdn_site.example <site_id>