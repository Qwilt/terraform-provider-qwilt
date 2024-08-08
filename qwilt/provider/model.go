package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// QwiltProviderModel describes the provider data model.
type QwiltProviderModel struct {
	// CDN
	EnvType   types.String `tfsdk:"env_type"`
	Username  types.String `tfsdk:"username"`
	Password  types.String `tfsdk:"password"`
	XApiToken types.String `tfsdk:"xapi_token"`
}
