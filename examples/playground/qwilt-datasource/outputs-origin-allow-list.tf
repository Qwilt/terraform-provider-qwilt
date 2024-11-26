# Conditionally show the complete origin allow list.
output "origin_allow_list" {
  value = var.show_origin_allow_list == false ? null : data.qwilt_cdn_origin_allow_list.detail
}

# Unconditionally output the complete origin allow list.
#output "all_origin_allow_list" {
#  value = data.qwilt_cdn_origin_allow_list.detail
#}
