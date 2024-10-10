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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/api"
	cdnclient "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/client"
	cdnmodel "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/model"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &siteConfigResource{}
	_ resource.ResourceWithConfigure   = &siteConfigResource{}
	_ resource.ResourceWithImportState = &siteConfigResource{}
	_ resource.ResourceWithModifyPlan  = &siteConfigResource{}
)

// NewSiteConfigResource is a helper function to simplify the provider implementation.
func NewSiteConfigResource() resource.Resource {
	return &siteConfigResource{}
}

// siteConfigResource is the resource implementation.
type siteConfigResource struct {
	client *cdnclient.SiteClientFacade
}

// Metadata returns the resource type name.
func (r *siteConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_site_configuration"
}

// Schema defines the schema for the resource.
func (r *siteConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Qwilt CDN site Configuration. [Learn how to prepare the configuration JSON.](https://api-docs.qwilt.cqloud.com/docs/CDN%20APIs/Sites%20API/prepare-the-configuration)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "For internal use only, for testing. Equals site_id:revision_id.",
				Computed:    true,
			},
			"site_id": schema.StringAttribute{
				Description: "The unique identifier of the Site.",
				Required:    true,
			},
			"revision_id": schema.StringAttribute{
				Description: "The unique identifier of the configuration version.",
				Computed:    true,
			},
			"revision_num": schema.Int64Attribute{
				Description: "The unique revision number of the configuration version.",
				Computed:    true,
			},
			"host_index": schema.StringAttribute{
				Description: "The SVTA metadata objects that define the delivery service configuration, in application/json format.",
				CustomType:  cdnmodel.HostIndexType{},
				Required:    true,
			},
			"change_description": schema.StringAttribute{
				Description: "Comments added by the user to the configuration JSON payload.",
				Required:    true,
			},
			"owner_org_id": schema.StringAttribute{
				Description: "The organization that owns the site.",
				Computed:    true,
			},
			"last_update_time_milli": schema.Int64Attribute{
				Description: "When the site configuration was last updated, in epoch time.",
				Computed:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *siteConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cdnmodel.SiteConfiguration
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	siteCreate := api.SiteConfigAddRequest{
		HostIndex:         json.RawMessage(plan.HostIndex.ValueString()),
		ChangeDescription: string(plan.ChangeDescription.ValueString()),
	}

	// Create new site
	siteResp, err := r.client.CreateSiteConfig(plan.SiteId.ValueString(), siteCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Qwilt CDN Site",
			"Could not create Qwilt CDN Site, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = cdnmodel.NewSiteConfigBuilder().
		WithCtx(ctx).
		WithSiteId(siteResp.SiteId).
		WithRevisionId(siteResp.RevisionId).
		WithHostIndex(siteCreate.HostIndex, false). //host index is not returned from QCon 'create'
		WithChangeDescription(siteResp.ChangeDescription).
		WithOwnerOrgId(siteResp.OwnerOrgId).
		WithRevisionNum(siteResp.RevisionNum).
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
func (r *siteConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state cdnmodel.SiteConfiguration
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed site value from CDN
	siteResp, err := r.client.GetSiteConfig(state.SiteId.ValueString(), state.RevisionId.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Qwilt CDN Site Configuration",
			"Could not read Qwilt CDN Site ID "+state.SiteId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state = cdnmodel.NewSiteConfigBuilder().
		WithSiteId(siteResp.SiteId).
		WithRevisionId(siteResp.RevisionId).
		WithHostIndex(siteResp.HostIndex, true).
		WithChangeDescription(siteResp.ChangeDescription).
		WithOwnerOrgId(siteResp.OwnerOrgId).
		WithRevisionNum(siteResp.RevisionNum).
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
func (r *siteConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan cdnmodel.SiteConfiguration
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	siteCreate := api.SiteConfigAddRequest{
		HostIndex:         json.RawMessage(plan.HostIndex.ValueString()),
		ChangeDescription: string(plan.ChangeDescription.ValueString()),
	}

	// Create new site - update is not supported for SiteConfiguration
	siteResp, err := r.client.CreateSiteConfig(plan.SiteId.ValueString(), siteCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Qwilt CDN Site",
			"Could not create Qwilt CDN Site, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = cdnmodel.NewSiteConfigBuilder().
		WithSiteId(siteResp.SiteId).
		WithRevisionId(siteResp.RevisionId).
		WithHostIndex(siteCreate.HostIndex, false).
		WithChangeDescription(siteResp.ChangeDescription).
		WithOwnerOrgId(siteResp.OwnerOrgId).
		WithRevisionNum(siteResp.RevisionNum).
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
func (r *siteConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state cdnmodel.SiteConfiguration
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No Deletion for existing site, just add a log
	params := map[string]any{"site_id": state.SiteId.ValueString(), "revision_id": state.RevisionNum}
	tflog.Info(ctx, "Qwilt CDN Site Deletion not supported: site_id=%s revision_num=%s", params)
}

// Configure adds the provider configured client to the resource.
func (r *siteConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cdnclient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *cdnclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = cdnclient.NewSiteFacadeClient(api.SITES_HOSTNAME, client)
}

func (r *siteConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	var site_id, revision_id string

	tflog.Info(ctx, "siteConfigResource:IMPORT: "+req.ID)
	//format:
	// option1: site_id:revision_id
	// option2: site_id (revision_id is defaulted to last active or latest)

	if len(idParts) == 2 && idParts[0] != "" && idParts[1] != "" {
		//option1: user wants to import a specific revision_id
		site_id = idParts[0]
		revision_id = idParts[1]
		tflog.Info(ctx, "import site "+site_id+" by explicit revision_id: "+revision_id)
	} else if len(idParts) == 1 && idParts[0] != "" {
		//option2: import latest or active revision_id
		site_id = idParts[0]

		tflog.Info(ctx, "import site "+site_id+" implicitly")

		//try to get latest or active revision
		siteResp, err := r.client.GetSite(idParts[0], "ga", false, false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Getting latest/active revision for Qwilt CDN Site",
				"Could not get active/latest revision for Qwilt CDN Site, unexpected error: "+err.Error(),
			)
			return
		}
		if siteResp.ActiveAndLastPublishingOperation != nil {
			tflog.Info(ctx, "retrieve publish info for site "+site_id)
			if siteResp.ActiveAndLastPublishingOperation.Active != nil && siteResp.ActiveAndLastPublishingOperation.Active.RevisionId != "" {
				revision_id = siteResp.ActiveAndLastPublishingOperation.Active.RevisionId
				tflog.Info(ctx, "found active revision id: "+revision_id+" for site "+site_id)
			} else if siteResp.ActiveAndLastPublishingOperation.Last != nil && siteResp.ActiveAndLastPublishingOperation.Last.RevisionId != "" {
				revision_id = siteResp.ActiveAndLastPublishingOperation.Last.RevisionId
				tflog.Info(ctx, "found last revision id: "+revision_id+" for site "+site_id)
			}
		} else {
			tflog.Info(ctx, "could not get publish info for site "+site_id+". use latest revision")

			//if we couldnt get active or latest, try to import the latest revision_num
			//try to get latest or active revision
			siteConfigsResp, err := r.client.GetSiteConfigs(idParts[0], true)
			if err != nil || len(siteConfigsResp) == 0 {
				resp.Diagnostics.AddError(
					"Error Getting site configuration revisions for Qwilt CDN Site",
					"Could not get site configuration revisions for Qwilt CDN Site, unexpected error: "+err.Error(),
				)
				return
			}
			tflog.Info(ctx, fmt.Sprintf("Found: %d configuration revisions for site %s", len(siteConfigsResp), site_id))
			maxRevisionNum := 0
			//find the latest revision_num and take it's revision_id
			for _, rev := range siteConfigsResp {
				if rev.RevisionNum > maxRevisionNum {
					maxRevisionNum = rev.RevisionNum
					revision_id = rev.RevisionId
				}
			}
		}
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: site_id, revision_id OR site_id. Got: %q", req.ID),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Importing: %s:%s", site_id, revision_id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("site_id"), site_id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("revision_id"), revision_id)...)
}

func (r *siteConfigResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {

	// A null plan means that the resource is being destroyed
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	// Get the plan
	var plan cdnmodel.SiteConfiguration
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the HostIndex is unknown, we can't do anything
	if plan.HostIndex.IsUnknown() {
		return
	}

	planRawJson := json.RawMessage([]byte(plan.HostIndex.ValueString()))
	planString, err := json.Marshal(planRawJson)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Validating Configured HostIndex",
			"Could not marshal configured plan HostIndex to JSON: "+err.Error(),
		)
		return
	}

	// Get the state
	var state cdnmodel.SiteConfiguration
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the HostIndex is unknown, we can't do anything
	if state.HostIndex.IsUnknown() {
		return
	}

	stateRawJson := json.RawMessage([]byte(state.HostIndex.ValueString()))
	stateString, err := json.Marshal(stateRawJson)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Validating HostIndex from State",
			"Could not marshal state HostIndex to JSON: "+err.Error(),
		)
		return
	}

	// Compare the plan and state hostIndex JSON
	hostIndexEqual, err := cdnmodel.JsonBytesEqual(planString, stateString)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Unmarshaling HostIndex for Comparison",
			"Could not compare plan and state HostIndex JSON: "+err.Error(),
		)
		return
	}

	// If ChangeDescription and HostIndex are both semantically equal,
	// i.e. no white space changes are detected, use values from the state
	// to suppress unnecessary updates.
	// Otherwise, allow all values to pass through unmodified.
	if (plan.ChangeDescription.ValueString() == state.ChangeDescription.ValueString()) &&
		hostIndexEqual {
		plan.HostIndex = state.HostIndex
		plan.RevisionId = state.RevisionId
		plan.RevisionNum = state.RevisionNum
		plan.OwnerOrgId = state.OwnerOrgId
		plan.LastUpdateTimeMilli = state.LastUpdateTimeMilli
		plan.Id = state.Id

		diags = resp.Plan.Set(ctx, &plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}
