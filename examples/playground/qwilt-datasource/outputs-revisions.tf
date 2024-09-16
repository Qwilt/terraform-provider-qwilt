# Output the revision list of a site if site_id is not "all", and revision_id is not defined
output "revisions_list" {
  value = var.site_id != "all" && data.qwilt_cdn_sites.detail.revision != null && var.revision_id == "all" ? [for revision in data.qwilt_cdn_sites.detail.revision : { revision_id = revision.revision_id, revision_num = revision.revision_num, change_description = revision.change_description }] : null
}

# Output the revision detail of a matching revision if site_id and revision_id are set to explicit values
output "revision_detail" {
  value = var.site_id == "all" || var.revision_id == "all" || data.qwilt_cdn_sites.detail.revision == null ? null : data.qwilt_cdn_sites.detail.revision[0]
}
