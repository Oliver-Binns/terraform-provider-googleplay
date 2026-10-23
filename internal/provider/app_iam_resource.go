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
	androidpublisher "google.golang.org/api/androidpublisher/v3"
)

var _ resource.Resource = &AppIAMResource{}
var _ resource.ResourceWithValidateConfig = &AppIAMResource{}
var _ resource.ResourceWithImportState = &AppIAMResource{}

func NewAppIAMResource() resource.Resource {
	return &AppIAMResource{}
}

type AppIAMResource struct {
	client *GooglePlayClient
}

type appIAMResourceModel struct {
	UserID              types.String `tfsdk:"user_id"`
	AppID               types.String `tfsdk:"app_id"`
	Permissions         types.Set    `tfsdk:"permissions"`
	ExpandedPermissions types.Set    `tfsdk:"expanded_permissions"`
	PackageName         types.String `tfsdk:"package_name"`
}

// grantID is the identifier used to address the grant: the resolved package
// name when known, otherwise the configured app_id.
func (m appIAMResourceModel) grantID() string {
	if m.PackageName.ValueString() != "" {
		return m.PackageName.ValueString()
	}
	return m.AppID.ValueString()
}

// grantPackageName returns the grant's package name, or null for draft apps
// which do not have one yet.
func grantPackageName(grant *androidpublisher.Grant) types.String {
	if grant.PackageName == "" {
		return types.StringNull()
	}
	return types.StringValue(grant.PackageName)
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
				MarkdownDescription: "The app to grant access to: either the numeric app ID or the package name (e.g. `uk.gov.android.onelogin`). The numeric app ID keeps working after a draft app is published, as the provider resolves and records the package name in `package_name`.",
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
			"package_name": schema.StringAttribute{
				MarkdownDescription: "The package name of the app, as reported by Google Play. This is empty for draft apps which do not have one yet.",
				Computed:            true,
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

	client, ok := req.ProviderData.(*GooglePlayClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *GooglePlayClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
	permissions := []AppLevelPermission{}
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

	permissions := []AppLevelPermission{}
	diag := data.Permissions.ElementsAs(ctx, &permissions, false)
	resp.Diagnostics.Append(diag...)

	grant, err := r.client.GrantAccess(
		ctx,
		data.UserID.ValueString(),
		data.AppID.ValueString(),
		permissions,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to grant access to app:",
			err.Error(),
		)
		return
	}

	// user_id and app_id are kept as planned: the grant name only carries them
	// as an email and a package name (or app ID for draft apps).
	data.PackageName = grantPackageName(grant)
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

	var found *androidpublisher.Grant
	for _, user := range users {
		if user.Email != userID {
			continue
		}
		for _, grant := range user.Grants {
			if grantMatchesApp(grant, appID) || (data.PackageName.ValueString() != "" && grantMatchesApp(grant, data.PackageName.ValueString())) {
				found = grant
				break
			}
		}
		break
	}

	if found == nil {
		found = r.resolvePublishedGrant(ctx, &data)
	}

	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var diag diag.Diagnostics
	data.PackageName = grantPackageName(found)
	data.ExpandedPermissions, diag = types.SetValueFrom(ctx, types.StringType, found.AppLevelPermissions)
	resp.Diagnostics.Append(diag...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// resolvePublishedGrant handles a draft app that has since been published: its
// grant is then listed under the package name, so a numeric app_id no longer
// matches the list. The API still accepts the numeric ID when addressing the
// grant directly, so re-send the recorded permissions (no change) to resolve it
// to the grant, which carries the package name. Returns nil if the grant does
// not exist or there are no recorded permissions to re-send.
func (r *AppIAMResource) resolvePublishedGrant(ctx context.Context, data *appIAMResourceModel) *androidpublisher.Grant {
	if data.ExpandedPermissions.IsNull() || data.ExpandedPermissions.IsUnknown() {
		return nil
	}

	recorded := []string{}
	if data.ExpandedPermissions.ElementsAs(ctx, &recorded, false).HasError() || len(recorded) == 0 {
		return nil
	}
	permissions := make([]AppLevelPermission, len(recorded))
	for i, p := range recorded {
		permissions[i] = AppLevelPermission(p)
	}

	grant, err := r.client.ModifyAccess(ctx, data.UserID.ValueString(), data.AppID.ValueString(), permissions)
	if err != nil {
		tflog.Debug(ctx, "grant not found when resolving by app_id", map[string]any{"error": err.Error()})
		return nil
	}
	return grant
}

// ImportState imports an existing grant using the ID "EMAIL/APP_ID".
// Only identity attributes are set here; Read populates the rest.
func (r *AppIAMResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	email, appID, found := strings.Cut(req.ID, "/")
	if !found || email == "" || appID == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier in the format \"email/app_id\", got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), email)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("app_id"), appID)...)
}

func (r *AppIAMResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data appIAMResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	permissions := []AppLevelPermission{}
	diag := data.Permissions.ElementsAs(ctx, &permissions, false)
	resp.Diagnostics.Append(diag...)

	grant, err := r.client.ModifyAccess(
		ctx,
		data.UserID.ValueString(),
		data.AppID.ValueString(),
		permissions,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update IAM permissions",
			err.Error(),
		)
		return
	}

	// user_id and app_id are kept as planned: the grant name only carries them
	// as an email and a package name (or app ID for draft apps).
	data.PackageName = grantPackageName(grant)
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
		ctx,
		data.UserID.ValueString(),
		data.grantID(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete permission, got error: %s", err))
		return
	}
}
