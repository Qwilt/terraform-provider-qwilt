// Package qwiltcdn
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
	cdnmodel "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/model"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	cdnclient "github.com/Qwilt/terraform-provider-qwilt/qwilt/cdn/client"
)

var (
	_ datasource.DataSource              = &qwiltSitesDataSource{}
	_ datasource.DataSourceWithConfigure = &qwiltSitesDataSource{}
)

// NewSitesDataSource is a helper function to simplify the provider implementation.
func NewSitesDataSource() datasource.DataSource {
	return &qwiltSitesDataSource{}
}

// qwiltSitesDataSource is the data source implementation.
type qwiltSitesDataSource struct {
	client *cdnclient.SiteClientFacade
}

// Metadata returns the data source type name.
func (d *qwiltSitesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_sites"
}

// Schema defines the schema for the data source.
func (d *qwiltSitesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches Qwilt Sites configuration.",
		Attributes: map[string]schema.Attribute{
			"site": schema.ListNestedAttribute{
				Description: "List of site configurations.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"site_id": schema.StringAttribute{
							Description: "The unique identifier of the site. The siteID will be needed when you add the site configuration and when you publish the site.",
							Computed:    true,
						},
						"owner_org_id": schema.StringAttribute{
							Description: "The name of the organization that owns the site.",
							Computed:    true,
						},
						"creation_time_milli": schema.Int64Attribute{
							Description: "When the site was created, in epoch time.",
							Computed:    true,
						},
						"last_update_time_milli": schema.Int64Attribute{
							Description: "When the site was last updated, in epoch time.",
							Computed:    true,
						},
						"created_user": schema.StringAttribute{
							Description: "The user who created the site.",
							Computed:    true,
						},
						"last_updated_user": schema.StringAttribute{
							Description: "The user who last updated the site.",
							Computed:    true,
						},
						"site_dns_cname_delegation_target": schema.StringAttribute{
							Description: "The CNAME you'll use direct traffic from your website to the CDN.",
							Computed:    true,
						},
						"site_name": schema.StringAttribute{
							Description: "The user-defined site name.",
							Computed:    true,
						},
						"api_version": schema.StringAttribute{
							Description: "The Media Delivery Configuration API version.",
							Computed:    true,
						},
						"service_type": schema.StringAttribute{
							Description: "The value will always be MEDIA_DELIVERY.",
							Computed:    true,
						},
						"routing_method": schema.StringAttribute{
							Description: "The routing method used for the site. DNS is the default routing mechanism.",
							Computed:    true,
						},
						"should_provision_to_third_party_cdn": schema.BoolAttribute{
							Description: "Indicates if the site should be provisioned to third-party CDNs.",
							Computed:    true,
						},
						"service_id": schema.StringAttribute{
							Description: "A system-generated unique identifier for the site configuration.",
							Computed:    true,
						},
						"is_self_service_blocked": schema.BoolAttribute{
							Description: "Indicates if site updates are currently allowed.",
							Computed:    true,
						},
						"is_deleted": schema.BoolAttribute{
							Description: "Indicates if site was marked for deletion.",
							Computed:    true,
						},
					},
				},
			},
			"revision": schema.ListNestedAttribute{
				Description: "Site configurations associated with the site.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"site_id": schema.StringAttribute{
							Description: "The unique identifier of the site.",
							Computed:    true,
						},
						"revision_id": schema.StringAttribute{
							Description: "The unique identifier of the configuration version.",
							Computed:    true,
						},
						"revision_num": schema.Int64Attribute{
							Description: "The unique revision number of the configuration version.",
							Computed:    true,
						},
						"owner_org_id": schema.StringAttribute{
							Description: "The name of the organization that owns the site.",
							Computed:    true,
						},
						"creation_time_milli": schema.Int64Attribute{
							Description: "The time when the configuration version was added, in epoch time.",
							Computed:    true,
						},
						"last_update_time_milli": schema.Int64Attribute{
							Description: "The time when the configuration version was added, in epoch time. (This will be the same as the creationTimeMilli value.)",
							Computed:    true,
						},
						"created_user": schema.StringAttribute{
							Description: "The user who created the site.",
							Computed:    true,
						},
						"host_index": schema.StringAttribute{
							Description: "The SVTA metadata objects that define the delivery service	   configuration, in application/json format.",
							Computed:    true,
						},
						"change_description": schema.StringAttribute{
							Description: "Comments added by the user to the configuration JSON payload.",
							Computed:    true,
						},
					},
				},
			},
			"publish_op": schema.ListNestedAttribute{
				Description: "List of publishing operations associated with site.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"publish_id": schema.StringAttribute{
							Description: "Unique identifier of the publishing operation.",
							Computed:    true,
						},
						"creation_time_milli": schema.Int64Attribute{
							Description: "The time when the publish operation was created, in epoch time.",
							Computed:    true,
						},
						"owner_org_id": schema.StringAttribute{
							Description: "The organization that owns the site.",
							Computed:    true,
						},
						"last_update_time_milli": schema.Int64Attribute{
							Description: "When the publishing operation was last updated, in epoch time.",
							Computed:    true,
						},
						"revision_id": schema.StringAttribute{
							Description: "Unique identifier of the configuration version that was published or unpublished.",
							Computed:    true,
						},
						"target": schema.StringAttribute{
							Description: "The value will 'ga' or 'staging'.",
							Computed:    true,
						},
						"username": schema.StringAttribute{
							Description: "Username that initiated the publishing operation.",
							Computed:    true,
						},
						"publish_state": schema.StringAttribute{
							Description: "For internal use. Use the 'publishStatus' values instead.",
							Computed:    true,
						},
						"publish_status": schema.StringAttribute{
							Description: "The publishing operation status.",
							Computed:    true,
						},
						"operation_type": schema.StringAttribute{
							Description: "The operation type (Publish, Unpublish) that was initiated with the request. An Unpublish operation removes a delivery service from the CDN.",
							Computed:    true,
						},
						"status_line": schema.ListAttribute{
							ElementType: types.StringType,
							Description: "Additional information related to the publish status.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Indicates if the configuration is active or inactive.",
							Computed:    true,
						},
						"validators_err_details": schema.StringAttribute{
							Description: "Details about errors generated during validation.",
							Computed:    true,
						},
					},
				},
			},
			"filter": schema.SingleNestedAttribute{
				Description: "Data source filter attributes.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"site_id": schema.StringAttribute{
						Description: "Filter sites based on a specific site ID.",
						Optional:    true,
					},
					"revision_id": schema.StringAttribute{
						Description: "Filter configurations based on a specific revision ID.",
						Optional:    true,
					},
					"publish_id": schema.StringAttribute{
						Description: "filter publishing operations based on a specific publish ID.",
						Optional:    true,
					},
					"truncate_host_index": schema.BoolAttribute{
						Description: "By default, the configuration details are included in the response, and you can exclude them by setting this to false.",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *qwiltSitesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cdnmodel.QwiltSitesDataSourceModel

	sites, err := d.client.GetSites(true, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Qwilt Sites",
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
	var filter cdnmodel.QwiltSitesFilterModel
	var siteIdFilter string
	var revisionIdFilter string
	var publishIdFilter string
	var truncateHostIndex bool

	if state.Filter.IsNull() {
		siteIdFilter = "all"
		revisionIdFilter = "all"
		publishIdFilter = "all"
		truncateHostIndex = false
	} else {
		diags := state.Filter.As(ctx, &filter, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		siteIdFilter = filter.SiteId.ValueString()
		revisionIdFilter = filter.RevisionId.ValueString()
		publishIdFilter = filter.PublishId.ValueString()
		truncateHostIndex = filter.TruncateHostIndex.ValueBool()
		if siteIdFilter == "" {
			siteIdFilter = "all"
		}
		if revisionIdFilter == "" {
			revisionIdFilter = "all"
		}
		if publishIdFilter == "" {
			publishIdFilter = "all"
		}
	}

	// Map response body to model
	// Get site(s)
	for _, site := range sites {
		if siteIdFilter != "all" && site.SiteId != siteIdFilter {
			continue
		}
		siteState := cdnmodel.SiteModel{
			SiteId:                         types.StringValue(site.SiteId),
			OwnerOrgId:                     types.StringValue(site.OwnerOrgId),
			CreationTimeMilli:              types.Int64Value(int64(site.CreationTimeMilli)),
			LastUpdateTimeMilli:            types.Int64Value(int64(site.LastUpdateTimeMilli)),
			CreatedUser:                    types.StringValue(site.CreatedUser),
			LastUpdatedUser:                types.StringValue(site.LastUpdatedUser),
			SiteDnsCnameDelegationTarget:   types.StringValue(site.SiteDnsCnameDelegationTarget),
			SiteName:                       types.StringValue(site.SiteName),
			ApiVersion:                     types.StringValue(site.ApiVersion),
			ServiceType:                    types.StringValue(site.ServiceType),
			RoutingMethod:                  types.StringValue(site.RoutingMethod),
			ShouldProvisionToThirdPartyCdn: types.BoolValue(site.ShouldProvisionToThirdPartyCdn),
			ServiceId:                      types.StringValue(site.ServiceId),
			IsSelfServiceBlocked:           types.BoolValue(site.IsSelfServiceBlocked),
			IsDeleted:                      types.BoolValue(site.IsDeleted),
		}

		state.Site = append(state.Site, siteState)
	}

	// Get revision(s)
	if siteIdFilter != "all" {
		siteConfigs, err := d.client.GetSiteConfigs(siteIdFilter, truncateHostIndex)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to Read Qwilt Site Config for SiteId %s", siteIdFilter),
				err.Error(),
			)
			return
		}
		for _, siteConfig := range siteConfigs {
			if revisionIdFilter != "all" && siteConfig.RevisionId != revisionIdFilter {
				continue
			}

			siteConfigState := cdnmodel.SiteConfigModel{
				SiteId:              types.StringValue(siteConfig.SiteId),
				RevisionId:          types.StringValue(siteConfig.RevisionId),
				RevisionNum:         types.Int64Value(int64(siteConfig.RevisionNum)),
				OwnerOrgId:          types.StringValue(siteConfig.OwnerOrgId),
				CreationTimeMilli:   types.Int64Value(int64(siteConfig.CreationTimeMilli)),
				LastUpdateTimeMilli: types.Int64Value(int64(siteConfig.LastUpdateTimeMilli)),
				CreatedUser:         types.StringValue(siteConfig.CreatedUser),
				HostIndex:           types.StringValue(string(siteConfig.HostIndex)),
				ChangeDescription:   types.StringValue(siteConfig.ChangeDescription),
			}

			state.Revision = append(state.Revision, siteConfigState)
		}
	}

	// Get publishing operation(s)
	if siteIdFilter != "all" {
		pubOps, err := d.client.GetPubOps(siteIdFilter, false, "")
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to Read Qwilt Publishing Operations for SiteId %s", siteIdFilter),
				err.Error(),
			)
			return
		}
		for _, pubOp := range pubOps {
			if publishIdFilter != "all" && pubOp.PublishId != publishIdFilter {
				continue
			}
			pubOpState := cdnmodel.PubOpModel{
				PublishId:            types.StringValue(pubOp.PublishId),
				CreationTimeMilli:    types.Int64Value(int64(pubOp.CreationTimeMilli)),
				OwnerOrgId:           types.StringValue(pubOp.OwnerOrgId),
				LastUpdateTimeMilli:  types.Int64Value(int64(pubOp.LastUpdateTimeMilli)),
				RevisionId:           types.StringValue(pubOp.RevisionId),
				Target:               types.StringValue(pubOp.Target),
				Username:             types.StringValue(pubOp.Username),
				PublishState:         types.StringValue(pubOp.PublishState),
				PublishStatus:        types.StringValue(pubOp.PublishStatus),
				OperationType:        types.StringValue(pubOp.OperationType),
				IsActive:             types.BoolValue(pubOp.IsActive),
				ValidatorsErrDetails: types.StringValue(string(pubOp.ValidatorsErrDetails)),
			}

			for _, statusLine := range pubOp.StatusLine {
				pubOpState.StatusLine = append(pubOpState.StatusLine, types.StringValue(statusLine))
			}

			state.PubOp = append(state.PubOp, pubOpState)
		}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *qwiltSitesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = cdnclient.NewSiteFacadeClient(api.SITES_HOSTNAME, client)
}
