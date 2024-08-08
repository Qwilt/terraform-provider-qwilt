#
#Notes
#this resource take a long time to fully apply.
#any attempt to apply site_activation with the same site_id might encounter a failure due to another publish operation in-progress
#the user can run terraform refresh to sync the state of this resource explicitly


resource "qwiltcdn_site_activation" "example" {
  site_id     = qwiltcdn_site_configuration.example.site_id
  revision_id = qwiltcdn_site_configuration.example.revision_id
  #certificate_id = qwiltcdn_certificate.example.cert_id
}