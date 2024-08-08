
resource "qwiltcdn_certificate" "example" {
  certificate       = filebase64("./domain.crt")
  certificate_chain = filebase64("./domain.crt")
  private_key       = filebase64("./domain.key")
  description       = "ccc"
}

