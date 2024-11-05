
resource "qwilt_cdn_certificate_template" "example" {
  common_name = "example.com"
  auto_managed_certificate_template = true
  country = "US"
  state = "New York"
  sans = ["www.example.com", "example.net"]
}

