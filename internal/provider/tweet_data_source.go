package provider

import (
	"context"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sebastianmarines/terraform-provider-twitter/internal/utils"
)

var _ tfsdk.DataSourceType = tweetDataSourceType{}
var _ tfsdk.DataSource = tweetDataSource{}

type tweetDataSourceType struct{}

func (t tweetDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Tweet data source",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The integer representation of the unique identifier for this Tweet.",
				Type:                types.Int64Type,
				Required:            true,
			},
			"text": {
				MarkdownDescription: "The actual UTF-8 text of the status update.",
				Type:                types.StringType,
				Computed:            true,
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

func (t tweetDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return tweetDataSource{
		provider: provider,
	}, diags
}

type tweetDataSourceData struct {
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

type tweetDataSource struct {
	provider provider
}

func (d tweetDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	err := utils.CheckProviderConfiguration(&resp.Diagnostics, d.provider.configured)
	if err != nil {
		return
	}

	var data tweetDataSourceData

	diags := req.Config.Get(ctx, &data)
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

	tweet, _, err := d.provider.client.Statuses.Show(data.ID.Value, params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read tweet",
			fmt.Sprintf("Unable to read tweet, got error: %s", err),
		)
		return
	}

	newTweet := &tweetDataSourceData{}

	newTweet.Text.Value = tweet.Text
	newTweet.UserID.Value = tweet.User.ID
	newTweet.Source.Value = tweet.Source
	newTweet.InReplyToStatusID.Value = tweet.InReplyToStatusID
	newTweet.InReplyToUserID.Value = tweet.InReplyToUserID
	newTweet.QuotedStatusID.Value = tweet.QuotedStatusID
	newTweet.QuoteCount.Value = int64(tweet.QuoteCount)
	newTweet.ReplyCount.Value = int64(tweet.ReplyCount)
	newTweet.RetweetCount.Value = int64(tweet.RetweetCount)
	newTweet.FavoriteCount.Value = int64(tweet.FavoriteCount)
	newTweet.PossiblySensitive.Value = tweet.PossiblySensitive
	newTweet.Lang.Value = tweet.Lang

	diags = resp.State.Set(ctx, &newTweet)
	resp.Diagnostics.Append(diags...)

}
