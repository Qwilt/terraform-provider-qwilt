# Qwilt Data Source Example

This is a practical example demonstrating how to use the *sites* and *certificates* data sources to query various resources within your configuration.

First, make sure the QCDN_API_KEY env variable is set. (This is the recommended authentication method.)
For more information on authentication, see the [User Guide](https://docs.qwilt.com/docs/terraform-user-guide#authentication), which also covers alternative methods.

By default, checking the plan will return an empty result:
```
$ terraform plan
```

## Site Lookup

A common use case for using the data source is to find your site information so that you can import it into a Terraform configuration.  

To search for a site by name, do this:

```
$ terraform plan -var="site_name=example.com"
```

The above command launches a substring match of the site_name attribute of each defined site and outputs the matching site_ids.


To retrieve a list of all sites, do this:

```
$ terraform plan -var="site_id=all"
```

Once you have identified the site_id that you are targeting, you may search for it directly:

```
$ terraform plan -var="site_id=<SiteID>"
```
   
The output includes details about the specified site_id, and a list of associated site configurations (i.e. revisions) and publishing operations (i.e. site activations).

## Site Configuration and Activation Lookup

To get details about a specific site configuration and/or publishing operation, you may use their specific variables:
```
$ terraform plan -var="site_id=<SiteID>" -var="revision_id=<RevisionID>" -var="publish_id=<PublishID>"
```
You may specify revision_id, publish_id, or both.

## Certificate Lookup

The same logic can be applied to certificates.  Certificates can be queried by domain name and certificate ID.

To search for certificates by domain name:

```
$ terraform plan -var="cert_domain=example.com"
```
The above command launches a substring match of the cert_domain attribute for each defined certificate and outputs the matching cert_ids.


To retrieve a list of all certificates, do this:

```
$ terraform plan -var="cert_id=all"
```

Once you have identified the cert_id that you are targeting, you may search for it directly:

```
$ terraform plan -var="cert_id=<CertificateID>"
```
The output provides details about the specified cert_id.

## Certificate Template Lookup

The certificate templates data source can be used to query the details of certificate templates for the CSR workflow.  Certificate templates can be queried by common name and template ID.

To search for certificates by common name:

```
$ terraform plan -var="cert_template_cn=example.com"
```
The above command launches a substring match of the common_name attribute for each defined certificate template and outputs the matching certificate_template_ids.

To retrieve a list of all certificate templates, do this:

```
$ terraform plan -var="cert_template_id=all"
```

Once you have identified the certificate_template_id that you are targeting, you may search for it directly:

```
$ terraform plan -var="cert_template_id=<CertificateTemplateID>"
```
The output provides details about the specified certificate_template_id.

## Origin Allow List

The origin allow list data source allows you to view the most current origin IP allow list.  This data can be used, for example, to update the firewall rules of your origin server.

To view the most current origin allow list, do this:

```
$ terraform plan -var="show_origin_allow_list=true"
```
The output provides details about current IPv4 and IPv6 addresses associated with each service provider in Qwilt's network.

## About This Example

This example uses some complex output definitions to query the sites, certs, cert templates, and origin allow list data sources for information.  It illustrates how variables can be defined to search a specific attribute and filter the results.  With some additional effort, the same can be done to query other site, revision, activation, certificate, certificate template, and origin allow list attributes.