
data "qwiltcdn_certificates" "certificates_list" {
  filter = {
    cert_id = "all"
  }
}

data "qwiltcdn_certificates" "detail" {
  filter = {
    cert_id = 25
  }
}