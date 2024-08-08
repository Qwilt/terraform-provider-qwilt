variable "site_id" {
  type    = string
  default = null
}

variable "revision_id" {
  type    = string
  default = "all"
}

variable "publish_id" {
  type    = string
  default = "all"
}

variable "truncate_host_index" {
  type    = bool
  default = false
}

variable "site_name" {
  type    = string
  default = null
}

data "qwiltcdn_sites" "detail" {
  filter = {
    site_id             = var.site_id
    revision_id         = var.revision_id
    publish_id          = var.publish_id
    truncate_host_index = var.truncate_host_index
  }
}
