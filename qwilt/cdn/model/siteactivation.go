package model

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SiteActivation struct {
	Id                      types.String `tfsdk:"id"`
	SiteId                  types.String `tfsdk:"site_id"`
	RevisionId              types.String `tfsdk:"revision_id"`
	CertificateId           types.Int64  `tfsdk:"certificate_id"`
	CsrId                   types.Int64  `tfsdk:"csr_id"`
	PublishId               types.String `tfsdk:"publish_id"`
	CreationTimeMilli       types.Int64  `tfsdk:"creation_time_milli"`
	OwnerOrgId              types.String `tfsdk:"owner_org_id"`
	LastUpdateTimeMilli     types.Int64  `tfsdk:"last_update_time_milli"`
	Target                  types.String `tfsdk:"target"`
	Username                types.String `tfsdk:"username"`
	PublishState            types.String `tfsdk:"publish_state"`
	PublishStatus           types.String `tfsdk:"publish_status"`
	PublishAcceptanceStatus types.String `tfsdk:"publish_acceptance_status"`
	OperationType           types.String `tfsdk:"operation_type"`
	//StatusLine          []types.String `tfsdk:"status_line"`
	IsActive           types.Bool   `tfsdk:"is_active"`
	ValidateErrDetails types.String `tfsdk:"validators_err_details"`
}

type SiteActivationBuilder struct {
	activation SiteActivation
	ctx        context.Context
}

func (b SiteActivationBuilder) Ctx(ctx context.Context) SiteActivationBuilder {
	b.ctx = ctx
	return b
}
func (b SiteActivationBuilder) PublishId(value string) SiteActivationBuilder {
	b.activation.PublishId = types.StringValue(value)
	return b
}
func (b SiteActivationBuilder) RevisionId(value string) SiteActivationBuilder {
	b.activation.RevisionId = types.StringValue(value)
	return b
}
func (b SiteActivationBuilder) SiteId(value string) SiteActivationBuilder {
	b.activation.SiteId = types.StringValue(value)
	return b
}
func (b SiteActivationBuilder) CertificateId(value int64) SiteActivationBuilder {
	if value != 0 {
		b.activation.CertificateId = types.Int64Value(value)
	} else {
		b.activation.CertificateId = types.Int64Null()
	}
	return b
}
func (b SiteActivationBuilder) CsrId(value int64) SiteActivationBuilder {
	if value != 0 {
		b.activation.CsrId = types.Int64Value(value)
	} else {
		b.activation.CsrId = types.Int64Null()
	}
	return b
}
func (b SiteActivationBuilder) PublishState(value string) SiteActivationBuilder {
	b.activation.PublishState = types.StringValue(value)
	return b
}
func (b SiteActivationBuilder) OperationType(value string) SiteActivationBuilder {
	b.activation.OperationType = types.StringValue(value)
	return b
}
func (b SiteActivationBuilder) LastUpdateTimeMilli(value int) SiteActivationBuilder {
	b.activation.LastUpdateTimeMilli = types.Int64Value(int64(value))
	return b
}
func (b SiteActivationBuilder) Username(value string) SiteActivationBuilder {
	b.activation.Username = types.StringValue(value)
	return b
}
func (b SiteActivationBuilder) PublishStatus(value string) SiteActivationBuilder {
	b.activation.PublishStatus = types.StringValue(value)
	return b
}
func (b SiteActivationBuilder) AcceptanceStatus(value string) SiteActivationBuilder {
	b.activation.PublishAcceptanceStatus = types.StringValue(value)
	return b
}
func (b SiteActivationBuilder) IsActive(value bool) SiteActivationBuilder {
	b.activation.IsActive = types.BoolValue(value)
	return b
}

//	func (b SiteActivationBuilder) StatusLine(values []string) SiteActivationBuilder {
//		typesStrSlice := make([]types.String, len(values))
//		for i, s := range values {
//			typesStrSlice[i] = types.StringValue(s)
//		}
//		b.activation.StatusLine = typesStrSlice
//		return b
//	}
func (b SiteActivationBuilder) Target(value string) SiteActivationBuilder {
	b.activation.Target = types.StringValue(value)
	return b
}
func (b SiteActivationBuilder) ValidateErrDetails(value json.RawMessage) SiteActivationBuilder {
	b.activation.ValidateErrDetails = types.StringValue(string(value))
	return b
}
func (b SiteActivationBuilder) Build() SiteActivation {
	id := b.activation.SiteId.ValueString() + ":" + b.activation.PublishId.ValueString()
	b.activation.Id = types.StringValue(id)
	return b.activation
}
