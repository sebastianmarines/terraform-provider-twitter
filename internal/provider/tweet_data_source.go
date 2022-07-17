package provider

import (
	"context"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
				MarkdownDescription: "Tweet id",
				Type:                types.Int64Type,
				Required:            true,
			},
			"text": {
				MarkdownDescription: "Tweet text",
				Type:                types.StringType,
				Computed:            true,
			},
			"user_id": {
				MarkdownDescription: "Tweet user id",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"favorite_count": {
				MarkdownDescription: "Tweet favorite count",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"reply_count": {
				MarkdownDescription: "Tweet reply count",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"retweet_count": {
				MarkdownDescription: "Tweet retweet count",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"full_text": {
				MarkdownDescription: "Tweet full text",
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
	ID            types.Int64  `tfsdk:"id"`
	Text          types.String `tfsdk:"text"`
	UserID        types.Int64  `tfsdk:"user_id"`
	FavoriteCount types.Int64  `tfsdk:"favorite_count"`
	ReplyCount    types.Int64  `tfsdk:"reply_count"`
	RetweetCount  types.Int64  `tfsdk:"retweet_count"`
	FullText      types.String `tfsdk:"full_text"`
}

type tweetDataSource struct {
	provider provider
}

func (d tweetDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
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

	data.Text = types.String{Value: tweet.Text}
	data.UserID = types.Int64{Value: tweet.User.ID}
	data.FavoriteCount = types.Int64{Value: int64(tweet.FavoriteCount)}
	data.ReplyCount = types.Int64{Value: int64(tweet.ReplyCount)}
	data.RetweetCount = types.Int64{Value: int64(tweet.RetweetCount)}
	data.FullText = types.String{Value: tweet.FullText}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

}
