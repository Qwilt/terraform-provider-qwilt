#keep an empty resource to import into
#after import is completed the user should manually set the required attributes in the resource from the imported state file
resource "qwilt_cdn_site" "example" {
}

terraform import qwilt_cdn_site.example <site_id>