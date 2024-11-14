package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SiteModel maps Site schema data.
type SiteModel struct {
	SiteId                         types.String `tfsdk:"site_id"`
	OwnerOrgId                     types.String `tfsdk:"owner_org_id"`
	CreationTimeMilli              types.Int64  `tfsdk:"creation_time_milli"`
	LastUpdateTimeMilli            types.Int64  `tfsdk:"last_update_time_milli"`
	CreatedUser                    types.String `tfsdk:"created_user"`
	LastUpdatedUser                types.String `tfsdk:"last_updated_user"`
	SiteDnsCnameDelegationTarget   types.String `tfsdk:"site_dns_cname_delegation_target"`
	SiteName                       types.String `tfsdk:"site_name"`
	ApiVersion                     types.String `tfsdk:"api_version"`
	ServiceType                    types.String `tfsdk:"service_type"`
	RoutingMethod                  types.String `tfsdk:"routing_method"`
	ShouldProvisionToThirdPartyCdn types.Bool   `tfsdk:"should_provision_to_third_party_cdn"`
	ServiceId                      types.String `tfsdk:"service_id"`
	IsSelfServiceBlocked           types.Bool   `tfsdk:"is_self_service_blocked"`
	IsDeleted                      types.Bool   `tfsdk:"is_deleted"`
}

// PubOpModel maps publishing operation schema data.
type PubOpModel struct {
	PublishId            types.String   `tfsdk:"publish_id"`
	CreationTimeMilli    types.Int64    `tfsdk:"creation_time_milli"`
	OwnerOrgId           types.String   `tfsdk:"owner_org_id"`
	LastUpdateTimeMilli  types.Int64    `tfsdk:"last_update_time_milli"`
	RevisionId           types.String   `tfsdk:"revision_id"`
	Target               types.String   `tfsdk:"target"`
	Username             types.String   `tfsdk:"username"`
	PublishState         types.String   `tfsdk:"publish_state"`
	PublishStatus        types.String   `tfsdk:"publish_status"`
	OperationType        types.String   `tfsdk:"operation_type"`
	StatusLine           []types.String `tfsdk:"status_line"`
	IsActive             types.Bool     `tfsdk:"is_active"`
	ValidatorsErrDetails types.String   `tfsdk:"validators_err_details"`
}

// SiteConfigModel maps site configuration schema data.
type SiteConfigModel struct {
	SiteId              types.String `tfsdk:"site_id"`
	RevisionId          types.String `tfsdk:"revision_id"`
	RevisionNum         types.Int64  `tfsdk:"revision_num"`
	OwnerOrgId          types.String `tfsdk:"owner_org_id"`
	CreationTimeMilli   types.Int64  `tfsdk:"creation_time_milli"`
	LastUpdateTimeMilli types.Int64  `tfsdk:"last_update_time_milli"`
	CreatedUser         types.String `tfsdk:"created_user"`
	HostIndex           types.String `tfsdk:"host_index"`
	ChangeDescription   types.String `tfsdk:"change_description"`
}

// certificateModel maps Certificate schema data.
type CertificateDataModel struct {
	CertId           types.Int64  `tfsdk:"cert_id"`
	Certificate      types.String `tfsdk:"certificate"`
	CertificateChain types.String `tfsdk:"certificate_chain"`
	//Email            types.String `tfsdk:"email"`
	Description types.String `tfsdk:"description"`
	PkHash      types.String `tfsdk:"pk_hash"`
	Tenant      types.String `tfsdk:"tenant"`
	Domain      types.String `tfsdk:"domain"`
	Status      types.String `tfsdk:"status"`
	Type        types.String `tfsdk:"type"`
}

// IP ALLOW List model
type NetworkIpDataModel struct {
	Ipv4 []types.String `tfsdk:"ipv4"`
	Ipv6 []types.String `tfsdk:"ipv6"`
}
