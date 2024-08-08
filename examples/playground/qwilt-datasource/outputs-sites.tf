# Output a site list if site_id is "all" and site_name isn't defined
output "sites_list" {
  value = var.site_id == "all" && var.site_name == null ? [for site in data.qwiltcdn_sites.detail.site : { site_name = site.site_name, site_id = site.site_id }] : null
}

# Output a matching site list if site_name is defined
output "sites_list_by_site_name" {
  value = var.site_name != null ? [for site in data.qwiltcdn_sites.detail.site : { site_name = site.site_name, site_id = site.site_id } if strcontains(site.site_name, var.site_name)] : null
}

# Output the detail of a matching site if site_id is set to an explicit value
output "site_detail" {
  value = var.site_id == "all" || var.site_id == null || data.qwiltcdn_sites.detail.site == null ? null : data.qwiltcdn_sites.detail.site[0]
}

# Output to dump every attribute of every site in the account (probably not what you want)
#output "all_sites" {
#  value = data.qwiltcdn_sites.detail
#}
