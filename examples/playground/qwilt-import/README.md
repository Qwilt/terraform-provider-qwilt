# Qwilt Import Resource Example

This is a simple example demonstrating how to import 1 site, site configuration, certificate, and activation resource.

First, make sure your QCDN_API_KEY env variable is set (this is the recommended method for authentication.
See other authentication alternatives in details in the provider documentation.

## Basic Import

Basic import uses site_id to detect revison_id and publish_id automatically and spare the user the hussle of more complex commands.
The logic for qwilt_cdn_site_configuration is:
1. use published active revision_id (if available)
2. use published last revision_id (if available)
3. use latest revision_id (max revision_number)

The import is a rather manual operation. To complete the import process the following steps should be done for each resource:
1. perform terraform import
2. verify import completed successfully
3. use attributes from terraform state to generate the configurable attributes in the .tf file for this resource
4. change the id's to reference to allow implicit dependencies:
   1. site_id in site_configuration
   2. site_id, revision_id and cert_id in site_activation
5. run terraform plan - expect no changes 

To use it for step (1), define your API token and run:
```
$ terraform import qwilt_cdn_certificate.example <cert_id> -var="token=<API Token>"
$ terraform import qwilt_cdn_site.example <site_id> -var="token=<API Token>"
$ terraform import qwilt_cdn_site_configuration.example <site_id> -var="token=<API Token>"
$ terraform import qwilt_cdn_site_activation.example <site_id> -var="token=<API Token>"
```

## Advanced Import

Advanced import lets the user explicitly select revision_id for site_configuration and publish_id for site_activation

The import is a rather manual operation. To complete the import process the following steps should be done for each resource:
1. perform terraform import
2. verify import completed successfully
3. use attributes from terraform state to generate the configurable attributes in the .tf file for this resource
4. change the id's to reference to allow implicit dependencies:
   1. site_id in site_configuration
   2. site_id, revision_id and cert_id in site_activation
5. run terraform plan - expect no changes

To use it for step (1), define your API token and run:
```
$ terraform import qwilt_cdn_certificate.example <cert_id> -var="token=<API Token>"
$ terraform import qwilt_cdn_site.example <site_id> -var="token=<API Token>"
$ terraform import qwilt_cdn_site_configuration.example <site_id>:<revision_id> -var="token=<API Token>"
$ terraform import qwilt_cdn_site_activation.example <site_id>:<publish_id> -var="token=<API Token>"
```