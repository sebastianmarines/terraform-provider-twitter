package utils

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type Response struct {
	State tfsdk.State
	Diags diag.Diagnostics
}

func CheckProviderConfiguration(d *diag.Diagnostics, configured bool) error {
	if !configured {
		d.AddError(
			"Provider not configured",
			"The provider is not configured. Please configure the provider before using it.",
		)
		return errors.New("Provider not configured")
	}

	return nil
}
