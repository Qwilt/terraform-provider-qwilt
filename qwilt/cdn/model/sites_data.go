package model

import "github.com/hashicorp/terraform-plugin-framework/types"

// QwiltSitesDataSourceModel maps the data source schema data.
type QwiltSitesDataSourceModel struct {
	Site     []SiteModel       `tfsdk:"site"`
	Revision []SiteConfigModel `tfsdk:"revision"`
	PubOp    []PubOpModel      `tfsdk:"publish_op"`
	Filter   types.Object      `tfsdk:"filter"`
}

// QwiltSitesFilterModel
type QwiltSitesFilterModel struct {
	SiteId            types.String `tfsdk:"site_id"`
	RevisionId        types.String `tfsdk:"revision_id"`
	PublishId         types.String `tfsdk:"publish_id"`
	TruncateHostIndex types.Bool   `tfsdk:"truncate_host_index"`
}
