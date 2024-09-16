# Examples

The examples in this directory are mostly used to generate documentation, with the exception of the files under Playground (see below).

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or are testable even if some parts are not relevant for the documentation.

* **provider/provider.tf** example file for the provider index page
* **data-sources/`full data source name`/data-source.tf** example file for the named data source page
* **resources/`full resource name`/resource.tf** example file for the named data source page

Run: go generate -v ./... to generate documentation


The sample configuration files under Playground demonstrate how to use the Qwilt Terraform provider. They can be used as starter files for provisioning and managing resources via the Terraform CLI. They are designed for customization-- replace placeholder values with your own specific configuration details. Replace the example certificate and key values with your own.
