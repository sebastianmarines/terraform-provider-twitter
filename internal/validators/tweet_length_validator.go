package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type tweetLengthValidator struct {
	Max int
	Min int
}

func (v tweetLengthValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Tweet length must be between %d and %d characters.", v.Min, v.Max)
}

func (v tweetLengthValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Tweet length must be between %d and %d characters.", v.Min, v.Max)
}

func (v tweetLengthValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if str.Unknown || str.Null {
		return
	}

	strLen := len(str.Value)

	if strLen < v.Min || strLen > v.Max {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Tweet Length",
			fmt.Sprintf("Tweet length must be between %d and %d characters, got: %d characters.", v.Min, v.Max, strLen),
		)

		return
	}
}

func TweetLength() tweetLengthValidator {
	maxLength := 280
	minLength := 1

	return tweetLengthValidator{
		Max: maxLength,
		Min: minLength,
	}
}
