variable "show_origin_allow_list" {
  type    = bool
  default = false
}

data "qwilt_cdn_origin_allow_list" "detail" {
}
