package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sebastianmarines/terraform-provider-twitter/internal/modifiers"
	"github.com/sebastianmarines/terraform-provider-twitter/internal/utils"
)

var _ tfsdk.ResourceType = listResourceType{}
var _ tfsdk.Resource = listResource{}

type listResourceType struct{}

func (t listResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "A list is a curated group of Twitter accounts. You can create your own lists or subscribe to lists created by others for the authenticated user. Viewing a list timeline will show you a stream of Tweets from only the accounts on that list.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The numerical id of the list.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"slug": {
				MarkdownDescription: "The slug of the list.",
				Type:                types.StringType,
				Computed:            true,
			},
			"name": {
				MarkdownDescription: "The name for the list. A list's name must start with a letter and can consist only of 25 or fewer letters, numbers, \"-\", or \"_\" characters.",
				Type:                types.StringType,
				Required:            true,
				Validators:          []tfsdk.AttributeValidator{
					// TODO: add validator for list name
				},
			},
			"created_at": {
				MarkdownDescription: "The date and time when the list was created.",
				Type:                types.StringType,
				Computed:            true,
			},
			"uri": {
				MarkdownDescription: "The URI of the list.",
				Type:                types.StringType,
				Computed:            true,
			},
			"subscriber_count": {
				MarkdownDescription: "The number of subscribers that the list has.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"member_count": {
				MarkdownDescription: "The number of members that the list has.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"mode": {
				MarkdownDescription: "Whether your list is public or private. Values can be public or private . If no mode is specified the list will be public.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.StringDefault("public"),
				},
				Validators: []tfsdk.AttributeValidator{
					// TODO: add validator for list mode
				},
			},
			"full_name": {
				MarkdownDescription: "The full name of the list.",
				Type:                types.StringType,
				Computed:            true,
			},
			"description": {
				MarkdownDescription: "The description of the list.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"user_id": {
				MarkdownDescription: "The user ID of the list's owner.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			// "members": {
			// 	MarkdownDescription: "The members of the list.",
			// 	Type:                types.SetType{ElemType: types.StringType},
			// 	Optional:            true,
			// 	Validators:          []tfsdk.AttributeValidator{
			// 		// TODO: add validator for list members length
			// 	},
			// },
		},
	}, nil
}

func (t listResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return listResource{
		provider: provider,
	}, diags
}

type listResourceData struct {
	ID              types.Int64  `tfsdk:"id"`
	Slug            types.String `tfsdk:"slug"`
	Name            types.String `tfsdk:"name"`
	CreatedAt       types.String `tfsdk:"created_at"`
	URI             types.String `tfsdk:"uri"`
	SubscriberCount types.Int64  `tfsdk:"subscriber_count"`
	MemberCount     types.Int64  `tfsdk:"member_count"`
	Mode            types.String `tfsdk:"mode"`
	FullName        types.String `tfsdk:"full_name"`
	Description     types.String `tfsdk:"description"`
	UserID          types.Int64  `tfsdk:"user_id"`
	// Members         types.Set    `tfsdk:"members"`
}

type listResource struct {
	provider provider
}

func (t listResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {

	err := utils.CheckProviderConfiguration(&resp.Diagnostics, t.provider.configured)
	if err != nil {
		return
	}

	var data listResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := &twitter.ListsCreateParams{
		Name: data.Name.Value,
	}

	if !data.Description.Null {
		params.Description = data.Description.Value
	}

	if !data.Mode.Null {
		params.Mode = data.Mode.Value
	}

	list, _, err := t.provider.client.Lists.Create(data.Name.Value, params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create list",
			fmt.Sprintf("Could not create list, got error: %s", err),
		)
		return
	}

	newList := &listResourceData{}

	newList.ID.Value = list.ID
	newList.Slug.Value = list.Slug
	newList.CreatedAt.Value = list.CreatedAt
	newList.URI.Value = list.URI
	newList.SubscriberCount.Value = int64(list.SubscriberCount)
	newList.MemberCount.Value = int64(list.MemberCount)
	newList.Mode.Value = list.Mode
	newList.FullName.Value = list.FullName
	newList.UserID.Value = int64(list.User.ID)
	// newList.Description.Value = list.Description
	newList.Description.Value = data.Description.Value
	newList.Name.Value = list.Name

	// Log the user ID
	log.Printf("[DEBUG] List user ID: %d", newList.UserID.Value)

	diags = resp.State.Set(ctx, newList)
	resp.Diagnostics.Append(diags...)
}

func (r listResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	err := utils.CheckProviderConfiguration(&resp.Diagnostics, r.provider.configured)
	if err != nil {
		return
	}

	var data listResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := &twitter.ListsShowParams{
		ListID: data.ID.Value,
	}

	if !data.Slug.Null {
		params.Slug = data.Slug.Value
	}

	if !data.UserID.Null {
		params.OwnerID = data.UserID.Value
	}

	list, _, err := r.provider.client.Lists.Show(params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read list",
			fmt.Sprintf("Could not read list, got error: %s", err),
		)
		return
	}

	newList := &listResourceData{}

	newList.ID.Value = list.ID
	newList.Slug.Value = list.Slug
	newList.CreatedAt.Value = list.CreatedAt
	newList.URI.Value = list.URI
	newList.SubscriberCount.Value = int64(list.SubscriberCount)
	newList.MemberCount.Value = int64(list.MemberCount)
	newList.Mode.Value = list.Mode
	newList.FullName.Value = list.FullName
	newList.UserID.Value = int64(list.User.ID)
	// newList.Description.Value = list.Description
	newList.Description.Value = data.Description.Value
	newList.Name.Value = list.Name

	diags = req.State.Set(ctx, newList)
	resp.Diagnostics.Append(diags...)
}

func (r listResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"Update is not supported for list resource",
	)
	return
}

func (r listResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	err := utils.CheckProviderConfiguration(&resp.Diagnostics, r.provider.configured)
	if err != nil {
		return
	}

	var data listResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := &twitter.ListsDestroyParams{
		ListID: data.ID.Value,
	}

	_, _, err = r.provider.client.Lists.Destroy(params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not delete list",
			fmt.Sprintf("Could not delete list, got error: %s", err),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}
