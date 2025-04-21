package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/oliver-binns/googleplay-go"
	"github.com/oliver-binns/googleplay-go/users"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &UserResource{}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

type UserResource struct {
	client *googleplay.Client
}

type userResourceModel struct {
	Name              types.String `tfsdk:"name"`
	Email             types.String `tfsdk:"email"`
	GlobalPermissions types.Set    `tfsdk:"global_permissions"`
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Manage user accounts in the Google Play Console",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the user",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email address for the user",
				Required:            true,
			},
			"global_permissions": schema.SetAttribute{
				MarkdownDescription: `Permissions for the user which apply across the developer account:
				https://developers.google.com/android-publisher/api-ref/rest/v3/users#DeveloperLevelPermission`,
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*googleplay.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *googleplay.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data userResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	permissions := []users.DeveloperLevelPermission{}
	diag := data.GlobalPermissions.ElementsAs(ctx, &permissions, false)
	resp.Diagnostics.Append(diag...)

	user, err := r.client.CreateUser(
		data.Email.ValueString(),
		permissions,
		ctx,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create user",
			err.Error(),
		)
		return
	}

	data.Name = types.StringValue(user.Name)
	data.Email = types.StringValue(user.Email)

	data.GlobalPermissions, diag = types.SetValueFrom(ctx, types.StringType, user.DeveloperAccountPermissions)
	resp.Diagnostics.Append(diag...)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data userResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	email := data.Email.ValueString()
	if email == "" {
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"Attribute 'email' is required to fetch user information.",
		)
		return
	}

	// List users from Google Play API
	users, err := r.client.ListUsers(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to fetch users",
			err.Error(),
		)
		return
	}

	for _, user := range users {
		if user.Email == email {
			data.Name = types.StringValue(user.Name)
			data.Email = types.StringValue(user.Email)

			var diag diag.Diagnostics
			data.GlobalPermissions, diag = types.SetValueFrom(ctx, types.StringType, user.DeveloperAccountPermissions)
			resp.Diagnostics.Append(diag...)
			break
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data userResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	permissions := []users.DeveloperLevelPermission{}
	diag := data.GlobalPermissions.ElementsAs(ctx, &permissions, false)
	resp.Diagnostics.Append(diag...)

	user, err := r.client.UpdateUser(
		data.Email.ValueString(),
		&permissions,
		ctx,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update user",
			err.Error(),
		)
		return
	}

	data.Email = types.StringValue(user.Email)
	data.Name = types.StringValue(user.Name)

	data.GlobalPermissions, diag = types.SetValueFrom(ctx, types.StringType, user.DeveloperAccountPermissions)
	resp.Diagnostics.Append(diag...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data userResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(data.Email.ValueString(), ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete user, got error: %s", err))
		return
	}
}
