# Output a cert template list if cert_template_id is "all" and cert_template_cn isn't defined.
output "certificate_templates_list" {
  value = var.cert_template_id == "all" && var.cert_template_cn == null ? [for certificate_template in data.qwilt_cdn_certificate_templates.detail.certificate_template : { common_name = certificate_template.common_name, certificate_template_id = certificate_template.certificate_template_id }] : null
}

# Output a matching cert template list if cert_template_cn is defined
output "certificate_templates_list_by_common_name" {
  value = var.cert_template_cn != null ? [for certificate_template in data.qwilt_cdn_certificate_templates.detail.certificate_template : { common_name = certificate_template.common_name, certificate_template_id = certificate_template.certificate_template_id } if strcontains(certificate_template.common_name, var.cert_template_cn)] : null
}

# Output the details of a matching cert template if cert_template_id is set to an explicit value.
output "certificate_template_detail" {
  value = var.cert_template_id == "all" || var.cert_template_id == null || data.qwilt_cdn_certificate_templates.detail.certificate_template == null ? null : data.qwilt_cdn_certificate_templates.detail.certificate_template[0]
}

# Output to dump every attribute of every cert template available to users in your org.
#output "all_certificate_templates" {
#    value = data.qwilt_cdn_certificate_templates.detail
#}
