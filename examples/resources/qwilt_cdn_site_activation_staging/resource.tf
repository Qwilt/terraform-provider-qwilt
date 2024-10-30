
#Activation cannot be initiated while a previous activation is still 
#in progress.


resource "qwilt_cdn_site_activation_staging" "example" {
  site_id     = qwilt_cdn_site_configuration.example.site_id
  revision_id = qwilt_cdn_site_configuration.example.revision_id
  #certificate_id = qwilt_cdn_certificate.example.cert_id
}



