package provider

import (
	"context"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ tfsdk.ResourceType = tweetResourceType{}
var _ tfsdk.Resource = tweetResource{}

type tweetResourceType struct{}

func (t tweetResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Tweet resource",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Tweet id",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"text": {
				MarkdownDescription: "Tweet text",
				Type:                types.StringType,
				Required:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
				// TODO Add length validation
			},
			"user_id": {
				MarkdownDescription: "Tweet user id",
				Type:                types.Int64Type,
				Computed:            true,
			},
		},
	}, nil
}

func (t tweetResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return tweetResource{
		provider: provider,
	}, diags
}

type tweetResourceData struct {
	ID     types.Int64  `tfsdk:"id"`
	Text   types.String `tfsdk:"text"`
	UserID types.Int64  `tfsdk:"user_id"`
}

type tweetResource struct {
	provider provider
}

func (t tweetResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data tweetResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := &twitter.StatusUpdateParams{
		Status:   data.Text.Value,
		TrimUser: twitter.Bool(true),
	}

	tweet, _, err := t.provider.client.Statuses.Update(data.Text.Value, params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create tweet",
			fmt.Sprintf("Unable to create tweet, got error %s", err.Error()),
		)
		return
	}

	data.ID = types.Int64{Value: tweet.ID}
	data.Text = types.String{Value: tweet.Text}
	data.UserID = types.Int64{Value: tweet.User.ID}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r tweetResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data tweetResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := &twitter.StatusShowParams{
		ID:               data.ID.Value,
		TrimUser:         twitter.Bool(true),
		IncludeMyRetweet: twitter.Bool(false),
		IncludeEntities:  twitter.Bool(false),
	}

	tweet, response, err := r.provider.client.Statuses.Show(data.ID.Value, params)

	if err != nil {
		if response.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(
				"Could not read tweet",
				fmt.Sprintf("Unable to read tweet, got error: %d, %s", response.StatusCode, err),
			)
			return
		}
	}

	data.ID = types.Int64{Value: tweet.ID}
	data.Text = types.String{Value: tweet.Text}
	data.UserID = types.Int64{Value: tweet.User.ID}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r tweetResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"Update is not supported for tweet resource",
	)
	return
}

func (r tweetResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data tweetResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := &twitter.StatusDestroyParams{
		TrimUser: twitter.Bool(true),
	}

	_, _, err := r.provider.client.Statuses.Destroy(data.ID.Value, params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not delete tweet",
			fmt.Sprintf("Unable to delete tweet with ID %s, got error: %s", data.ID.String(), err),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}
