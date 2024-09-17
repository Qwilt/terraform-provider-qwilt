# Qwilt Provider Usage Examples

The examples in this directory demonstrate how to use the Qwilt Terraform Provider.  They cover various common use cases.

Inside this directory, you will find the following examples:
* **provider-install-verification**:  Demonstrates how to confirm that the Qwilt Terraform Provider is installed correctly.
* **qwilt-config-basic**:  Demonstrates how to configure a single site, site configuration, certificate, and activation resource.
* **qwilt-config-multi**: Demonstrates how to manage two sites, site configurations, certificates, and activation resources within a single Terraform configuration.
* **qwilt-config-multi-workspace**: (Advanced) Demonstrates how to use the Terraform Workspaces capability to manage multiple instances of the same configuration. Specifically, this example illustrates how to manage both production and non-production instances of a configuration that defines two sites.
* **qwilt-datasource**:  A practical example demonstrating how to use the *sites* and *certificates* data sources to query various resources.


**TIP:** The sample configuration files in this directory can also be used as starter files for provisioning and managing resources via the Terraform CLI. They are designed for customization-- replace placeholder values with your own specific configuration details. Replace the example certificate and key values with your own.