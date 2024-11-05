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

	cdnmodel "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/model"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	cdnclient "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/client"
)

var (
	_ datasource.DataSource              = &qwiltCertificateTemplatesDataSource{}
	_ datasource.DataSourceWithConfigure = &qwiltCertificateTemplatesDataSource{}
)

// NewCertificatesDataSource is a helper function to simplify the provider implementation.
func NewCertificateTemplatesDataSource() datasource.DataSource {
	return &qwiltCertificateTemplatesDataSource{}
}

// qwiltCertificateTemplatesDataSource is the data source implementation.
type qwiltCertificateTemplatesDataSource struct {
	client *cdnclient.CertificateTemplatesClient
}

// qwiltCertificateTemplatesDataSourceModel maps the data source schema data.
type qwiltCertificateTemplatesDataSourceModel struct {
	Cert   []cdnmodel.CertificateTemplateDataModel `tfsdk:"certificate"`
	Filter types.Object                            `tfsdk:"filter"`
}

// qwiltCertificatesFilterModel
type qwiltCertificateTemplatesFilterModel struct {
	CertificateTemplateId types.Int64 `tfsdk:"certificate_template_id"`
}

// Metadata returns the data source type name.
func (d *qwiltCertificateTemplatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_certificate_templates"
}

// Schema defines the schema for the data source.
func (d *qwiltCertificateTemplatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the certificate template defined by your organization and the associated metadata.",
		Attributes: map[string]schema.Attribute{
			"certificateTemplate": schema.ListNestedAttribute{
				Description: "List of certificate templates.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"autoManagedCertificateTemplate": schema.BoolAttribute{
							Description: "Indicates whether the certificate template is managed by Qwilt.",
							Required:    true,
						},
						"common_name": schema.StringAttribute{
							Description: "The fully qualified domain name for which the certificate is issued.",
							Required:    true,
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
						},
						"tenant": schema.StringAttribute{
							Description: "The organization your user is assigned to.",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "The full name of the state or province where the organization or entity requesting the certificate is located.",
							Computed:    false,
							Optional:    true,
						},
						"locality": schema.StringAttribute{
							Description: "The city or locality where the organization or entity requesting the certificate is located.",
							Computed:    false,
							Optional:    true,
						},
						"organization_name": schema.StringAttribute{
							Description: "The legal name of the organization or entity applying for the certificate.",
							Computed:    false,
							Optional:    true,
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
				},
			},
			"filter": schema.SingleNestedAttribute{
				Description: "Data source filter attributes.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"certificate_template_id": schema.Int64Attribute{
						Description: "The ID of the specific certificate template you want to retrieve.",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *qwiltCertificateTemplatesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state qwiltCertificateTemplatesDataSourceModel

	certs, err := d.client.GetCertificateTemplates()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Qwilt Certificate Templates",
			err.Error(),
		)
		return
	}

	// Get state
	diags := resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Filter variables
	var filter qwiltCertificateTemplatesFilterModel
	var idFilter types.Int64

	if !state.Filter.IsNull() {
		diags := state.Filter.As(ctx, &filter, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		idFilter = filter.CertificateTemplateId
	}

	// Map response body to model
	// Get cert(s)
	for _, cert := range certs {
		if !idFilter.IsNull() && cert.CertificateTemplateID != idFilter.ValueInt64() {
			continue
		}
		certState := cdnmodel.CertificateTemplateDataModel{
			CertificateTemplateId:          types.Int64Value(cert.CertificateTemplateID),
			Country:                        types.StringValue(cert.Country),
			Tenant:                         types.StringValue(cert.Tenant),
			State:                          types.StringValue(cert.State),
			Locality:                       types.StringValue(cert.Locality),
			OrganizationName:               types.StringValue(cert.OrganizationName),
			CommonName:                     types.StringValue(cert.CommonName),
			AutoManagedCertificateTemplate: types.BoolValue(cert.AutoManagedCertificateTemplate),
			LastCertificateID:              types.Int64Value(cert.LastCertificateID),
		}
		for _, san := range cert.SANs {
			certState.SANs = append(certState.SANs, types.StringValue(san))
		}

		for _, csr := range cert.CsrIds {
			certState.CsrIds = append(certState.CsrIds, types.Int64Value(csr))
		}

		state.Cert = append(state.Cert, certState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *qwiltCertificateTemplatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cdnclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *cdnclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = cdnclient.NewCertificateTemplateClient(client)
}
