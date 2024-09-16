
#Create an empty resource to import into.
#After the import is complete, manually set the required attributes in the resource based on the imported state.

resource "qwilt_cdn_certificate" "example" {

  #After import, the private_key remains empty in the state.
  #To prevent this from being detected as a change every time you run ```terraform plan```, add this lifecycle section to the imported resource:

  lifecycle {
   ignore_changes = [
      private_key
    ]
  }

}

terraform import qwilt_cdn_certificate.example <cert_id>