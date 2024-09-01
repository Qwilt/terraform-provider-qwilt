#> ⚠️ **Disclaimer**: this resource is experimental and should not be used!

#Notes
#this resource is similar to site_activation resource but activates the selected configuration to staging segment only.
#Note that activation cannot be performed if a previous activation is still in-progress


resource "qwilt_cdn_site_activation_staging" "example" {
  site_id     = qwilt_cdn_site_configuration.example.site_id
  revision_id = qwilt_cdn_site_configuration.example.revision_id
  #certificate_id = qwilt_cdn_certificate.example.cert_id
}