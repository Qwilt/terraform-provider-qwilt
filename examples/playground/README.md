# Qwilt Provider Usage Examples

This directory contains examples that comprehensively demonstrate how to use Qwilt's Terraform Provider.  These examples demonstrate various common use cases.

Inside this directory, you will find the following examples:
* *provider-install-verification*:  Very bare-bones example to verify that you have the Qwilt Terraform Provider installed correctly.
* *qwilt-config-basic*:  A simple example with 1 site, site configuration, certificate, and activation resource.
* *qwilt-config-multi*:  A more complex example involving 2 sites, site configurations, certificates, and activation resources.
* *qwilt-config-multi-workspace*:  An advanced demonstration for how to use the Terraform workspaces capability to manage multiple instances of the same configuration.  This manages 2 configurations, and allows you to create production and non-production instances.
* *qwilt-datasource*:  A practical demonstration of how to use the sites and certificates data sources to query for various resources.
