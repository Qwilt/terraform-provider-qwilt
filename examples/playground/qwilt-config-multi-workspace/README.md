# Qwilt Multi-Config Resource Example Using Workspaces

This is our most complex example.  If you are not comfortable with how to use the Qwilt Terraform Provider, we suggest reviewing some of the other examples first.

This is an advanced demonstration for how to use the Terraform workspaces capability to manage multiple instances of the same configuration.  This example manages 2 configurations, and allows you to create production and non-production instances.  This strategy may be used to manage more complex workflows, where different teams (e.g. dev, qa, etc.) may need to work on different instances of a site.

This example may be more complex than your individual needs.  Qwilt's system provides a method for simple staging of sites for verification prior to production that does not require workspaces.  If your needs are more modest, we encourage you to try this staging capability instead.

In this example, configurations are managed through separate Terraform configuration files (e.g. examplesite.tf and examplesite2.tf), and SVTA configurations are managed through corresponding JSON files that are included in the configuration.  Both the configuration and JSON configuration are templated to allow variable substitution based on what type of site is being managed.  Each site and instance has its own associated certificate and private key.

First, make sure your QCDN_XAPI_TOKEN env variable is set (this is the recommended method for authentication.
See other authentication alternatives in details in the provider documentation.

To get started, you must create your workspaces.  For example:
```
$ terraform workspace new prod
Created and switched to workspace "prod"!

You're now on a new, empty workspace. Workspaces isolate their state,
so if you run "terraform plan" Terraform will not see any existing state
for this configuration.

$ terraform workspace new dev
Created and switched to workspace "dev"!

You're now on a new, empty workspace. Workspaces isolate their state,
so if you run "terraform plan" Terraform will not see any existing state
for this configuration.

$ terraform workspace list
  default
* dev
  prod
```

Terraform will now track separate states for each of these environments.  You may now switch to an environment to operate within it.  This will switch to the production workspace so that you can apply the production configuration:
```
$ terraform workspace select prod
Switched to workspace "prod".
```

Now, you may apply the configuration as follows:
```
$ terraform apply
```

The same may be done within the "dev" workspace.  Managing more than 2 workspaces (prod and non-prod) will require further changes to the configuration.