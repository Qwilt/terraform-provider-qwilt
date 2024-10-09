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
		MarkdownDescription: "The Qwilt Terraform provider is used to interact with the Qwilt Sites and Certificate Manager services. <br><br>" +
			"The Qwilt Terraform Provider allows you to manage your site configurations using a declarative configuration language. You can store, manage, and version your configuration data in the source control system of your choice.",
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
