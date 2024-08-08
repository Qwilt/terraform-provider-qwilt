# Qwilt Multi-Config Resource Example

This is a more complex example demonstrating how to manage 2 sites, site configurations, certificates, and activation resources in a single configuration.

Configurations are managed through separate Terraform configuration files (e.g. examplesite.tf and examplesite2.tf), and SVTA configurations are managed through corresponding JSON files that are included in the configuration.  Each site has its own associated certificate and private key.

To use it, define your API token and run 'apply':
```
$ terraform apply -var="token=<API Token>"
```