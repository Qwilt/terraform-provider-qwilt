variable "cert_template_id" {
  type    = string
  default = null
}

variable "cert_template_cn" {
  type    = string
  default = null
}

data "qwilt_cdn_certificate_templates" "detail" {
  filter = {
    certificate_template_id = var.cert_template_id != null && var.cert_template_id != "all" ? try(tonumber(var.cert_template_id), -1) : null
  }
}
