#Notes:
#- This resource takes a long time to fully apply.

#- Any attempt to apply site_activation with the same site_id might encounter 
#  a failure due to another publish operation in-progress.

#- Run terraform refresh to sync the state of this resource explicitly.


resource "qwilt_cdn_site_activation" "example" {
  site_id     = qwilt_cdn_site_configuration.example.site_id
  revision_id = qwilt_cdn_site_configuration.example.revision_id
  #certificate_id = qwilt_cdn_certificate.example.cert_id
}