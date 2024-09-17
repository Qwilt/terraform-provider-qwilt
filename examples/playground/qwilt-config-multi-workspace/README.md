# Qwilt Use Workspaces to Manage Multiple Instances of a TF Config Example

⚠️This is our most complex example.  If you are not yet comfortable using the Qwilt Terraform Provider, we suggest reviewing the other examples first.

This example demonstrates how to use Terraform workspaces to manage multiple instances of the same configuration. Specifically, it illustrates how to manage both production and non-production instances of a Terraform configuration that defines two sites.

This approach is ideal for managing advanced use cases where different teams (e.g. dev, qa, etc.) need to work on different instances of a site.

If this level of complexity is not needed for your use case, Qwilt provides a simpler method for staging sites for verification before production that does not require workspaces.  For more straightforward requirements, we recommend using the staging capability.

Configurations are managed through separate Terraform configuration files (e.g. examplesite.tf and examplesite2.tf), and the SVTA configurations are managed through the corresponding JSON files (e.g. examplesite.json, examplesite2.json). Both the Terraform and JSON configuration are templated to allow variable substitution based on the type of site being managed. Each site and instance has its own associated certificate and private key.

To use this example:

1. Make sure your QCDN_API_KEY env variable is set. (This is the recommended authentication method.)

    For more information on authentication, see the provider documentation, which also covers alternative methods. 

2. Create your workspaces.  

    For example, to create workspaces for production and development, run the following commands:

  

    ```
    $ terraform workspace new prod
    ```
      You created and switched to the workspace "prod."

      Since you're in a new, empty workspace and workspaces keep their state isolated, running terraform plan will not detect any existing state for this configuration.

    ```
    $ terraform workspace new dev
    ```

      You created and switched to the workspace "dev."

      Since you're in a new, empty workspace and workspaces keep their state isolated, running terraform plan will not detect any existing state for this configuration.

    ```
     $ terraform workspace list
       default
     * dev
       prod
    ```

    Terraform will now track separate states for each environment.  You must switch to an environment to operate within it. For example, to apply the production configuration, switch to the production workspace:

    ```
    $ terraform workspace select prod    
    ```

3. After switching to the desired workspace, apply the configuration:

    ```
    $ terraform apply
     ```

    Repeat for the "dev" workspace if needed. 
    
  Managing more than two workspaces (e.g. additional non-production environments) requires further changes to the configuration.