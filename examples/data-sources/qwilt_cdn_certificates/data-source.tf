
data "qwilt_cdn_certificates" "certificates_list" {
  filter = {
    cert_id = "all"
  }
}

data "qwilt_cdn_certificates" "detail" {
  filter = {
    cert_id = 25
  }
}