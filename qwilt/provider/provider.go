// Package qwiltcdn
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.
package provider

import (
	"context"
	"fmt"
	"github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn"
	cdnclient "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/client"
	"github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/model"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
	"strings"
)

// Ensure qwiltCDNProvider satisfies various provider interfaces.
var _ provider.Provider = &qwiltCDNProvider{}

// NewQwiltProvider is a helper function to simplify provider server and testing implementation.
func NewQwiltProvider(version string) *qwiltCDNProvider {
	return &qwiltCDNProvider{
		version: version,
	}
}

// qwiltCDNProvider defines the provider implementation.
type qwiltCDNProvider struct {
	// version is set to the provider version on release, "dev" when the provider
	// is built and ran locally, and "test" when running acceptance testing.
	version string
}

func (p *qwiltCDNProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "qwilt"
	resp.Version = p.version
}

func (p *qwiltCDNProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	AddResponseSchema(resp)
}

func (p *qwiltCDNProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		cdn.NewCertificateResource,
		cdn.NewSiteActivationResource,
		cdn.NewSiteActivationStagingResource,
		cdn.NewSiteConfigResource,
		cdn.NewSiteResource,
	}
}

func (p *qwiltCDNProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		cdn.NewSitesDataSource,
		cdn.NewCertificatesDataSource,
	}
}

func (p *qwiltCDNProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Qwilt provider")

	// Retrieve provider data from configuration
	var config QwiltProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert provider configuration for Qwilt CDN.
	cfg := p.parseConfig(config)

	p.isConfigValid(cfg, resp)

	// Check for unknown errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	if cfg.EnvType == "" {
		cfg.EnvType = "prod"
	}

	logMsg := fmt.Sprintf("Creating Qwilt CDN API client with env_type: %s, xapi_token length: %d", cfg.EnvType, len(cfg.XApiToken))
	tflog.Debug(ctx, logMsg)

	client, err := cdnclient.NewClient(cfg.EnvType, cfg.Username, cfg.Password, cfg.XApiToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Qwilt CDN API Client",
			"An unexpected error occurred when creating the Qwilt API client. "+
				"If the error is not clear, please contact Qwilt customer support.\n\n"+
				"Qwilt API Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Qwilt CDN client", map[string]any{"success": true})
}

func (p *qwiltCDNProvider) parseConfig(config QwiltProviderModel) model.Settings {
	// Default values to environment variables, but override with Terraform
	// configuration value if set.
	cfg := model.Settings{
		EnvType:   os.Getenv("QCDN_ENVTYPE"),
		Username:  os.Getenv("QCDN_USERNAME"),
		Password:  os.Getenv("QCDN_PASSWORD"),
		XApiToken: os.Getenv("QCDN_XAPI_TOKEN"),
	}

	if !config.EnvType.IsNull() {
		cfg.EnvType = config.EnvType.ValueString()
	}

	if !config.Username.IsNull() {
		cfg.Username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		cfg.Password = config.Password.ValueString()
	}

	if !config.XApiToken.IsNull() {
		cfg.XApiToken = config.XApiToken.ValueString()
	}

	return cfg
}

func (p *qwiltCDNProvider) isConfigValid(cfg model.Settings, resp *provider.ConfigureResponse) bool {
	//token is the preferred authentication option.if token is unavailable use legacy user/pass authentication with explicit env_type
	if cfg.XApiToken == "" {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("token"),
			"Unknown QC Services login token",
			"The provider is missing value for the login token, which is the preferred authentication method. "+
				"It is recommended to Either set the value statically in the configuration, or use the QCDN_XAPI_TOKEN environment variable.",
		)
		if cfg.Username == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("username"),
				"Unknown QC Services username",
				"The provider cannot create the Qwilt CDN Sites API client as there is an unknown configuration value for the username. "+
					"Either set the value statically in the configuration, or use the QCDN_USERNAME environment variable.",
			)
		} else {
			if cfg.Password == "" {
				resp.Diagnostics.AddAttributeError(
					path.Root("password"),
					"Unknown QC Services password",
					"The provider cannot create the Qwilt CDN Sites API client as there is an unknown configuration value for the password. "+
						"Either set the value statically in the configuration, or use the QCDN_PASSWORD environment variable.",
				)
			}
			if strings.Contains(cfg.Username, "qwilt.com") {
				if cfg.EnvType == "" {
					resp.Diagnostics.AddAttributeError(
						path.Root("env_type"),
						"Unknown env_type",
						"The provider cannot create the Qwilt CDN Sites API client as there is an unknown configuration value for the env_type. "+
							"Either set the value statically in the configuration, or use the QCDN_ENVTYPE environment variable.",
					)
				}
			}
		}
	}
	return true
}
