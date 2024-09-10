package model

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// SiteConfiguration maps site configuration schema data.
type SiteConfiguration struct {
	Id                  types.String `tfsdk:"id"`
	SiteId              types.String `tfsdk:"site_id"`
	RevisionId          types.String `tfsdk:"revision_id"`
	RevisionNum         types.Int64  `tfsdk:"revision_num"`
	OwnerOrgId          types.String `tfsdk:"owner_org_id"`
	HostIndex           types.String `tfsdk:"host_index"`
	ChangeDescription   types.String `tfsdk:"change_description"`
	LastUpdateTimeMilli types.Int64  `tfsdk:"last_update_time_milli"`
}

type SiteConfigBuilder struct {
	cfg SiteConfiguration
	ctx context.Context
}

func (b SiteConfigBuilder) WithCtx(ctx context.Context) SiteConfigBuilder {
	b.ctx = ctx
	return b
}
func (b SiteConfigBuilder) LastUpdateTimeMilli(value int) SiteConfigBuilder {
	b.cfg.LastUpdateTimeMilli = types.Int64Value(int64(value))
	return b
}
func (b SiteConfigBuilder) WithSiteId(siteId string) SiteConfigBuilder {
	b.cfg.SiteId = types.StringValue(siteId)
	return b
}
func (b SiteConfigBuilder) WithOwnerOrgId(ownerOrgId string) SiteConfigBuilder {
	b.cfg.OwnerOrgId = types.StringValue(ownerOrgId)
	return b
}
func (b SiteConfigBuilder) WithRevisionId(revision string) SiteConfigBuilder {
	b.cfg.RevisionId = types.StringValue(revision)
	return b
}
func (b SiteConfigBuilder) WithRevisionNum(revision int) SiteConfigBuilder {
	b.cfg.RevisionNum = types.Int64Value(int64(revision))
	return b
}
func (b SiteConfigBuilder) WithHostIndex(hostIndex json.RawMessage, indent bool) SiteConfigBuilder {

	if indent {
		//Format the HostIndex JSON string from the API.
		// Consistent formatting is important so that Terraform does not continue
		// trying to update the HostIndex attribute.
		// An additional newline character is added to match the state input.
		var hostIndexIndented string
		hostIndexIndentedJson, err := json.MarshalIndent(hostIndex, "", "\t")
		if err != nil {
			tflog.Info(b.ctx, "Failed to format HostIndex string")
		} else {
			hostIndexIndented = string(hostIndexIndentedJson) + "\n"
		}

		b.cfg.HostIndex = types.StringValue(hostIndexIndented)
	} else {
		b.cfg.HostIndex = types.StringValue(string(hostIndex))
	}
	return b
}
func (b SiteConfigBuilder) WithChangeDescription(desc string) SiteConfigBuilder {
	b.cfg.ChangeDescription = types.StringValue(desc)
	return b
}

func (b SiteConfigBuilder) Build() SiteConfiguration {
	id := b.cfg.SiteId.ValueString() + ":" + b.cfg.RevisionId.ValueString()
	b.cfg.Id = types.StringValue(id)
	return b.cfg
}
