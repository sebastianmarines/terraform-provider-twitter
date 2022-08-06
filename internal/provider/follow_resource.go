package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/avast/retry-go"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sebastianmarines/terraform-provider-twitter/internal/utils"
)

var _ tfsdk.ResourceType = profileResourceType{}
var _ tfsdk.Resource = profileResource{}

type followResourceType struct{}

func (t followResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Allows the authenticating user to follow (friend) a user.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the user being followed.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"screen_name": {
				MarkdownDescription: "The screen name of the user being followed.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"user_id": {
				MarkdownDescription: "The ID of the user being followed.",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"pending": {
				MarkdownDescription: "Whether the authenticated user is pending approval to follow the user.",
				Type:                types.BoolType,
				Computed:            true,
			},
		},
	}, nil
}

func (t followResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return followResource{
		provider: provider,
	}, diags
}

type followResourceData struct {
	ID         types.Int64  `tfsdk:"id"`
	ScreenName types.String `tfsdk:"screen_name"`
	UserId     types.Int64  `tfsdk:"user_id"`
	Pending    types.Bool   `tfsdk:"pending"`
}

type followResource struct {
	provider provider
}

func (t followResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {

	err := utils.CheckProviderConfiguration(&resp.Diagnostics, t.provider.configured)
	if err != nil {
		return
	}

	var data followResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ScreenName.Null && data.UserId.Null {
		resp.Diagnostics.AddError(
			"Could not follow user",
			"Must specify either screen_name or user_id",
		)
		return
	}

	params := &twitter.FriendshipCreateParams{
		Follow: twitter.Bool(true),
	}

	if !data.ScreenName.Null {
		params.ScreenName = data.ScreenName.Value
	}
	if !data.UserId.Null {
		params.UserID = data.UserId.Value
	}

	user, _, err := t.provider.client.Friendships.Create(params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not follow user",
			fmt.Sprintf("Unable to follow user, got error %s", err),
		)
		return
	}

	follow := &followResourceData{}
	follow.Pending.Value = user.FollowRequestSent
	follow.ScreenName.Value = user.ScreenName
	follow.UserId.Value = user.ID
	follow.ID.Value = user.ID

	diags = resp.State.Set(ctx, follow)
	resp.Diagnostics.Append(diags...)
}

func (r followResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	err := utils.CheckProviderConfiguration(&resp.Diagnostics, r.provider.configured)
	if err != nil {
		return
	}

	var data followResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := &twitter.UserShowParams{}

	if !data.ScreenName.Null {
		params.ScreenName = data.ScreenName.Value
	}
	if !data.UserId.Null {
		params.UserID = data.UserId.Value
	}

	user, _, err := r.provider.client.Users.Show(params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read profile",
			fmt.Sprintf("Unable to read profile, got error %s", err),
		)
		return
	}

	if !user.FollowRequestSent && !user.Following {
		resp.State.RemoveResource(ctx)
		return
	}

	follow := &followResourceData{}
	follow.Pending.Value = user.FollowRequestSent
	follow.ScreenName.Value = user.ScreenName
	follow.UserId.Value = user.ID
	follow.ID.Value = user.ID

	diags = resp.State.Set(ctx, &follow)
	resp.Diagnostics.Append(diags...)
}

func (r followResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"Update is not supported for follow resource",
	)
	return
}

func (r followResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	err := utils.CheckProviderConfiguration(&resp.Diagnostics, r.provider.configured)
	if err != nil {
		return
	}

	var data followResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := &twitter.FriendshipDestroyParams{}

	if !data.ScreenName.Null {
		params.ScreenName = data.ScreenName.Value
	}
	if !data.UserId.Null {
		params.UserID = data.UserId.Value
	}

	_, _, err = r.provider.client.Friendships.Destroy(params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not unfollow user",
			fmt.Sprintf("Unable to unfollow user, got error %s", err),
		)
		return
	}

	err = retry.Do(
		func() error {
			user, _, err := r.provider.client.Users.Show(&twitter.UserShowParams{
				ScreenName: data.ScreenName.Value,
				UserID:     data.UserId.Value,
			})
			if err != nil {
				return err
			}

			if user.FollowRequestSent || user.Following {
				return errors.New("Unable to unfollow user")
			}

			return nil
		},
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not unfollow user",
			fmt.Sprintf("Unable to unfollow user"),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}
