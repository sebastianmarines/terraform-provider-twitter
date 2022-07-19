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

type userDataSourceType struct{}

func (t userDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "User data source",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The integer representation of the unique identifier for this User.",
				Type:                types.Int64Type,
				Optional:            true,
			},
			"screen_name": {
				MarkdownDescription: "The screen name, handle, or alias that this user identifies themselves with.",
				Type:                types.StringType,
				Optional:            true,
			},
			"name": {
				MarkdownDescription: "The name of the user, as they’ve defined it.",
				Type:                types.StringType,
				Computed:            true,
			},
			"location": {
				MarkdownDescription: "The user-defined location for this account’s profile. Not necessarily a location, nor machine-parseable.",
				Type:                types.StringType,
				Computed:            true,
			},
			"url": {
				MarkdownDescription: "The URL provided by the user in association with their profile. May be absent.",
				Type:                types.StringType,
				Computed:            true,
			},
			"description": {
				MarkdownDescription: "The user-defined UTF-8 string describing their account.",
				Type:                types.StringType,
				Computed:            true,
			},
			"protected": {
				MarkdownDescription: "When true, indicates that this user has chosen to protect their Tweets.",
				Type:                types.BoolType,
				Computed:            true,
			},
			"verified": {
				MarkdownDescription: "When true, indicates that this user has a verified account.",
				Type:                types.BoolType,
				Computed:            true,
			},
			"followers_count": {
				MarkdownDescription: "The number of followers this account currently has.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"friends_count": {
				MarkdownDescription: "The number of users this account is following.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"statuses_count": {
				MarkdownDescription: "The number of Tweets (including retweets) issued by the user.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"favorites_count": {
				MarkdownDescription: "The number of Tweets this user has liked in the account’s lifetime.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"profile_banner_url": {
				MarkdownDescription: "The HTTPS-based URL pointing to the standard web representation of the user’s uploaded profile banner.",
				Type:                types.StringType,
				Computed:            true,
			},
			"profile_image_url": {
				MarkdownDescription: "A HTTPS-based URL pointing to the user’s profile image.",
				Type:                types.StringType,
				Computed:            true,
			},
			"default_profile": {
				MarkdownDescription: "When true, indicates that the user has not altered the theme or background of their user profile.",
				Type:                types.BoolType,
				Computed:            true,
			},
			"default_profile_image": {
				MarkdownDescription: "When true, indicates that the user has not uploaded their own profile image and a default image is used instead.",
				Type:                types.BoolType,
				Computed:            true,
			},
		},
	}, nil
}

func (t userDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return userDataSource{
		provider: provider,
	}, diags
}

type userDataSourceData struct {
	ID                   types.Int64  `tfsdk:"id"`
	ScreenName           types.String `tfsdk:"screen_name"`
	Name                 types.String `tfsdk:"name"`
	Location             types.String `tfsdk:"location"`
	URL                  types.String `tfsdk:"url"`
	Description          types.String `tfsdk:"description"`
	Protected            types.Bool   `tfsdk:"protected"`
	Verified             types.Bool   `tfsdk:"verified"`
	Followers            types.Int64  `tfsdk:"followers_count"`
	Friends              types.Int64  `tfsdk:"friends_count"`
	Statuses             types.Int64  `tfsdk:"statuses_count"`
	Favorites            types.Int64  `tfsdk:"favorites_count"`
	ProfileBannerURL     types.String `tfsdk:"profile_banner_url"`
	ProfileImageURLHttps types.String `tfsdk:"profile_image_url"`
	DefaultProfile       types.Bool   `tfsdk:"default_profile"`
	DefaultProfileImage  types.Bool   `tfsdk:"default_profile_image"`
}

type userDataSource struct {
	provider provider
}

func (d userDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data userDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.Null && data.ScreenName.Null {
		resp.Diagnostics.AddError(
			"Missing required argument",
			"Either a id or screen_name is required, but no definition was found.",
		)
		return
	}

	params := &twitter.UserShowParams{
		IncludeEntities: twitter.Bool(false),
	}

	if !data.ID.Null {
		params.UserID = data.ID.Value
	}

	if !data.ScreenName.Null {
		params.ScreenName = data.ScreenName.Value
	}

	tweet, _, err := d.provider.client.Users.Show(params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read user",
			fmt.Sprintf("Unable to read user, got error: %s", err),
		)
		return
	}

	data.ID = types.Int64{Value: tweet.ID}
	data.ScreenName = types.String{Value: tweet.ScreenName}
	data.Name = types.String{Value: tweet.Name}
	data.Location = types.String{Value: tweet.Location}
	data.URL = types.String{Value: tweet.URL}
	data.Description = types.String{Value: tweet.Description}
	data.Protected = types.Bool{Value: tweet.Protected}
	data.Verified = types.Bool{Value: tweet.Verified}
	data.Followers = types.Int64{Value: int64(tweet.FollowersCount)}
	data.Friends = types.Int64{Value: int64(tweet.FriendsCount)}
	data.Statuses = types.Int64{Value: int64(tweet.StatusesCount)}
	data.Favorites = types.Int64{Value: int64(tweet.FavouritesCount)}
	data.ProfileBannerURL = types.String{Value: tweet.ProfileBannerURL}
	data.ProfileImageURLHttps = types.String{Value: tweet.ProfileImageURLHttps}
	data.DefaultProfile = types.Bool{Value: tweet.DefaultProfile}
	data.DefaultProfileImage = types.Bool{Value: tweet.DefaultProfileImage}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

}
