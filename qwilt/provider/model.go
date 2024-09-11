// Package provider
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.

package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// QwiltProviderModel describes the provider data model.
type QwiltProviderModel struct {
	// CDN
	EnvType   types.String `tfsdk:"env_type"`
	Username  types.String `tfsdk:"username"`
	Password  types.String `tfsdk:"password"`
	XApiToken types.String `tfsdk:"api_key"`
}
