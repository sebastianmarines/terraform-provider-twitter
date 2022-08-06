package provider

import (
	"context"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sebastianmarines/terraform-provider-twitter/internal/validators"
)

var _ tfsdk.ResourceType = tweetResourceType{}
var _ tfsdk.Resource = tweetResource{}

type tweetResourceType struct{}

func (t tweetResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Tweet resource",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The integer representation of the unique identifier for this Tweet.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"text": {
				MarkdownDescription: "The actual UTF-8 text of the status update. Should not exceed 280 characters.",
				Type:                types.StringType,
				Required:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
				Validators: []tfsdk.AttributeValidator{
					validators.TweetLength(),
				},
			},
			"user_id": {
				MarkdownDescription: "The integer representation of the unique identifier for the user who posted this Tweet.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"source": {
				MarkdownDescription: "Utility used to post the Tweet, as an HTML-formatted string. ",
				Type:                types.StringType,
				Computed:            true,
			},
			"in_reply_to_status_id": {
				MarkdownDescription: "If the represented Tweet is a reply, this field will contain the integer representation of the original Tweet’s ID.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"in_reply_to_user_id": {
				MarkdownDescription: "If the represented Tweet is a reply, this field will contain the integer representation of the original Tweet’s author ID.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"quoted_status_id": {
				MarkdownDescription: "This field only surfaces when the Tweet is a quote Tweet. This field contains the integer value Tweet ID of the quoted Tweet.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"quote_count": {
				MarkdownDescription: "Indicates approximately how many times this Tweet has been quoted by Twitter users.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"reply_count": {
				MarkdownDescription: "Number of times this Tweet has been replied to.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"retweet_count": {
				MarkdownDescription: "Number of times this Tweet has been retweeted.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"favorite_count": {
				MarkdownDescription: "Indicates approximately how many times this Tweet has been liked by Twitter users.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"possibly_sensitive": {
				MarkdownDescription: "An indicator that the URL contained in the Tweet may contain content or media identified as sensitive content.",
				Type:                types.BoolType,
				Computed:            true,
			},
			"lang": {
				MarkdownDescription: "When present, indicates a BCP 47 language identifier corresponding to the machine-detected language of the Tweet text, or und if no language could be detected.",
				Type:                types.StringType,
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
	ID                types.Int64  `tfsdk:"id"`
	Text              types.String `tfsdk:"text"`
	UserID            types.Int64  `tfsdk:"user_id"`
	Source            types.String `tfsdk:"source"`
	InReplyToStatusID types.Int64  `tfsdk:"in_reply_to_status_id"`
	InReplyToUserID   types.Int64  `tfsdk:"in_reply_to_user_id"`
	QuotedStatusID    types.Int64  `tfsdk:"quoted_status_id"`
	QuoteCount        types.Int64  `tfsdk:"quote_count"`
	ReplyCount        types.Int64  `tfsdk:"reply_count"`
	RetweetCount      types.Int64  `tfsdk:"retweet_count"`
	FavoriteCount     types.Int64  `tfsdk:"favorite_count"`
	PossiblySensitive types.Bool   `tfsdk:"possibly_sensitive"`
	Lang              types.String `tfsdk:"lang"`
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

	data.ID.Value = tweet.ID
	data.Text.Value = tweet.Text
	data.UserID.Value = tweet.User.ID
	data.Source.Value = tweet.Source
	data.InReplyToStatusID.Value = tweet.InReplyToStatusID
	data.InReplyToUserID.Value = tweet.InReplyToUserID
	data.QuotedStatusID.Value = tweet.QuotedStatusID
	data.QuoteCount.Value = int64(tweet.QuoteCount)
	data.ReplyCount.Value = int64(tweet.ReplyCount)
	data.RetweetCount.Value = int64(tweet.RetweetCount)
	data.FavoriteCount.Value = int64(tweet.FavoriteCount)
	data.PossiblySensitive.Value = tweet.PossiblySensitive
	data.Lang.Value = tweet.Lang

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

	data.ID.Value = tweet.ID
	data.Text.Value = tweet.Text
	data.UserID.Value = tweet.User.ID

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
