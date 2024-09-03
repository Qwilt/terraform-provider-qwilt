package cdn

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MutualExclusiveValidator implements the validator.Int64 interface
type MutualExclusiveValidator struct {
	OtherKey path.Path
}

func (v MutualExclusiveValidator) Description(ctx context.Context) string {
	return "Mutual-Exclusion validator"
}

func (v MutualExclusiveValidator) MarkdownDescription(ctx context.Context) string {
	return "Mutual-Exclusion validator"
}

// ValidateInt64 performs the validation for mutual exclusivity
func (v MutualExclusiveValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	// Skip validation if the current value is null or unknown
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Retrieve the value of the other attribute
	otherVal := types.Int64{}
	diags := req.Config.GetAttribute(ctx, v.OtherKey, &otherVal)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if both attributes are set
	if !otherVal.IsNull() && !otherVal.IsUnknown() {
		resp.Diagnostics.AddError(
			"Mutual Exclusivity Error",
			fmt.Sprintf("Only one of '%s' or '%s' can be set. Please unset one of them.", req.Path, v.OtherKey),
		)
	}
}

// NewMutualExclusiveValidator creates a new MutualExclusiveValidator
func NewMutualExclusiveValidator(otherKey path.Path) MutualExclusiveValidator {
	return MutualExclusiveValidator{
		OtherKey: otherKey,
	}
}
