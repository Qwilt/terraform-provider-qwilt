// Package cdn
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.
package cdn

import (
	"context"
	"fmt"
	"github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/api"
	cdnclient "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/client"
	cdnmodel "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/model"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &siteResource{}
	_ resource.ResourceWithConfigure   = &siteResource{}
	_ resource.ResourceWithImportState = &siteResource{}
)

// NewSiteResource is a helper function to simplify the provider implementation.
func NewSiteResource() resource.Resource {
	return &siteResource{}
}

// siteResource is the resource implementation.
type siteResource struct {
	client *cdnclient.SiteClient
}

// Metadata returns the resource type name.
func (r *siteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_site"
}

// Schema defines the schema for the resource.
func (r *siteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Qwilt CDN site.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the site. Equals site_id. Required for testing infra",
				Computed:    true,
			},
			"site_id": schema.StringAttribute{
				Description: "The unique identifier of the site. The siteID will be needed when you add the site configuration and when you publish the site.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"site_dns_cname_delegation_target": schema.StringAttribute{
				Description: "The CNAME you'll use direct traffic from your website to the cdnclient.",
				Computed:    true,
			},
			"site_name": schema.StringAttribute{
				Description: "The user-defined site name.",
				Required:    true,
			},
			"routing_method": schema.StringAttribute{
				Description: "The routing method used for the site. It is defaulted to 'DNS'.",
				Computed:    true,
			},
			"owner_org_id": schema.StringAttribute{
				Description: "The organization that owns the site.",
				Computed:    true,
			},
			"last_update_time_milli": schema.Int64Attribute{
				Description: "When the site last updated, in epoch time.",
				Computed:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *siteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cdnmodel.Site
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "siteResource: create\n")

	// Generate API request body from plan
	siteCreate := api.SiteCreateRequest{
		SiteName: plan.SiteName.ValueString(),
	}

	// Create new site
	siteResp, err := r.client.CreateSite(siteCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Qwilt CDN Site",
			"Could not create Qwilt CDN Site, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = cdnmodel.SiteBuilder{}.
		SiteId(siteResp.SiteId).
		OwnerOrgId(siteResp.OwnerOrgId).
		SiteName(siteResp.SiteName).
		RoutingMethod(siteResp.RoutingMethod).
		SiteDnsCnameDelegationTarget(siteResp.SiteDnsCnameDelegationTarget).
		LastUpdateTimeMilli(siteResp.LastUpdateTimeMilli).
		Build()

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *siteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state cdnmodel.Site
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "siteResource: read\n")

	// Get refreshed site value from QC
	siteResp, err := r.client.GetSite(state.SiteId.ValueString(), "", false, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Qwilt CDN Site",
			"Could not read Qwilt CDN Site ID "+state.SiteId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state = cdnmodel.SiteBuilder{}.SiteId(siteResp.SiteId).
		OwnerOrgId(siteResp.OwnerOrgId).
		SiteName(siteResp.SiteName).
		RoutingMethod(siteResp.RoutingMethod).
		SiteDnsCnameDelegationTarget(siteResp.SiteDnsCnameDelegationTarget).
		LastUpdateTimeMilli(siteResp.LastUpdateTimeMilli).
		Build()

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *siteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan cdnmodel.Site
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "siteResource: update\n")

	// Retrieve values from state
	var state cdnmodel.Site
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	siteUpdate := api.SiteUpdateRequest{
		SiteName: plan.SiteName.ValueString(),
	}

	siteResp, err := r.client.UpdateSite(plan.SiteId.ValueString(), siteUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Qwilt CDN Site",
			"Could not update Qwilt CDN Site, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = cdnmodel.SiteBuilder{}.
		SiteId(siteResp.SiteId).
		OwnerOrgId(siteResp.OwnerOrgId).
		SiteName(siteResp.SiteName).
		RoutingMethod(state.RoutingMethod.ValueString()). //use the routing-method from the state in case it had an override in QC
		SiteDnsCnameDelegationTarget(siteResp.SiteDnsCnameDelegationTarget).
		LastUpdateTimeMilli(siteResp.LastUpdateTimeMilli).
		Build()

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete just removes the Terraform state on success. No real deletion for this object
func (r *siteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state cdnmodel.Site
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "siteResource: delete\n")

	// Delete the site - but rename it afterwards. This due to the fact that sites aren't really deleted but only marked for deletion,
	//so avoid future calls failure on name conflicts
	err := r.client.DeleteAndRenameSite(state.SiteId.ValueString(), state.SiteName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Qwilt CDN Site",
			"Could not deleting Qwilt CDN Site, unexpected error: "+err.Error(),
		)
		return
	}
}

// ure adds the provider configured client to the resource.
func (r *siteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cdnclient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source ure Type",
			fmt.Sprintf("Expected *qwiltcdnapi.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = cdnclient.NewSiteClient(api.SITES_HOSTNAME, client)
}

func (r *siteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	// Retrieve import ID and save to SiteId attribute
	path := path.Root("site_id")
	tflog.Info(ctx, "siteResource: ImportState "+path.String()+"\n")

	resource.ImportStatePassthroughID(ctx, path, req, resp)
}
