package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.Provider = &provider{}

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type provider struct {
	// client can contain the upstream provider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.
	//
	client twitter.Client

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	ApiKey       types.String `tfsdk:"api_key"`
	ApiSecretKey types.String `tfsdk:"api_secret_key"`
	AccessToken  types.String `tfsdk:"access_token"`
	AccessSecret types.String `tfsdk:"access_token_secret"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var apiKey string
	var apiSecretKey string
	var accessToken string
	var accessTokenSecret string

	if data.ApiKey.Unknown {
		resp.Diagnostics.AddWarning(
			"Missing Twitter API key",
			"The Twitter API key is not configured. The Twitter provider will not be able to function.",
		)
		return
	}

	if data.ApiKey.Null {
		apiKey = os.Getenv("TWITTER_API_KEY")
	} else {
		apiKey = data.ApiKey.Value
	}

	if apiKey == "" {
		resp.Diagnostics.AddWarning(
			"Missing Twitter API key",
			"The Twitter API key is not configured. The Twitter provider will not be able to function.",
		)
		return
	}

	if data.ApiSecretKey.Unknown {
		resp.Diagnostics.AddError(
			"Missing Twitter API secret key",
			"The Twitter API secret key is not configured. The Twitter provider will not be able to function.",
		)
		return
	}

	if data.ApiSecretKey.Null {
		apiSecretKey = os.Getenv("TWITTER_API_SECRET_KEY")
	} else {
		apiSecretKey = data.ApiSecretKey.Value
	}

	if apiSecretKey == "" {
		resp.Diagnostics.AddError(
			"Missing Twitter API secret key",
			"The Twitter API secret key is not configured. The Twitter provider will not be able to function.",
		)
		return
	}

	if data.AccessToken.Unknown {
		resp.Diagnostics.AddWarning(
			"Missing Twitter access token",
			"The Twitter access token is not configured. The Twitter provider will not be able to function.",
		)
		return
	}

	if data.AccessToken.Null {
		accessToken = os.Getenv("TWITTER_ACCESS_TOKEN")
	} else {
		accessToken = data.AccessToken.Value
	}

	if accessToken == "" {
		resp.Diagnostics.AddError(
			"Missing Twitter access token",
			"The Twitter access token is not configured. The Twitter provider will not be able to function.",
		)
		return
	}

	if data.AccessSecret.Unknown {
		resp.Diagnostics.AddWarning(
			"Missing Twitter access secret",
			"The Twitter access secret is not configured. The Twitter provider will not be able to function.",
		)
		return
	}

	if data.AccessSecret.Null {
		accessTokenSecret = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	} else {
		accessTokenSecret = data.AccessSecret.Value
	}

	if accessTokenSecret == "" {
		resp.Diagnostics.AddError(
			"Missing Twitter access secret",
			"The Twitter access secret is not configured. The Twitter provider will not be able to function.",
		)
		return
	}

	config := oauth1.NewConfig(apiKey, apiSecretKey)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	p.client = *client

	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"twitter_tweet": tweetResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"twitter_tweet": tweetDataSourceType{},
	}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"api_key": {
				MarkdownDescription: "Twitter API key",
				Optional:            true,
				Type:                types.StringType,
				Sensitive:           true,
			},
			"api_secret_key": {
				MarkdownDescription: "Twitter API secret key",
				Optional:            true,
				Type:                types.StringType,
				Sensitive:           true,
			},
			"access_token": {
				MarkdownDescription: "Twitter access token",
				Optional:            true,
				Type:                types.StringType,
				Sensitive:           true,
			},
			"access_token_secret": {
				MarkdownDescription: "Twitter access token secret",
				Optional:            true,
				Type:                types.StringType,
				Sensitive:           true,
			},
		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
