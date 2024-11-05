// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2024 Qwilt Inc.
package cdn

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/api"
	cdnclient "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/client"
	cdnmodel "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/model"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &certificateTemplateResource{}
	_ resource.ResourceWithConfigure   = &certificateTemplateResource{}
	_ resource.ResourceWithImportState = &certificateTemplateResource{}
)

// NewCertificateTemplateResource is a helper function to simplify the provider implementation.
func NewCertificateTemplateTemplateResource() resource.Resource {
	return &certificateTemplateResource{}
}

// certificateResource is the resource implementation.
type certificateTemplateResource struct {
	client *cdnclient.CertificateTemplatesClient
}

// Metadata returns the resource type name.
func (r *certificateTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_certificate_template"
}

// Schema defines the schema for the resource.
func (r *certificateTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Qwilt CDN CertificateTemplateTemplate.",
		Attributes: map[string]schema.Attribute{
			"autoManagedCertificateTemplate": schema.BoolAttribute{
				Description: "Indicates whether the certificate template is managed by Qwilt.",
				Required:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"common_name": schema.StringAttribute{
				Description: "The fully qualified domain name for which the certificate is issued.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate_template_id": schema.Int64Attribute{
				Description: "The unique identifier of the certificate template. This identifier will be needed when you add the certificate configuration and when you assign it to a site.",
				Computed:    true,
			},
			"country": schema.StringAttribute{
				Description: "The two-letter ISO 3166-1 country code that represents the country where the organization or entity requesting the certificate is located.",
				Computed:    false,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 2),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tenant": schema.StringAttribute{
				Description: "The organization your user is assigned to.",
				Computed:    true,
			},
			"state": schema.StringAttribute{
				Description: "The full name of the state or province where the organization or entity requesting the certificate is located.",
				Computed:    false,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"locality": schema.StringAttribute{
				Description: "The city or locality where the organization or entity requesting the certificate is located.",
				Computed:    false,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"organization_name": schema.StringAttribute{
				Description: "The legal name of the organization or entity applying for the certificate.",
				Computed:    false,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sans": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "Additional domains that the certificate should cover.",
				Computed:    false,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(100),
					listvalidator.UniqueValues(),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"last_certificate_id": schema.Int64Attribute{
				Description: "The unique identifier of the last certificate generated from this template.",
				Computed:    true,
			},
			"csr_ids": schema.ListAttribute{
				ElementType: types.Int64Type,
				Description: "The unique identifiers of the certificate signing requests (CSRs) that are associated with this certificate template.",
				Computed:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *certificateTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cdnmodel.CertificateTemplate
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	certRequest := api.CertificateTemplateCreateRequest{
		Country:                        plan.Country.ValueString(),
		State:                          plan.State.ValueString(),
		Locality:                       plan.Locality.ValueString(),
		OrganizationName:               plan.OrganizationName.ValueString(),
		CommonName:                     plan.CommonName.ValueString(),
		AutoManagedCertificateTemplate: plan.AutoManagedCertificateTemplate.ValueBool(),
	}

	for _, san := range plan.SANs {
		certRequest.SANs = append(certRequest.SANs, san.ValueString())
	}

	// Create new site
	certResp, err := r.client.CreateCertificateTemplate(certRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Qwilt CDN Certificate Template",
			"Could not create Qwilt CDN Certificate Template, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	newPlan := cdnmodel.NewCertificateTemplateBuilder().
		CertificateTemplateId(certResp.CertificateTemplateID).
		Country(certResp.Country).
		Tenant(certResp.Tenant).
		State(certResp.State).
		Locality(certResp.Locality).
		OrganizationName(certResp.OrganizationName).
		CommonName(certResp.CommonName).
		AutoManagedCertificateTemplate(certResp.AutoManagedCertificateTemplate).
		LastCertificateID(certResp.LastCertificateID).
		AddSANs(certResp.SANs...).
		AddCsrIds(certResp.CsrIds...).
		Build()

	// Set state to fully populated data
	diags = resp.State.Set(ctx, newPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *certificateTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state cdnmodel.CertificateTemplate
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed certificate value from client
	certResp, err := r.client.GetCertificateTemplate(state.CertificateTemplateId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Qwilt CDN CertificateTemplate",
			"Could not read Qwilt CDN CertificateTemplate ID "+state.CertificateTemplateId.String()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state = cdnmodel.NewCertificateTemplateBuilder().
		CertificateTemplateId(certResp.CertificateTemplateID).
		Country(certResp.Country).
		Tenant(certResp.Tenant).
		State(certResp.State).
		Locality(certResp.Locality).
		OrganizationName(certResp.OrganizationName).
		CommonName(certResp.CommonName).
		AutoManagedCertificateTemplate(certResp.AutoManagedCertificateTemplate).
		LastCertificateID(certResp.LastCertificateID).
		AddSANs(certResp.SANs...).
		AddCsrIds(certResp.CsrIds...).
		Build()

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Updates the resource and sets the updated Terraform state on success.
func (r *certificateTemplateResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

// Deletes the resource and removes the Terraform state on success.
func (r *certificateTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state cdnmodel.CertificateTemplate
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing site
	err := r.client.DeleteCertificateTemplate(state.CertificateTemplateId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Qwilt CDN CertificateTemplate",
			"Could not delete Qwilt CDN CertificateTemplate, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *certificateTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = cdnclient.NewCertificateTemplateClient(client)
}

func (r *certificateTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to CertId attribute
	importIdIntValue, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Missing Attribute Path",
			"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Resource ImportState method call to ImportStatePassthroughID path must be set to a valid attribute path that can accept a string value.",
		)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("certificate_template_id"), importIdIntValue)...)
	}
}
