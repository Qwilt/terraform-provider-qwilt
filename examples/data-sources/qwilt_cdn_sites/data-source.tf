
data "qwilt_cdn_sites" "example" {
  filter = {
    # Specify "all" for site_id, revision_id, or publish_id to list all instances.
    site_id             = "65fc70554a1c9c72079eb803"
    revision_id         = "65fc70561c17716a2c468839"
    publish_id          = "27ae501c-ac0f-430a-876c-42e6e78d2969"
    truncate_host_index = false
  }
}
