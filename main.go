package main

import (
	"context"
	"flag"
	cdn "github.com/Qwilt/terraform-provider-qwilt/qwilt/provider"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name qwilt

var (

	//- use this when developing locally
	//providerName = "qwilt.com/qwiltinc/qwilt"

	//use this when testing in released version
	providerName = "Qwilt/qwilt"

	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debugFlag bool

	flag.BoolVar(&debugFlag, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()

	providers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(cdn.NewQwiltProvider(version)),
		//providerserver.NewProtocol6(oeprovider.NewQwiltProvider(version)),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err.Error())
	}

	var serveOpts []tf6server.ServeOpt

	if debugFlag {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(providerName, muxServer.ProviderServer, serveOpts...)
	if err != nil {
		log.Fatal(err.Error())
	}
}
