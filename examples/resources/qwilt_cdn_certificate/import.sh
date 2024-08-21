
#keep an empty resource to import into (Test Esther Delete This)
#after import is completed the user should manually set the required attributes in the resource from the imported state file
resource "qwilt_cdn_certificate" "example" {

  #after import, the private_key will remain empty in the state.
  #to avoid plan keep detecting this as a change, add this lifecycle section to the imported resource
  lifecycle {
      ignore_changes = [
        private_key
      ]
    }
}

terraform import qwilt_cdn_certificate.example <cert_id>