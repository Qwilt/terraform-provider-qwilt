package custome_modifiers

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// CustomPlanModifierSuppressDiff is a custom PlanModifier that suppresses diffs for private_key
type CustomPlanModifierSuppressDiff struct{}

func (m CustomPlanModifierSuppressDiff) Description(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will ignore this diff."
}

func (m CustomPlanModifierSuppressDiff) MarkdownDescription(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will ignore this diff."
}

// Modify the plan based on the current state and plan inputs
func (m CustomPlanModifierSuppressDiff) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the private_key is set in the plan but not in the state, suppress the diff
	if !req.ConfigValue.IsNull() && req.StateValue.IsNull() {
		resp.PlanValue = req.ConfigValue // Use the provided value from the config (input)
	}
}
