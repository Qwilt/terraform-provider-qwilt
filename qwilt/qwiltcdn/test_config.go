package qwiltcdn

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// QwiltCdnProviderConfig is a shared configuration to combine with the actual
	// test configuration so the CDN client is properly configured.
	// It is also possible to use the QCDN_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	QwiltCdnProviderConfig = `
provider "qwiltcdn" {
	env_type = "dev"
}
`

	// QwiltCdnProviderConfig is a shared configuration to combine with the actual
	// test configuration so the CDN client is properly configured.
	// It is also possible to use the QCDN_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	QwiltCdnFullProviderConfig = `
terraform {
  required_providers {
    qwiltcdn = {
      source = "qwilt.com/qwiltinc/qwilt"
    }
  }
}

provider "qwiltcdn" {
	env_type = "dev"
}
`
)

var (

	// TestAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"qwiltcdn": providerserver.NewProtocol6WithError(NewCdnProvider("test"))}
)
