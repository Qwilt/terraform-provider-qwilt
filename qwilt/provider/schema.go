// Package provider
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

func AddResponseSchema(resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Qwilt Terraform Provider integrates with the Qwilt Sites API, Certificate Manager API, and Origin Allow List API, simplifying infrastructure management, allowing you to define configurations as code and apply them consistently across environments. <br><br>" +
			"[Qwilt Terraform Provider User Guide](https://docs.qwilt.com/docs/terraform-user-guide)<br><br>" +
			"**KEY FEATURES**<br>" +
			"**Site Resources** <br>" +
			"Create, update, and manage site configurations, including activation workflows.<br>" +
			"You can store, manage, and version your configuration data in the source control system of your choice.<br>" +
			"To manage your sites with the Qwilt Terraform provider, they must first be created in Terraform or imported into Terraform.<br>" +
			"[Learn how to create or import a site. ](https://docs.qwilt.com/docs/terraform-user-guide#create-or-import-a-site)<br><br>" +
			"**Certificate and Certificate Template Resources** <br>" +
			"Facilitate certificate management, including support for setting up Qwilt-managed Certificate Signing Request (CSR) workflows to streamline certificate lifecycle management.<br>" +
			"[Certificate Management in Terraform](https://docs.qwilt.com/docs/certificate-management-in-terraform) <br><br> " +
			"**Origin Allow List Data Source**<br>" +
			"Retrieves a comprehensive list of IP addresses that Qwilt CDN may use to fetch content from your origin, grouped by ISP network.<br><br>" +
			"**AUTHENTICATION**<br>" +
			"The Qwilt Terraform provider supports two authentication methods: API key-based authentication (the preferred method) and login with user name and password. You can set the authentication parameters inside the provider configuration or as environment variables. We recommend setting env variables.<br>" +
			"[Learn more about the supported authentication methods.](https://docs.qwilt.com/docs/terraform-user-guide#authentication)<br><br>" +
			"**QUICK START**<br>" +
			"The sample configuration files in our playground on GitHub demonstrate how to use the Qwilt Terraform provider. They can be used as starter files for provisioning and managing resources via the Terraform CLI. They are designed for customization-- replace placeholder values with your own configuration details. Replace the example certificate and key values with your own.<br>" +
			"[Explore the ReadMe files and examples in our playground on GitHub.](https://github.com/Qwilt/terraform-provider-qwilt/blob/main/examples/playground/README.md)",
		Attributes: map[string]schema.Attribute{
			// CDN
			"env_type": schema.StringAttribute{
				Description:         "Qwilt CDN environment [prod,prestg,stage,dev]. May also be set by the QCDN_ENVTYPE environment variable.",
				MarkdownDescription: "FOR INTERNAL USE ONLY!! The Qwilt CDN environment [prod,prestg,stage,dev]. May also be set by the QCDN_ENVTYPE environment variable.",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				Description:         "Username for Qwilt CDN Sites API. May also be provided via QCDN_USERNAME environment variable.",
				MarkdownDescription: "QC services username.  May also be set by the QCDN_USERNAME environment variable.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				Description:         "Password for Qwilt CDN Sites API. May also be provided via QCDN_PASSWORD environment variable.",
				MarkdownDescription: "QC services password. May also be set by the QCDN_PASSWORD environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"api_key": schema.StringAttribute{
				Description:         "API key for Qwilt CDN Sites API. May also be provided via QCDN_API_KEY environment variable.",
				MarkdownDescription: "API key for Qwilt CDN Sites API. May also be set by the QCDN_API_KEY environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}
