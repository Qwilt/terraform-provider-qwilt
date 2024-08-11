// Package qwiltcdn
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.
package qwiltcdn

import (
	"context"
	cdnclient "github.com/Qwilt/terraform-provider-qwilt/qwilt/qwiltcdn/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	_ resource.Resource                = &siteActivationStagingResource{}
	_ resource.ResourceWithConfigure   = &siteActivationStagingResource{}
	_ resource.ResourceWithImportState = &siteActivationStagingResource{}
)

func NewSiteActivationStagingResource() resource.Resource {
	r := siteActivationStagingResource{}
	r.target = cdnclient.TARGET_STAGING
	return &r
}

// siteActivationStagingResource is the resource implementation.
type siteActivationStagingResource struct {
	siteActivationResource
}

// Metadata returns the resource type name.
func (r *siteActivationStagingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site_activation_staging"
}

// Schema defines the schema for the resource.
func (r *siteActivationStagingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.siteActivationResource.Schema(ctx, req, resp)
	resp.Schema.Attributes["target"] = schema.StringAttribute{
		Description: "The value will always be 'staging'.",
		Computed:    true,
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *siteActivationStagingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.siteActivationResource.Create(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *siteActivationStagingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.siteActivationResource.Read(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *siteActivationStagingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.siteActivationResource.Update(ctx, req, resp)
}

// Delete just removes the Terraform state on success. No real deletion for this object
func (r *siteActivationStagingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.siteActivationResource.Delete(ctx, req, resp)
}

// Configure adds the provider configured client to the resource.
func (r *siteActivationStagingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.siteActivationResource.Configure(ctx, req, resp)
}

func (r *siteActivationStagingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.siteActivationResource.ImportState(ctx, req, resp)
}
