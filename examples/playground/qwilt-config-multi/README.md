# Qwilt Multi-Config Resource Example

This is a more complex example demonstrating how to manage 2 sites, site configurations, certificates, and activation resources in a single configuration.

Configurations are managed through separate Terraform configuration files (e.g. examplesite.tf and examplesite2.tf), and SVTA configurations are managed through corresponding JSON files that are included in the configuration.  Each site has its own associated certificate and private key.

First, make sure your QCDN_XAPI_TOKEN env variable is set (this is the recommended method for authentication. 
See other authentication alternatives in details in the provider documentation.

To use it, run 'apply':
```
$ terraform apply
```