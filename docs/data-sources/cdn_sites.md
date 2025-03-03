---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "qwilt_cdn_sites Data Source - qwilt"
subcategory: ""
description: |-
  Retrieves the sites available to your user. By default, the output includes site metadata, details about the associated site configuration versions, and details about the associated publishing operations. You can apply filters to the data source.
---

# qwilt_cdn_sites (Data Source)

Retrieves the sites available to your user. By default, the output includes site metadata, details about the associated site configuration versions, and details about the associated publishing operations. You can apply filters to the data source.

## Example Usage

```terraform
data "qwilt_cdn_sites" "example" {
  filter = {
    # Specify "all" for site_id, revision_id, or publish_id to list all instances.
    site_id             = "65fc70554a1c9c72079eb803"
    revision_id         = "65fc70561c17716a2c468839"
    publish_id          = "27ae501c-ac0f-430a-876c-42e6e78d2969"
    truncate_host_index = false
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Attributes) Data source filter attributes. (see [below for nested schema](#nestedatt--filter))

### Read-Only

- `publish_op` (Attributes List) List of publishing operations associated with site. (see [below for nested schema](#nestedatt--publish_op))
- `revision` (Attributes List) Site configurations associated with the site. (see [below for nested schema](#nestedatt--revision))
- `site` (Attributes List) Site metadata. (see [below for nested schema](#nestedatt--site))

<a id="nestedatt--filter"></a>
### Nested Schema for `filter`

Optional:

- `publish_id` (String) Filter publishing operations based on a specific publish ID.
- `revision_id` (String) Filter configurations based on a specific revision ID.
- `site_id` (String) Filter sites based on a specific site ID.
- `truncate_host_index` (Boolean) By default, the configuration details are included in the response, and you can exclude them by setting this to false.


<a id="nestedatt--publish_op"></a>
### Nested Schema for `publish_op`

Read-Only:

- `creation_time_milli` (Number) The time when the publish operation was created, in epoch time.
- `is_active` (Boolean) Indicates if the configuration is active or inactive.
- `last_update_time_milli` (Number) When the publishing operation was last updated, in epoch time.
- `operation_type` (String) The operation type (Publish, Unpublish) that was initiated with the request. An Unpublish operation removes a delivery service from the CDN.
- `owner_org_id` (String) The organization that owns the site.
- `publish_id` (String) Unique identifier of the publishing operation.
- `publish_state` (String) For internal use. Use the 'publishStatus' values instead.
- `publish_status` (String) The publishing operation status.
- `revision_id` (String) Unique identifier of the configuration version that was published or unpublished.
- `status_line` (List of String) Additional information related to the publish status.
- `target` (String) The value will 'ga' or 'staging'.
- `username` (String) Username that initiated the publishing operation.
- `validators_err_details` (String) Details about errors generated during validation.


<a id="nestedatt--revision"></a>
### Nested Schema for `revision`

Read-Only:

- `change_description` (String) Comments added by the user to the configuration JSON payload.
- `created_user` (String) The user who created the site.
- `creation_time_milli` (Number) The time when the configuration version was added, in epoch time.
- `host_index` (String) The SVTA metadata objects that define the delivery service configuration, in application/json format.
- `last_update_time_milli` (Number) The time when the configuration version was added, in epoch time. (This will be the same as the creationTimeMilli value.)
- `owner_org_id` (String) The name of the organization that owns the site.
- `revision_id` (String) The unique identifier of the configuration version.
- `revision_num` (Number) The unique revision number of the configuration version.
- `site_id` (String) The unique identifier of the site.


<a id="nestedatt--site"></a>
### Nested Schema for `site`

Read-Only:

- `api_version` (String) The Media Delivery Configuration API version.
- `created_user` (String) The user who created the site.
- `creation_time_milli` (Number) When the site was created, in epoch time.
- `is_deleted` (Boolean) Indicates if site was marked for deletion.
- `is_self_service_blocked` (Boolean) Indicates if site updates are currently allowed.
- `last_update_time_milli` (Number) When the site was last updated, in epoch time.
- `last_updated_user` (String) The user who last updated the site.
- `owner_org_id` (String) The name of the organization that owns the site.
- `routing_method` (String) The routing method used for the site. DNS is the default routing mechanism.
- `service_id` (String) A system-generated unique identifier for the site configuration.
- `service_type` (String) The value will always be MEDIA_DELIVERY.
- `should_provision_to_third_party_cdn` (Boolean) Indicates if the site should be provisioned to third-party CDNs.
- `site_dns_cname_delegation_target` (String) The CNAME you'll use direct traffic from your website to the CDN.
- `site_id` (String) The unique identifier of the site. The siteID will be needed when you add the site configuration and when you publish the site.
- `site_name` (String) The user-defined site name.
