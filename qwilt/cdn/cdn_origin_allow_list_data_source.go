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
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	cdnclient "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/client"
)

var (
	_ datasource.DataSource              = &qwiltOriginAllowListDataSource{}
	_ datasource.DataSourceWithConfigure = &qwiltOriginAllowListDataSource{}
)

// NewOriginAllowListDataSource is a helper function to simplify the provider implementation.
func NewOriginAllowListDataSource() datasource.DataSource {
	return &qwiltOriginAllowListDataSource{}
}

// qwiltOriginAllowListDataSource is the data source implementation.
type qwiltOriginAllowListDataSource struct {
	client *cdnclient.DeviceIpsClient
}

// qwiltOriginAllowListDataSourceModel maps the data source schema data.
type qwiltOriginAllowListDataSourceModel struct {
	IpData           map[string]cdnmodel.NetworkIpDataModel `tfsdk:"ip_data"`
	Md5              types.String                           `tfsdk:"md5"`
	CreateTimeMillis types.Int64                            `tfsdk:"create_time_millis"`
}

// Metadata returns the data source type name.
func (d *qwiltOriginAllowListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_origin_allow_list"
}

// Schema defines the schema for the data source.
func (d *qwiltOriginAllowListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a comprehensive list of IP addresses that Qwilt CDN may use to request content from your origin, grouped by ISP network.<br><br>" +
			"This list includes the IPs from all the networks where Qwilt CDN is present. Because of Qwilt's unique routing strategy, any of the IPs from any of the ISP networks may potentially send a request to your origin.<br><br>" +
			"If your security policies limit origin access to approved IPs only, make sure to add all of these IP addresses to your allow list.<br><br>" +
			"The list is verified by Qwilt once per hour.",
		Attributes: map[string]schema.Attribute{
			"md5": schema.StringAttribute{
				Description: "A unique identifier for this instance of the IP address list.",
				Computed:    true,
			},
			"create_time_millis": schema.Int64Attribute{
				Description: "The time this instance of the IP address list was generated.",
				Computed:    true,
			},
			"ip_data": schema.MapNestedAttribute{
				Description: "A dictionary structure where each key is a network name, and the value is an object comprised of two arrays; one for the IPv4 addresses and one for the IPv6 addresses in the network that the Qwilt CDN may use to request content from your origin.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ipv4": schema.ListAttribute{
							ElementType: types.StringType,
							Description: "The IPv4 addresses in the network that Qwilt CDN may use to request content from your origin.",
							Computed:    true,
						},
						"ipv6": schema.ListAttribute{
							ElementType: types.StringType,
							Description: "The IPv6 addresses in the network that Qwilt CDN may use to request content from your origin.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *qwiltOriginAllowListDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state qwiltOriginAllowListDataSourceModel

	deviceIpsResp, err := d.client.GetOriginAllowList()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Qwilt IpData",
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

	// Map response body to model
	state.Md5 = types.StringValue(deviceIpsResp.Md5)
	state.CreateTimeMillis = types.Int64Value(int64(deviceIpsResp.CreateTimeMillis))

	ipData := make(map[string]cdnmodel.NetworkIpDataModel, 0)

	for network, networkDeviceIps := range deviceIpsResp.IpData {
		entry := cdnmodel.NetworkIpDataModel{}
		for _, ip := range networkDeviceIps.Ipv4 {
			entry.Ipv4 = append(entry.Ipv4, types.StringValue(ip))
		}
		for _, ip := range networkDeviceIps.Ipv4 {
			entry.Ipv6 = append(entry.Ipv6, types.StringValue(ip))
		}
		ipData[network] = entry
	}
	state.IpData = ipData

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *qwiltOriginAllowListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = cdnclient.NewDeviceIpsClient(client)
}
