
resource "qwilt_cdn_certificate_template" "example" {
  common_name = "example.com"
  auto_managed_certificate_template = true
}

