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
				Computed:            true,
			},
			"screen_name": {
				MarkdownDescription: "The screen name, handle, or alias that this user identifies themselves with.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
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

	user, _, err := d.provider.client.Users.Show(params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read user",
			fmt.Sprintf("Unable to read user, got error: %s", err),
		)
		return
	}

	newUser := &userDataSourceData{}

	newUser.ID.Value = user.ID
	newUser.ScreenName.Value = user.ScreenName
	newUser.Name.Value = user.Name
	newUser.Location.Value = user.Location
	newUser.URL.Value = user.URL
	newUser.Description.Value = user.Description
	newUser.Protected.Value = user.Protected
	newUser.Verified.Value = user.Verified
	newUser.Followers.Value = int64(user.FollowersCount)
	newUser.Friends.Value = int64(user.FriendsCount)
	newUser.Statuses.Value = int64(user.StatusesCount)
	newUser.Favorites.Value = int64(user.FavouritesCount)
	newUser.ProfileBannerURL.Value = user.ProfileBannerURL
	newUser.ProfileImageURLHttps.Value = user.ProfileImageURLHttps
	newUser.DefaultProfile.Value = user.DefaultProfile
	newUser.DefaultProfileImage.Value = user.DefaultProfileImage

	diags = resp.State.Set(ctx, &newUser)
	resp.Diagnostics.Append(diags...)

}
