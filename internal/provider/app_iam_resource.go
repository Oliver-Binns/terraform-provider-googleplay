package provider

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/oliver-binns/googleplay-go"
	"github.com/oliver-binns/googleplay-go/users"
)

var _ resource.Resource = &AppIAMResource{}
var _ resource.ResourceWithValidateConfig = &AppIAMResource{}

func NewAppIAMResource() resource.Resource {
	return &AppIAMResource{}
}

type AppIAMResource struct {
	client *googleplay.Client
}

type appIAMResourceModel struct {
	UserID              types.String `tfsdk:"user_id"`
	AppID               types.String `tfsdk:"app_id"`
	Permissions         types.Set    `tfsdk:"permissions"`
	ExpandedPermissions types.Set    `tfsdk:"expanded_permissions"`
}

func (m *appIAMResourceModel) Name(developerID string) string {
	return fmt.Sprintf("developers/%s/users/%s/grants/%s", developerID, m.UserID, m.AppID)
}

func (r *AppIAMResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_iam"
}

func (r *AppIAMResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Manage app level access to the Google Play Console",

		Attributes: map[string]schema.Attribute{
			"user_id": schema.StringAttribute{
				MarkdownDescription: "The ID for the user: this is the email they use to login to Google Play",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "The app / package ID to grant access to",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"permissions": schema.SetAttribute{
				MarkdownDescription: `Permissions for the user which apply to this specific app:
				https://developers.google.com/android-publisher/api-ref/rest/v3/grants#applevelpermission`,
				ElementType: types.StringType,
				Required:    true,
			},
			"expanded_permissions": schema.SetAttribute{
				MarkdownDescription: `Permissions for the user which apply to this specific app:
				https://developers.google.com/android-publisher/api-ref/rest/v3/grants#applevelpermission`,
				ElementType: types.StringType,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					expandAppPermissionsPlanModifier(path.Root("permissions")),
				},
			},
		},
	}
}

func (r *AppIAMResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AppIAMResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data appIAMResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Each grant must contain a valid permission
	if len(data.Permissions.Elements()) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("permissions"),
			"Invalid permissions configuration",
			"permissions must contain at least one permission.",
		)
	}

	// Warn user if they are using an implicit permission
	permissions := []users.AppLevelPermission{}
	diag := data.Permissions.ElementsAs(ctx, &permissions, false)
	resp.Diagnostics.Append(diag...)
	for _, permission := range permissions {
		for _, inherited := range permission.Expand() {
			if !slices.Contains(permissions, inherited) {
				resp.Diagnostics.AddWarning(
					"Granting implicit permission",
					fmt.Sprintf(
						"The permission '%s' is inherited from '%s', but it is not explicitly granted.",
						inherited, permission,
					),
				)
			}
		}
	}
}

func (r *AppIAMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data appIAMResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	permissions := []users.AppLevelPermission{}
	diag := data.Permissions.ElementsAs(ctx, &permissions, false)
	resp.Diagnostics.Append(diag...)

	grant, err := r.client.GrantAccess(
		data.UserID.ValueString(),
		data.AppID.ValueString(),
		permissions,
		ctx,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to grant access to app:",
			err.Error(),
		)
		return
	}

	// App ID is actually:
	// developers/DEVELOPER_ID/users/EMAIL/grants/APP_ID
	components := strings.Split(grant.Name, "/")
	data.UserID = types.StringValue(components[3])
	data.AppID = types.StringValue(components[5])

	data.ExpandedPermissions, diag = types.SetValueFrom(ctx, types.StringType, grant.AppLevelPermissions)
	resp.Diagnostics.Append(diag...)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppIAMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data appIAMResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userID := data.UserID.ValueString()
	if userID == "" {
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"Attribute 'userID' is required to fetch user information.",
		)
		return
	}

	appID := data.AppID.ValueString()
	if appID == "" {
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"Attribute 'appID' is required to fetch IAM information.",
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
		if user.Email == userID {
			for _, grant := range user.Grants {
				if grant.Name == appID {
					// App ID is actually:
					// developers/DEVELOPER_ID/users/EMAIL/grants/APP_ID
					components := strings.Split(grant.Name, "/")
					data.UserID = types.StringValue(components[3])
					data.AppID = types.StringValue(components[5])

					var diag diag.Diagnostics
					data.Permissions, diag = types.SetValueFrom(ctx, types.StringType, grant.AppLevelPermissions)
					resp.Diagnostics.Append(diag...)
					break
				}
			}
			break
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppIAMResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data appIAMResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	permissions := []users.AppLevelPermission{}
	diag := data.Permissions.ElementsAs(ctx, &permissions, false)
	resp.Diagnostics.Append(diag...)

	grant, err := r.client.ModifyAccess(
		data.UserID.ValueString(),
		data.AppID.ValueString(),
		permissions,
		ctx,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update IAM permissions",
			err.Error(),
		)
		return
	}

	// App ID is actually:
	// developers/DEVELOPER_ID/users/EMAIL/grants/APP_ID
	components := strings.Split(grant.Name, "/")
	data.UserID = types.StringValue(components[3])
	data.AppID = types.StringValue(components[5])

	data.ExpandedPermissions, diag = types.SetValueFrom(ctx, types.StringType, grant.AppLevelPermissions)
	resp.Diagnostics.Append(diag...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppIAMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data appIAMResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.RevokeAccess(
		data.UserID.ValueString(),
		data.AppID.ValueString(),
		ctx,
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete permission, got error: %s", err))
		return
	}
}
