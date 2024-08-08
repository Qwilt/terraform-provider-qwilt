variable "cert_id" {
  type    = string
  default = null
}

variable "cert_domain" {
  type    = string
  default = null
}

data "qwiltcdn_certificates" "detail" {
  filter = {
    cert_id = var.cert_id != null && var.cert_id != "all" ? try(tonumber(var.cert_id), -1) : null
  }
}
