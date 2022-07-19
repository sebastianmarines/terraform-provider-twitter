package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ tfsdk.ResourceType = profileResourceType{}
var _ tfsdk.Resource = profileResource{}

type profileResourceType struct{}

func (t profileResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Sets some values that users are able to set under the \"Account\" tab of their settings page. ",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The integer representation of the unique identifier for this User.",
				Type:                types.Int64Type,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Full name associated with the profile.",
				Type:                types.StringType,
				Optional:            true,
			},
			"url": {
				MarkdownDescription: "URL associated with the profile.",
				Type:                types.StringType,
				Optional:            true,
			},
			"location": {
				MarkdownDescription: "The city or country describing where the user of the account is located. The contents are not normalized or geocoded in any way.",
				Type:                types.StringType,
				Optional:            true,
			},
			"description": {
				MarkdownDescription: "A description of the user owning the account.",
				Type:                types.StringType,
				Optional:            true,
			},
		},
	}, nil
}

func (t profileResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return profileResource{
		provider: provider,
	}, diags
}

type profileResourceData struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	URL         types.String `tfsdk:"url"`
	Location    types.String `tfsdk:"location"`
	Description types.String `tfsdk:"description"`
}

type profileResource struct {
	provider provider
}

func (t profileResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data profileResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := &twitter.AccountUpdateProfileParams{}

	if !data.Name.Null {
		params.Name = data.Name.Value
	}

	if !data.URL.Null {
		params.URL = data.URL.Value
	}

	if !data.Location.Null {
		params.Location = data.Location.Value
	}

	if !data.Description.Null {
		params.Description = data.Description.Value
	}

	user, _, err := t.provider.client.Accounts.UpdateProfile(params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not update profile",
			fmt.Sprintf("Unable to update profile, got error %s", err.Error()),
		)
		return
	}

	data.ID = types.Int64{Value: user.ID}

	if !data.Name.Null {
		data.Name = types.String{Value: user.Name}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r profileResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data profileResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := &twitter.UserShowParams{
		UserID: data.ID.Value,
	}

	user, _, err := r.provider.client.Users.Show(params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read user",
			fmt.Sprintf("Unable to read user, got error: %s", err),
		)
		return
	}

	// data.ID = types.Int64{Value: user.ID}
	data.Name = types.String{Value: user.Name}
	data.URL = types.String{Value: user.URL}
	data.Location = types.String{Value: user.Location}
	data.Description = types.String{Value: user.Description}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r profileResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data profileResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	urlLocation := "https://api.twitter.com/1.1/account/update_profile.json?name=%s&url=%s&location=%s&description=%s"

	urlLocation = fmt.Sprintf(urlLocation, url.QueryEscape(data.Name.Value), url.QueryEscape(data.URL.Value), url.QueryEscape(data.Location.Value), url.QueryEscape(data.Description.Value))

	_req, _ := http.NewRequest("POST", urlLocation, nil)

	_res, _ := r.provider.httpClient.Do(_req)

	if _res.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Could not update profile",
			fmt.Sprintf("Unable to update profile, got error %s", _res.Status),
		)
		return
	}

	defer _res.Body.Close()

	body, _ := ioutil.ReadAll(_res.Body)

	profile := &profileHTTPResponse{}

	_ = json.Unmarshal(body, profile)

	data.ID = types.Int64{Value: profile.ID}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r profileResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data profileResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	urlLocation := "https://api.twitter.com/1.1/account/update_profile.json?url=&location=&description="

	_req, _ := http.NewRequest("POST", urlLocation, nil)

	_res, _ := r.provider.httpClient.Do(_req)

	if _res.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Could not delete profile",
			fmt.Sprintf("Unable to delete profile, got error %s", _res.Status),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

type profileHTTPResponse struct {
	ID int64 `json:"id"`
}
