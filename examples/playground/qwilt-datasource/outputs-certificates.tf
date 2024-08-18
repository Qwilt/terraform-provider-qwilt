# Output a cert list if cert_id is "all" and cert_domain isn't defined
output "certificates_list" {
  value = var.cert_id == "all" && var.cert_domain == null ? [for cert in data.qwiltcdn_certificates.detail.certificate : { domain = cert.domain, cert_id = cert.cert_id, csr_id = cert.csr_id }] : null
}

# Output a matching cert list if cert_domain is defined
output "cert_list_by_cert_domain" {
  value = var.cert_domain != null ? [for cert in data.qwiltcdn_certificates.detail.certificate : { cert_domain = cert.domain, cert_id = cert.cert_id } if strcontains(cert.domain, var.cert_domain)] : null
}

# Output the detail of a matching cert if cert_id is set to an explicit value
output "cert_detail" {
  value = var.cert_id == "all" || var.cert_id == null || data.qwiltcdn_certificates.detail.certificate == null ? null : data.qwiltcdn_certificates.detail.certificate[0]
}

# Output to dump every attribute of every cert in the account (probably not what you want)
#output "all_certs" {
#    value = data.qwiltcdn_certificates.detail
#}
