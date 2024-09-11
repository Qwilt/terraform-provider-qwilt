# Qwilt Data Source Example

This is a practical demonstration of how to use the sites and certificates data sources to query for various resources within your configuration.

First, make sure your QCDN_API_KEY env variable is set (this is the recommended method for authentication.
See other authentication alternatives in details in the provider documentation.

By default, checking the plan will return an empty result:
```
$ terraform plan
```

## Site Lookup

A common use case for using the data source is to find your site information so that you can import it into a Terraform configuration.  To search sites by name, you may do this:
```
$ terraform plan -var="site_name=example.com"
```
The command above will perform a substring match of the site_name attribute of each defined site and output matching site_id's.

If you prefer searching through a giant list of all sites, you may output all sites as follows:
```
$ terraform plan -var="site_id=all"
```

Once you have identified the site_id that you are targeting, you may search for it directly:
```
$ terraform plan -var="site_id=<SiteID>"
```
This will provide details about a specific site_id, and a list of associated site configurations (i.e. revisions) and publishing operations (i.e. site activations).

## Site Configuration and Activation Lookup

To get details about a specific site configuration and/or publishing operation, you may use their specific variables:
```
$ terraform plan -var="site_id=<SiteID>" -var="revision_id=<RevisionID>" -var="publish_id=<PublishID>"
```
You may specific either revision_id, publish_id, or both.

## Certificate Lookup

The same logic can be applied to certificates.  Certificates can be queried by domain name and certificate ID.

To search for certificates by domain name, you may do this:
```
$ terraform plan -var="cert_domain=example.com"
```
This command will perform a substring match of the cert_domain attribute for each defined certificate and output the matching cert_id's.

If you prefer searching through a giant list of all certificates, you may output all certificates as follows:
```
$ terraform plan -var="cert_id=all"
```

Once you have identified the cert_id that you are targeting, you may search for it directly:
```
$ terraform plan -var="cert_id=<CertificateID>"
```
This will provide details about a specific cert_id.

## About This Example

This example uses some complex output definitions to query the sites and certs data sources for information.  It illustrates how variables can be defined to search a specific attribute and filter the results.  With some additional effort, the same can be done to query other site, revision, activation, and certificate attributes.