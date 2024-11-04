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
	_ datasource.DataSource              = &qwiltIpAllowListDataSource{}
	_ datasource.DataSourceWithConfigure = &qwiltIpAllowListDataSource{}
)

// NewIpAllowListDataSource is a helper function to simplify the provider implementation.
func NewIpAllowListDataSource() datasource.DataSource {
	return &qwiltIpAllowListDataSource{}
}

// qwiltIpAllowListDataSource is the data source implementation.
type qwiltIpAllowListDataSource struct {
	client *cdnclient.DeviceIpsClient
}

// qwiltIpAllowListDataSourceModel maps the data source schema data.
type qwiltIpAllowListDataSourceModel struct {
	Md5              types.String                  `tfsdk:"md5"`
	CreateTimeMillis types.Int64                   `tfsdk:"createTimeMillis"`
	Networks         []cdnmodel.NetworkIpDataModel `tfsdk:"networks"`
}

// Metadata returns the data source type name.
func (d *qwiltIpAllowListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_ip_allow_list"
}

// Schema defines the schema for the data source.
func (d *qwiltIpAllowListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the device ip's to be added to origin allow iist.",
		Attributes: map[string]schema.Attribute{
			"md5": schema.StringAttribute{
				Description: "MD5 hash value",
				Computed:    true,
			},
			"create_time_millis": schema.Int64Attribute{
				Description: "Creation time in milliseconds",
				Computed:    true,
			},
			//"networks": schema.SingleNestedAttribute{
			//	Description: "IP data containing IPv4 and IPv6 addresses",
			//	Computed:    true,
			//	Attributes: map[string]schema.Attribute{
			//		"ipv4": schema.ListAttribute{
			//			ElementType: types.StringType,
			//			Description: "List of IPv4 addresses",
			//			Computed:    true,
			//		},
			//		"ipv6": schema.ListAttribute{
			//			ElementType: types.StringType,
			//			Description: "List of IPv6 addresses",
			//			Computed:    true,
			//		},
			//	},
			//},
			"networks": schema.ListNestedAttribute{
				Description: "List of networks and their device-ips.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The network name",
							Computed:    true,
						},
						"ipv4": schema.ListAttribute{
							ElementType: types.StringType,
							Description: "List of IPv4 addresses",
							Computed:    true,
						},
						"ipv6": schema.ListAttribute{
							ElementType: types.StringType,
							Description: "List of IPv6 addresses",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *qwiltIpAllowListDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state qwiltIpAllowListDataSourceModel

	deviceIpsResp, err := d.client.GetOriginAllowList()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Qwilt IpAllowList",
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
	// Get networks ips
	deviceIpsState := cdnmodel.DeviceIpsDataModel{
		Md5:              types.StringValue(deviceIpsResp.Md5),
		CreateTimeMillis: types.Int64Value(int64(deviceIpsResp.CreateTimeMillis)),
	}
	networks := make([]cdnmodel.NetworkIpDataModel, 0)

	for network, networkDeviceIps := range deviceIpsResp.IpData {
		entry := cdnmodel.NetworkIpDataModel{
			Name: types.StringValue(network),
		}
		for _, ip := range networkDeviceIps.Ipv4 {
			entry.Ipv4 = append(entry.Ipv4, types.StringValue(ip))
		}
		for _, ip := range networkDeviceIps.Ipv4 {
			entry.Ipv6 = append(entry.Ipv6, types.StringValue(ip))
		}
		networks = append(networks, entry)
	}
	deviceIpsState.Networks = networks
	//state.DeviceIpsDataModel = deviceIpsState

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *qwiltIpAllowListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
