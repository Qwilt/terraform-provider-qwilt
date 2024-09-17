# Output the publish_ops list of a site if site_id is not "all", but revision_id and publish_id are not defined.
output "publish_ops_list" {
  value = var.site_id != "all" && data.qwilt_cdn_sites.detail.publish_op != null && var.publish_id == "all" && var.revision_id == "all" ? [for publish_op in data.qwilt_cdn_sites.detail.publish_op : publish_op] : null
}

# Output the publish_ops belonging to revision_id of a site if site_id and revision_id are defined, and publish_id is not defined. 
output "publish_ops_list_by_revision_id" {
  value = data.qwilt_cdn_sites.detail.publish_op != null && var.site_id != "all" && var.revision_id != "all" && var.publish_id == "all" ? [for publish_op in data.qwilt_cdn_sites.detail.publish_op : publish_op if publish_op.revision_id == var.revision_id] : null
}

# Output the publish_ops detail for an active job of a revision_id.
output "publish_ops_list_by_revision_id_active" {
  value = data.qwilt_cdn_sites.detail.publish_op != null && var.site_id != "all" && var.revision_id != "all" && var.publish_id == "all" ? [for publish_op in data.qwilt_cdn_sites.detail.publish_op : publish_op if publish_op.revision_id == var.revision_id && publish_op.is_active == true] : null
}

# Output the publish_ops detail of a matching operation if site_id and publish_id are set to explicit values.
output "publish_op_detail" {
  value = var.site_id == "all" || var.publish_id == "all" || data.qwilt_cdn_sites.detail.publish_op == null ? null : data.qwilt_cdn_sites.detail.publish_op[0]
}
