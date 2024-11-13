
# Output to dump every attribute of every certificate available to users in your org.
output "all_origin_allow_list" {
  value = data.qwilt_cdn_origin_allow_list.detail
}
