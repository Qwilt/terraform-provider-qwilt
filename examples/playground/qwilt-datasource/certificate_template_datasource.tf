variable "certificate_template_id" {
  type    = string
  default = null
}

variable "certificate_template_common_name" {
  type    = string
  default = null
}

data "qwilt_cdn_certificate_templates" "detail" {
  filter = {
    id = var.certificate_template_id != null && var.certificate_template_id != "all" ? try(tonumber(var.certificate_template_id), -1) : null
  }
}
