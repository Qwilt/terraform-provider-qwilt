# Qwilt Multi-Site TF Config Example

This is a more complex example demonstrating how to manage two sites, site configurations, certificates, and activation resources in a single configuration.

Configurations are managed through separate Terraform configuration files (e.g. examplesite.tf and examplesite2.tf), and the SVTA configurations are managed through the corresponding JSON files (e.g. examplesite.json, examplesite2.json) that are included in the configuration.  Each site has its own associated certificate and private key.

To use this example:

1. Make sure your QCDN_API_KEY env variable is set. (This is the recommended authentication method.)

    For more information on authentication, see the provider documentation, which also covers alternative methods.

2. Replace the placeholder values with your own specific configuration details. 

    Replace the example certificate and key values with your own.

3. Apply the configuration:

    ```
    $ terraform apply
    ```