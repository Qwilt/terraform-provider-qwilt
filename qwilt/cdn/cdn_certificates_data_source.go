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
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	cdnclient "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/client"
)

var (
	_ datasource.DataSource              = &qwiltCertificatesDataSource{}
	_ datasource.DataSourceWithConfigure = &qwiltCertificatesDataSource{}
)

// NewCertificatesDataSource is a helper function to simplify the provider implementation.
func NewCertificatesDataSource() datasource.DataSource {
	return &qwiltCertificatesDataSource{}
}

// qwiltCertificatesDataSource is the data source implementation.
type qwiltCertificatesDataSource struct {
	client *cdnclient.CertificatesClient
}

// qwiltCertificatesDataSourceModel maps the data source schema data.
type qwiltCertificatesDataSourceModel struct {
	Cert   []cdnmodel.CertificateDataModel `tfsdk:"certificate"`
	Filter types.Object                    `tfsdk:"filter"`
}

// qwiltCertificatesFilterModel
type qwiltCertificatesFilterModel struct {
	CertId types.Int64 `tfsdk:"cert_id"`
}

// Metadata returns the data source type name.
func (d *qwiltCertificatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_certificates"
}

// Schema defines the schema for the data source.
func (d *qwiltCertificatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the certificates uploaded to Qwilt CDN by your organization and the associated metadata.",
		Attributes: map[string]schema.Attribute{
			"certificate": schema.ListNestedAttribute{
				Description: "List of certificates.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cert_id": schema.Int64Attribute{
							Description: "The unique identifier of the certificate. The certId will be needed when you add the certificate configuration and when you assign it to a site.",
							Computed:    true,
						},
						"certificate": schema.StringAttribute{
							Description: "A single X.509 signed PEM certificate, encoded in Base64.",
							Computed:    false,
							Required:    true,
						},
						"certificate_chain": schema.StringAttribute{
							Description: "An ordered concatenation of PEM-encoded signed certificates. The first is the signer of the imported certificate, and the last is an intermediate CA signed by a well known Root CA. The whole string must be Base64 encoded.",
							Computed:    false,
							Required:    true,
						},
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
							Description: "The certificate's type.",
							Computed:    true,
						},
					},
				},
			},
			"filter": schema.SingleNestedAttribute{
				Description: "Data source filter attributes.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"cert_id": schema.Int64Attribute{
						Description: "The ID of the specific certificate you want to retrieve.",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *qwiltCertificatesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state qwiltCertificatesDataSourceModel

	certs, err := d.client.GetCertificates(true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Qwilt Certificates",
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
	var filter qwiltCertificatesFilterModel
	var certIdFilter types.Int64

	if !state.Filter.IsNull() {
		diags := state.Filter.As(ctx, &filter, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		certIdFilter = filter.CertId
	}

	// Map response body to model
	// Get cert(s)
	for _, cert := range certs {
		if !certIdFilter.IsNull() && cert.CertId != certIdFilter.ValueInt64() {
			continue
		}
		certState := cdnmodel.CertificateDataModel{
			CertId:           types.Int64Value(cert.CertId),
			Certificate:      types.StringValue(cert.Certificate),
			CertificateChain: types.StringValue(cert.CertificateChain),
			//Email:            types.StringValue(cert.Email),
			Description: types.StringValue(cert.Description),
			PkHash:      types.StringValue(cert.PkHash),
			Tenant:      types.StringValue(cert.Tenant),
			Domain:      types.StringValue(cert.Domain),
			Status:      types.StringValue(cert.Status),
			Type:        types.StringValue(cert.Type),
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

func (d *qwiltCertificatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = cdnclient.NewCertificatesClient(client)
}
