# Qwilt Basic Resource Example

This is a simple example demonstrating how to configure 1 site, site configuration, certificate, and activation resource.

The example by default uses the embedded host_index JSON configuration in main.tf.  Optionally, this can be removed, and "examplesitejson.json" can be included instead to maintain the SVTA configuration in a separate file.

To use it, define your API token and run 'apply':
```
$ terraform apply -var="token=<API Token>"
```