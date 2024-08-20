package model

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Site struct {
	LastUpdated                  types.String `tfsdk:"last_updated"`
	Id                           types.String `tfsdk:"id"`
	SiteId                       types.String `tfsdk:"site_id"`
	OwnerOrgId                   types.String `tfsdk:"owner_org_id"`
	SiteDnsCnameDelegationTarget types.String `tfsdk:"site_dns_cname_delegation_target"`
	SiteName                     types.String `tfsdk:"site_name"`
	RoutingMethod                types.String `tfsdk:"routing_method"`
	LastUpdateTimeMilli          types.Int64  `tfsdk:"last_update_time_milli"`
}

type SiteBuilder struct {
	site Site
	ctx  context.Context
}

func (b SiteBuilder) WithCtx(ctx context.Context) SiteBuilder {
	b.ctx = ctx
	return b
}
func (b SiteBuilder) LastUpdateTimeMilli(value int) SiteBuilder {
	b.site.LastUpdateTimeMilli = types.Int64Value(int64(value))
	return b
}
func (b SiteBuilder) SiteId(value string) SiteBuilder {
	b.site.SiteId = types.StringValue(value)
	b.site.Id = b.site.SiteId
	return b
}
func (b SiteBuilder) OwnerOrgId(value string) SiteBuilder {
	b.site.OwnerOrgId = types.StringValue(value)
	return b
}
func (b SiteBuilder) SiteName(value string) SiteBuilder {
	b.site.SiteName = types.StringValue(value)
	return b
}
func (b SiteBuilder) RoutingMethod(value string) SiteBuilder {
	b.site.RoutingMethod = types.StringValue(value)
	return b
}
func (b SiteBuilder) SiteDnsCnameDelegationTarget(value string) SiteBuilder {
	b.site.SiteDnsCnameDelegationTarget = types.StringValue(value)
	return b
}
func (b SiteBuilder) Build() Site {
	return b.site
}
