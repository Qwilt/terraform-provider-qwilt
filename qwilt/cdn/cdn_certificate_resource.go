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
	"strconv"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &certificateResource{}
	_ resource.ResourceWithConfigure   = &certificateResource{}
	_ resource.ResourceWithImportState = &certificateResource{}
)

// NewCertificateResource is a helper function to simplify the provider implementation.
func NewCertificateResource() resource.Resource {
	return &certificateResource{}
}

// certificateResource is the resource implementation.
type certificateResource struct {
	client *cdnclient.CertificatesClient
}

// Metadata returns the resource type name.
func (r *certificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_certificate"
}

// Schema defines the schema for the resource.
func (r *certificateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Qwilt CDN Certificate.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The unique identifier of the site. Equals cert_id. Required for testing infra",
				Computed:    true,
			},
			"cert_id": schema.Int64Attribute{
				Description: "The unique identifier of the certificate. The certId will be needed when you add the certificate configuration and when you assign it to a site.",
				Computed:    true,
				//PlanModifiers: []planmodifier.Int64{
				//	setplanmodifier.UseStateForUnknown(),
				//},
			},
			"certificate": schema.StringAttribute{
				Description: "A single X.509 signed PEM certificate, encoded in Base64.",
				Computed:    false,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate_chain": schema.StringAttribute{
				Description: "An ordered concatenation of PEM-encoded signed certificates. The first is the signer of the imported certificate, and the last is an intermediate CA signed by a well known Root CA. The whole string must be Base64 encoded.",
				Computed:    false,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"private_key": schema.StringAttribute{
				Description: "A PEM private key, which pairs with the public key that is embedded in the certificate. The entire string must be Base64 encoded.",
				Computed:    false,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			//"email": schema.StringAttribute{
			//	Description: "The email of the user owning the certificate.",
			//	Computed:    false,
			//	Optional:    true,
			//},
			"description": schema.StringAttribute{
				Description: "The certificate description.",
				Computed:    false,
				Optional:    true,
			},
			"pk_hash": schema.StringAttribute{
				Description: "A unique identifier for the private key that does not expose the actual key itself.",
				Computed:    true,
			},
			"tenant": schema.StringAttribute{
				Description: "The organization your user is assigned to.",
				Computed:    true,
			},
			"domain": schema.StringAttribute{
				Description: "The site host domain the certificate is linked to.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the certificate:[\"ISSUED\",\n          \"ACTIVE\",\n          \"EXPIRED\",\n          \"REVOKED\"].",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The certificate type.",
				Computed:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *certificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cdnmodel.Certificate
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	certRequest := api.CertificateCreateRequest{
		PrivateKey:       plan.PrivateKey.ValueString(),
		Certificate:      plan.Certificate.ValueString(),
		CertificateChain: plan.CertificateChain.ValueString(),
		Description:      plan.Description.ValueString(),
		//Email:            plan.Email.ValueString(),
	}

	// Create new site
	certResp, err := r.client.CreateCertificate(certRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Qwilt CDN Certificate",
			"Could not create Qwilt CDN Certificate, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	newPlan := cdnmodel.CertificateBuilder{}.
		CertificateId(certResp.CertId).
		CertificateChain(plan.CertificateChain.ValueString()).
		Certificate(plan.Certificate.ValueString()).
		Description(certResp.Description).
		PrivateKey(plan.PrivateKey.ValueString()).
		Type(certResp.Type).
		Status(certResp.Status).
		PkHash(certResp.PkHash).
		Tenant(certResp.Tenant).
		Domain(certResp.Domain).
		Build()

	// Set state to fully populated data
	diags = resp.State.Set(ctx, newPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *certificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state cdnmodel.Certificate
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed certificate value from client
	certResp, err := r.client.GetCertificate(state.CertId, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Qwilt CDN Certificate",
			"Could not read Qwilt CDN Certificate ID "+state.CertId.String()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state = cdnmodel.CertificateBuilder{}.
		CertificateId(certResp.CertId).
		CertificateChain(certResp.CertificateChain).
		Certificate(certResp.Certificate).
		Description(certResp.Description).
		PrivateKey(state.PrivateKey.ValueString()).
		Type(certResp.Type).
		Status(certResp.Status).
		PkHash(certResp.PkHash).
		Tenant(certResp.Tenant).
		Domain(certResp.Domain).
		Build()

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *certificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan cdnmodel.Certificate
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from state
	var state cdnmodel.Certificate
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	certRequest := api.CertificateUpdateRequest{
		PrivateKey:       plan.PrivateKey.ValueString(),
		Certificate:      plan.Certificate.ValueString(),
		CertificateChain: plan.CertificateChain.ValueString(),
		Description:      plan.Description.ValueString(),
		//Email:            plan.Email.ValueString(),
	}

	// Create new site
	certResp, err := r.client.UpdateCertificate(state.CertId.ValueInt64(), certRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Qwilt CDN Certificate",
			"Could not update Qwilt CDN Certificate, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = cdnmodel.CertificateBuilder{}.
		CertificateId(certResp.CertId).
		CertificateChain(plan.CertificateChain.ValueString()).
		Certificate(plan.Certificate.ValueString()).
		Description(certResp.Description).
		PrivateKey(plan.PrivateKey.ValueString()).
		Type(certResp.Type).
		Status(certResp.Status).
		PkHash(certResp.PkHash).
		Tenant(certResp.Tenant).
		Domain(certResp.Domain).
		Build()

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *certificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state cdnmodel.Certificate
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing site
	err := r.client.DeleteCertificate(state.CertId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Qwilt CDN Certificate",
			"Could not delete Qwilt CDN Certificate, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *certificateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = cdnclient.NewCertificatesClient(client)
}

func (r *certificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to CertId attribute
	importIdIntValue, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Missing Attribute Path",
			"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Resource ImportState method call to ImportStatePassthroughID path must be set to a valid attribute path that can accept a string value.",
		)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cert_id"), importIdIntValue)...)
	}
}
