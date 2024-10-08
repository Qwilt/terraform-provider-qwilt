# Qwilt Basic Resource Example

This is a simple example demonstrating how to configure a single site, site configuration, certificate, and activation resource.

In this example, the host_index JSON configuration is embedded in *main.tf*.  

If you prefer to maintain the SVTA configuration in a separate file, remove the embedded host_index, and instead use the reference to *examplesitebasic.json*: 

    ```
    host_index = file("./examplesitebasic.json")
    ```



To use this example:

1. Make sure your QCDN_API_KEY env variable is set. (This is the recommended authentication method.)

    For more information on authentication, see the [User Guide](https://docs.qwilt.com/docs/terraform-user-guide#authentication), which also covers alternative methods.

2. Replace the placeholder values with your own specific configuration details. 

    Replace the example certificate and key values with your own.

3. Apply the configuration:

    ```
    $ terraform apply
    ```