package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type blankNameValidator struct {
	Max int
	Min int
}

func (v blankNameValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Tweet length must be between %d and %d characters.", v.Min, v.Max)
}

func (v blankNameValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Tweet length must be between %d and %d characters.", v.Min, v.Max)
}

func (v blankNameValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if str.Unknown {
		return
	}

	if str.Null || len(str.Value) == 0 {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid account name.",
			fmt.Sprintf("Name can't be blank."),
		)

		return
	}
}

func BlankName() blankNameValidator {
	return blankNameValidator{}
}
