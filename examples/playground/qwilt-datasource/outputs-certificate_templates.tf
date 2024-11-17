# Output a cert list if cert_id is "all" and cert_domain isn't defined.
output "certificate_templates_list" {
  value = var.certificate_template_id == "all" && var.certificate_template_common_name == null ? [for certificate_template in data.qwilt_cdn_certificate_templates.detail.certificate_template : { common_name = certificate_template.domain, id = certificate_template.cert_id }] : null
}

# Output a matching cert list if cert_domain is defined
output "certificate_templates_list_by_common_name" {
  value = var.certificate_template_common_name != null ? [for certificate_template in data.qwilt_cdn_certificate_templates.detail.certificate_template : { common_name = certificate_template.common_name, id = certificate_template.certificate_template_id } if strcontains(certificate_template.common_name, var.certificate_template_common_name)] : null
}

# Output the details of a matching cert if cert_id is set to an explicit value.
output "certificate_template_detail" {
  value = var.certificate_template_id == "all" || var.certificate_template_id == null || data.qwilt_cdn_certificate_templates.detail.certificate_template == null ? null : data.qwilt_cdn_certificate_templates.detail.certificate_template[0]
}

# Output to dump every attribute of every certificate available to users in your org.
#output "all_certificate_templates" {
#    value = data.qwilt_cdn_certificate_templates.detail
#}
