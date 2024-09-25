package model

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// SiteConfiguration maps site configuration schema data.
type SiteConfiguration struct {
	Id                  types.String    `tfsdk:"id"`
	SiteId              types.String    `tfsdk:"site_id"`
	RevisionId          types.String    `tfsdk:"revision_id"`
	RevisionNum         types.Int64     `tfsdk:"revision_num"`
	OwnerOrgId          types.String    `tfsdk:"owner_org_id"`
	HostIndex           HostIndexString `tfsdk:"host_index"`
	ChangeDescription   types.String    `tfsdk:"change_description"`
	LastUpdateTimeMilli types.Int64     `tfsdk:"last_update_time_milli"`
}

// SiteConfigBuilder is a builder for SiteConfiguration
type SiteConfigBuilder struct {
	cfg SiteConfiguration
	ctx context.Context
}

// Custom type definition definition of HostIndexString and corresponding HostIndexType
var _ basetypes.StringTypable = (*HostIndexType)(nil)
var _ basetypes.StringValuable = (*HostIndexString)(nil)

type HostIndexType struct {
	basetypes.StringType
}

type HostIndexString struct {
	basetypes.StringValue
}

// JsonBytesEqual compares the JSON in two byte slices for equality
func JsonBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}

// HostIndexType custom type methods
func (t HostIndexType) Equal(o attr.Type) bool {
	other, ok := o.(HostIndexType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t HostIndexType) String() string {
	return "HostIndexType"
}

func (t HostIndexType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := HostIndexString{
		StringValue: in,
	}

	return value, nil
}

func (t HostIndexType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

func (t HostIndexType) ValueType(ctx context.Context) attr.Value {
	return HostIndexString{}
}

// HostIndexString custom value methods
func (v HostIndexString) Equal(o attr.Value) bool {
	other, ok := o.(HostIndexString)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v HostIndexString) Type(ctx context.Context) attr.Type {
	return HostIndexType{}
}

// StringSemanticEquals returns true if the given string value is semantically equal to the current string value
func (v HostIndexString) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(HostIndexString)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)
		return false, diags
	}

	// Compare two JSON strings to determine if they are semantically equal
	hostIndexEqual, err := JsonBytesEqual([]byte(newValue.ValueString()), []byte(v.ValueString()))
	if err != nil {
		diags.AddError(
			"Error Unmarshaling HostIndex for Comparison",
			"Could not compare HostIndex JSON: "+err.Error(),
		)
		return false, diags
	}

	return hostIndexEqual, diags
}

// SiteConfigBuilder methods
func NewSiteConfigBuilder() *SiteConfigBuilder {
	b := SiteConfigBuilder{}
	return &b
}

func (b *SiteConfigBuilder) WithCtx(ctx context.Context) *SiteConfigBuilder {
	b.ctx = ctx
	return b
}
func (b *SiteConfigBuilder) LastUpdateTimeMilli(value int) *SiteConfigBuilder {
	b.cfg.LastUpdateTimeMilli = types.Int64Value(int64(value))
	return b
}
func (b *SiteConfigBuilder) WithSiteId(siteId string) *SiteConfigBuilder {
	b.cfg.SiteId = types.StringValue(siteId)
	return b
}
func (b *SiteConfigBuilder) WithOwnerOrgId(ownerOrgId string) *SiteConfigBuilder {
	b.cfg.OwnerOrgId = types.StringValue(ownerOrgId)
	return b
}
func (b *SiteConfigBuilder) WithRevisionId(revision string) *SiteConfigBuilder {
	b.cfg.RevisionId = types.StringValue(revision)
	return b
}
func (b *SiteConfigBuilder) WithRevisionNum(revision int) *SiteConfigBuilder {
	b.cfg.RevisionNum = types.Int64Value(int64(revision))
	return b
}
func (b *SiteConfigBuilder) WithHostIndex(hostIndex json.RawMessage, indent bool) *SiteConfigBuilder {

	if indent {
		// Format the HostIndex JSON string from the API.
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

		b.cfg.HostIndex = HostIndexString{types.StringValue(hostIndexIndented)}
	} else {
		b.cfg.HostIndex = HostIndexString{types.StringValue(string(hostIndex))}
	}
	return b
}
func (b *SiteConfigBuilder) WithChangeDescription(desc string) *SiteConfigBuilder {
	b.cfg.ChangeDescription = types.StringValue(desc)
	return b
}

func (b *SiteConfigBuilder) Build() SiteConfiguration {
	id := b.cfg.SiteId.ValueString() + ":" + b.cfg.RevisionId.ValueString()
	b.cfg.Id = types.StringValue(id)
	return b.cfg
}
