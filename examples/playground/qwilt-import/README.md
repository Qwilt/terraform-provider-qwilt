# Qwilt Import Resource Example

This is a simple example demonstrating how to import a single site, site configuration, certificate, and activation resource.

First, make sure your QCDN_API_KEY env variable is set. (This is the recommended authentication method.)
For more information on authentication, see the provider documentation, which also covers alternative methods.

## Basic Import

Basic import uses the site_id to automatically detect the revison_id and publish_id.

The logic for determining which qwilt_cdn_site_configuration to import is:
1. If there is an active published version, it is imported.
2. If there is no active published version, the most recently published version is imported.
3. If the site has never been published, the most recently saved configuration version is imported. (The version with the highest revision_Id value - max revision_number.)

To complete the import process, implement the following steps for each resource:
1. Make sure your QCDN_API_KEY env variable is set.
2. Perform terraform import. (See below.)
2. Verify the import was completed successfully.
3. Use the attribute values from the terraform state to generate the configurable attributes in the .tf file for the particular resource.
4. Change the following ids to references to allow for implicit dependencies:
   - *site_id* in the *site_configuration* resource.
   - *site_id*, *revision_id*, and *certificate_id* in the *site_activation* resource.
5. Run ```terraform plan``` - expect no changes.


The Terraform import command for each of the resources:

```
$ terraform import qwilt_cdn_certificate.example <cert_id> -var="token=<API Token>"

$ terraform import qwilt_cdn_site.example <site_id> -var="token=<API Token>"

$ terraform import qwilt_cdn_site_configuration.example <site_id> -var="token=<API Token>"

$ terraform import qwilt_cdn_site_activation.example <site_id> -var="token=<API Token>"
```

## Advanced Import

Advanced import lets the user explicitly select *revision_id* for site_configuration and *publish_id* for site_activation.


1. Make sure your QCDN_API_KEY env variable is set.
2. Perform terraform import. (See below.)
2. Verify the import was completed successfully.
3. Use the attribute values from the terraform state to generate the configurable attributes in the .tf file for the particular resource.
4. Change the following id's to references to allow for implicit dependencies:
   - *site_id* in the *site_configuration* resource.
   - *site_id*, *revision_id*, and *certificate_id* in the *site_activation* resource.
5. Run ```terraform plan``` - expect no changes.



The Terraform import command for each of the resources:
```
$ terraform import qwilt_cdn_certificate.example <cert_id> -var="token=<API Token>"

$ terraform import qwilt_cdn_site.example <site_id> -var="token=<API Token>"

$ terraform import qwilt_cdn_site_configuration.example <site_id>:<revision_id> -var="token=<API Token>"

$ terraform import qwilt_cdn_site_activation.example <site_id>:<publish_id> -var="token=<API Token>"
```