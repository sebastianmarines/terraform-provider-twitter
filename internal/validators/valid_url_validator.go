package validators

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type validURLValidator struct {
	Max int
	Min int
}

func (v validURLValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("The URL must be a valid URL.")
}

func (v validURLValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("The URL must be a valid URL.")
}

func (v validURLValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var _url types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &_url)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if _url.Unknown {
		return
	}

	if _url.Null {
		return
	}

	u, err := url.Parse(_url.Value)

	if err != nil || u.Scheme == "" || u.Host == "" {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid URL.",
			fmt.Sprintf("The URL is not a valid URL"),
		)

		return
	}
}

func ValidURL() validURLValidator {
	return validURLValidator{}
}
